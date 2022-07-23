package ormie

import (
	"database/sql"

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
