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
				err := i.AfterQuery(s)
				if err != nil {
					return
				}
			}
		case BeforeQuery:
			if i, ok := param.Interface().(IBeforeQuery); ok {
				err := i.BeforeQuery(s)
				if err != nil {
					return
				}
			}
		case AfterInsert:
			if i, ok := param.Interface().(IAfterInsert); ok {
				err := i.AfterInsert(s)
				if err != nil {
					return
				}
			}
		case BeforeInsert:
			if i, ok := param.Interface().(IBeforeInsert); ok {
				err := i.BeforeInsert(s)
				if err != nil {
					return
				}
			}
		case AfterDelete:
			if i, ok := param.Interface().(IAfterDelete); ok {
				err := i.AfterDelete(s)
				if err != nil {
					return
				}
			}
		case BeforeDelete:
			if i, ok := param.Interface().(IBeforeDelete); ok {
				err := i.BeforeDelete(s)
				if err != nil {
					return
				}
			}
		case AfterUpdate:
			if i, ok := param.Interface().(IAfterUpdate); ok {
				err := i.AfterUpdate(s)
				if err != nil {
					return
				}
			}
		case BeforeUpdate:
			if i, ok := param.Interface().(IBeforeUpdate); ok {
				err := i.BeforeUpdate(s)
				if err != nil {
					return
				}
			}
		default:
			panic(any("Unsupported hooks"))
		}
		return
	}
}
