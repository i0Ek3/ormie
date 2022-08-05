package session

import (
	"reflect"

	"github.com/i0Ek3/ormie/log"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

type IBeforeQuery interface {
	BeforeQuery(s *Session) error
}

type IAfterQuery interface {
	AfterQuery(s *Session) error
}

type IBeforeUpdate interface {
	BeforeUpdate(s *Session) error
}

type IAfterUpdate interface {
	AfterUpdate(s *Session) error
}

type IBeforeDelete interface {
	BeforeDelete(s *Session) error
}

type IAfterDelete interface {
	AfterDelete(s *Session) error
}

type IBeforeInsert interface {
	BeforeInsert(s *Session) error
}

type IAfterInsert interface {
	AfterInsert(s *Session) error
}

func (s *Session) CallMethod(method string, value any) {
	if s.hookGraceful {
		// use the MethodByName to reflect the method of the object
		fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
		if value != nil {
			fm = reflect.ValueOf(value).MethodByName(method)
		}
		// construct the parameter of the object
		param := []reflect.Value{reflect.ValueOf(s)}
		if fm.IsValid() {
			// calling fm function
			if v := fm.Call(param); len(v) > 0 {
				if err, ok := v[0].Interface().(error); ok {
					log.Error(err)
				}
			}
		}
	} else {
		param := reflect.ValueOf(value)
		switch method {
		case AfterQuery:
			if i, ok := param.Interface().(IAfterQuery); ok {
				i.AfterQuery(s)
			}
		case BeforeQuery:
			if i, ok := param.Interface().(IBeforeQuery); ok {
				i.BeforeQuery(s)
			}
		case AfterInsert:
			if i, ok := param.Interface().(IAfterInsert); ok {
				i.AfterInsert(s)
			}
		case BeforeInsert:
			if i, ok := param.Interface().(IBeforeInsert); ok {
				i.BeforeInsert(s)
			}
		case AfterDelete:
			if i, ok := param.Interface().(IAfterDelete); ok {
				i.AfterDelete(s)
			}
		case BeforeDelete:
			if i, ok := param.Interface().(IBeforeDelete); ok {
				i.BeforeDelete(s)
			}
		case AfterUpdate:
			if i, ok := param.Interface().(IAfterUpdate); ok {
				i.AfterUpdate(s)
			}
		case BeforeUpdate:
			if i, ok := param.Interface().(IBeforeUpdate); ok {
				i.BeforeUpdate(s)
			}
		default:
			panic(any("Unsupported hooks"))
		}
		return
	}
}
