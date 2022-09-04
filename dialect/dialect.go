// Package dialect abstract the difference between different databases
package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	// DataTypeOf converts the type of Go language to the data type of the database
	DataTypeOf(typ reflect.Value) string
	// TableExistSQL returns whether a table exists
	TableExistSQL(tableName string) (string, []any)
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]

	return
}
