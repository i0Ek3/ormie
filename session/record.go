package session

import (
	"errors"
	"reflect"

	"github.com/i0Ek3/ormie/clause"
)

// Insert flatten the value of each field of the existing object
func (s *Session) Insert(values ...any) (int64, error) {
	recordValues := make([]any, 0)
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert, "")
	return result.RowsAffected()
}

// Find constructs an object from the values of the flattened fields
func (s *Session) Find(values any) error {
	s.CallMethod(BeforeQuery, "")
	dstSlice := reflect.Indirect(reflect.ValueOf(values))
	// get the type of single element of a slice
	dstType := dstSlice.Type().Elem()
	// Model according given parameters mapped out table structure by RefTable()
	table := s.Model(reflect.New(dstType).Elem().Interface()).RefTable()

	// according to the table structure, use clause to construct a SELECT statement
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	// query to all rows that meet the criteria
	rows, err := s.Raw(sql, vars...).Query()
	if err != nil {
		return err
	}
	// traverse each row of records
	for rows.Next() {
		// use reflection to create an instance dst of dstType
		dst := reflect.New(dstType).Elem()
		var values []any
		// flatten all fields of dst, and construct slice values.
		for _, name := range table.FieldNames {
			values = append(values, dst.FieldByName(name).Addr().Interface())
		}
		// assign the value of each column of the
		// row record to each field in values in turn
		if err := rows.Scan(values...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery, dst.Addr().Interface())
		dstSlice.Set(reflect.Append(dstSlice, dst))
	}
	return rows.Close()
}

func (s *Session) Update(kv ...any) (int64, error) {
	s.CallMethod(BeforeUpdate, "")
	// assert m whether it is a map
	m, ok := kv[0].(map[string]any)
	// if not, convert m into flatten key-value pair
	if !ok {
		m = make(map[string]any)
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, "")
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, "")
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, "")
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...any) *Session {
	var vars []any
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) First(value any) error {
	dst := reflect.Indirect(reflect.ValueOf(value))
	dstSlice := reflect.New(reflect.SliceOf(dst.Type())).Elem()
	if err := s.Limit(1).Find(dstSlice.Addr().Interface()); err != nil {
		return err
	}
	if dstSlice.Len() == 0 {
		return errors.New("not found")
	}
	dst.Set(dstSlice.Index(0))
	return nil
}
