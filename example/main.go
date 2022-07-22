package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/i0Ek3/ormie"
	// register sqlite3 driver, required
	_ "github.com/mattn/go-sqlite3"
)

func main() {
    fmt.Println("---------------------------")
    fmt.Println("sqlUsage()")
    fmt.Println("---------------------------")
	sqlUsage()
    fmt.Println("---------------------------")
    fmt.Println("ormieTest()")
    fmt.Println("---------------------------")
	ormieTest()
}

func ormieTest() {
	e, _ := ormie.NewEngine("sqlite3", "ormie.db")
	defer e.Close()
	s := e.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}

func sqlUsage() {
	// Open used to connect sqlite3 database
	db, _ := sql.Open("sqlite3", "ormie.db")
	defer func() { _ = db.Close() }()

	// Exec used to execute SQL statements
	_, _ = db.Exec("DROP TABLE IF EXISTS User;")
	_, _ = db.Exec("CREATE TABLE User(Name text);")
	// placeholders ? are generally used to prevent SQL injection
	result, err := db.Exec("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam")
	if err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}

	// Query/QuerRow used to query SQL statements but the
	// former can return multiple records, and the latter
	// only returns one record which type is *sql.Row
	row := db.QueryRow("SELECT Name FROM User LIMIT 1")
	var name string
	if err := row.Scan(&name); err == nil {
		log.Println(name)
	}
}
