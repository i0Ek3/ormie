package schema

import (
	"go/ast"
	"reflect"

	"github.com/i0Ek3/ormie/dialect"
)

type Field struct {
	// field name
	Name string
	// field type
	Type string
	// restrictions
	Tag string
}

type Schema struct {
	// mapped object
	Model any
	// table name
	Name string
	// all fields
	Fields []*Field
	// all field names(columns)
	FieldNames []string
	// record field name and Field mapping relationship
	fieldMap map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// Parse parses any object into Schema instance
func Parse(dst any, d dialect.Dialect) *Schema {
	// get the instance pointed to by the pointer through Indirect()
	modelType := reflect.Indirect(reflect.ValueOf(dst)).Type()
	schema := &Schema{
		Model:    dst,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		// get a specific field by subscripting
		f := modelType.Field(i)
		if !f.Anonymous && ast.IsExported(f.Name) {
			// convert to database support field
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
