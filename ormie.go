package ormie

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/i0Ek3/ormie/dialect"
	"github.com/i0Ek3/ormie/log"
	"github.com/i0Ek3/ormie/session"
)

type Engine struct {
	db           *sql.DB
	dialect      dialect.Dialect
	hookGraceful bool
}

// NewEngine connects the database and checks if it alive, and also get
// the dialect corresponding to the driver
func NewEngine(driver string, src string) (e *Engine, err error) {
	db, err := sql.Open(driver, src)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s not found", driver)
		return
	}
	e = &Engine{
		db:      db,
		dialect: dial,
	}
	log.Info("Database connected")
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Database closed")
}

// NewSession creates session instance used to interact with database
// and then pass dialect to constructor New()
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect, e.hookGraceful)
}

type TxFunc func(*session.Session) (any, error)

func (e *Engine) Transaction(f TxFunc) (result any, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()
			err = s.Commit()
		}
	}()
	return f(s)
}

func findDiff(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate migrates the new fields into the new table and delete the old table
func (e *Engine) Migrate(value any) error {
	_, err := e.Transaction(func(s *session.Session) (result any, err error) {
		//
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).Query()
		columns, _ := rows.Columns()
		// calculate the new fields and deleted fields
		addCols := findDiff(table.FieldNames, columns)
		delCols := findDiff(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			// add new fields by ALTER
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}
		if len(delCols) == 0 {
			return
		}
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		// create a new table and delete the old table, then renamed new table to the old table
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}
