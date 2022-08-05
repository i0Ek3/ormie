package session

import (
	"database/sql"
	"strings"

	"github.com/i0Ek3/ormie/clause"
	"github.com/i0Ek3/ormie/dialect"
	"github.com/i0Ek3/ormie/log"
	"github.com/i0Ek3/ormie/schema"
)

type Session struct {
	// db is the pointer returned after the Sql.Open()
	// method successfully connects to the database
	db      *sql.DB
	dialect dialect.Dialect
	// transaction
	tx       *sql.Tx
	refTable *schema.Schema
	clause   clause.Clause
	// sql used to concatenate SQL Statements
	sql strings.Builder
	// sqlVars is the corresponding value of
	// the placeholder in the SQL statement
	sqlVars []any
	// hookGraceful denotes which method to use for hooking
	hookGraceful bool
}

func New(db *sql.DB, dialect dialect.Dialect, hookGraceful bool) *Session {
	return &Session{
		db:           db,
		dialect:      dialect,
		hookGraceful: true,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
	s.hookGraceful = false
}

type CommonDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

// this line used to check if *sql.DB/*sql.Tx type implement CommonDB interface
var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) Raw(sql string, values ...any) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// QueryRow gets a record from db
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// Query get all records from db
func (s *Session) Query() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
