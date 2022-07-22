package ormie

import (
	"database/sql"

	"github.com/i0Ek3/ormie/log"
	"github.com/i0Ek3/ormie/session"
)

type Engine struct {
	db *sql.DB
}

// NewEngine connects the database and checks if it alive
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
	e = &Engine{db: db}
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
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}
