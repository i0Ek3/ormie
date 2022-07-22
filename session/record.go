package session

import (
	"reflect"

	"github.com/i0Ek3/ormie/clause"
)

// Insert flatten the value of each field of the existing object
func (s *Session) Insert(values ...any) (int64, error) {
	recordValues := make([]any, 0)
	for _, value := range values {
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
	return result.RowsAffected()
}

// Find constructs an object from the values ​​of the flattened fields
func (s *Session) Find(values any) error {
	dstSlice := reflect.Indirect(reflect.ValueOf(values))
	// get the type of a single element of a slice
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
		// row record to each field in values ​​in turn
		if err := rows.Scan(values...); err != nil {
			return err
		}
		dstSlice.Set(reflect.Append(dstSlice, dst))
	}
	return rows.Close()
}
