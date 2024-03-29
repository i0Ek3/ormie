package schema

import (
	"go/ast"
	"reflect"

	"github.com/i0Ek3/ormie/dialect"
)

type Field struct {
	// Field name
	Name string
	// Field type
	Type string
	// Restrictions
	Tag string
}

type Schema struct {
	// Mapped object
	Model any
	// Table name
	Name string
	// All fields
	Fields []*Field
	// All field names(columns)
	FieldNames []string
	// Record field name and Field mapping relationship
	fieldMap map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// RecordValues according to the order of the columns in the database,
// find the corresponding value from the object and tile it in order
func (schema *Schema) RecordValues(dst any) []any {
	dstValue := reflect.Indirect(reflect.ValueOf(dst))
	var fieldValues []any
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, dstValue.FieldByName(field.Name).Interface())
	}

	return fieldValues
}

type ITableName interface {
	TableName() string
}

// Parse parses any object into Schema instance
func Parse(dst any, d dialect.Dialect) *Schema {
	// Get the instance pointed to by the pointer through Indirect()
	modelType := reflect.Indirect(reflect.ValueOf(dst)).Type()
	var tableName string
	t, ok := dst.(ITableName)
	if !ok {
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}
	schema := &Schema{
		Model:    dst,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		// Get a specific field by subscripting
		f := modelType.Field(i)
		if !f.Anonymous && ast.IsExported(f.Name) {
			// Convert to database support field
			field := &Field{
				Name: f.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(f.Type))),
			}
			if v, ok := f.Tag.Lookup("ormie"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, f.Name)
			schema.fieldMap[f.Name] = field
		}
	}

	return schema
}
