# ormie

`ormie` is a simple ORM framework for learning the fundamentals of ORM frameworks.

## Feature

- CURD for records
- Create/Drop/Migration database table
- Chained calls of query conditions
- Hooks and transaction support
- Support SQLite database

## Usage

```Go
package main

import (
    "github.com/i0Ek3/ormie/engine"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    e, _ := engine.NewEngine("sqlite3", "ormie.db")
	defer e.Close()
	s := e.NewSession()
	_, _ = s.Raw("/* Your SQL Statement here */").Exec()
    
    // ...
}
```

More details please check source code and usage example in folder example.

## Credit

[geektutu](https://github.com/geektutu)
