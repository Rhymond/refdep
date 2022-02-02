package refdep

import (
	"errors"
	"fmt"
	"reflect"
)

type Container struct {
	dependencies map[string]interface{}
}

func New() *Container {
	return &Container{dependencies: make(map[string]interface{})}
}

const nameGlue = "."

func (c *Container) Add(fn interface{}) error {
	ref := reflect.TypeOf(fn)
	val := reflect.ValueOf(fn)
	if ref.Kind() != reflect.Func {
		return fmt.Errorf("only functions allowed")
	}

	rets := val.Call([]reflect.Value{})
	for _, ret := range rets {
		if ret.Kind() == reflect.Pointer {
			ret = ret.Elem()
		}

		switch ret.Kind() {
		case reflect.Struct, reflect.Interface:
		default:
			continue
		}

		if ret.Kind() == reflect.Interface {
			if err, ok := ret.Interface().(error); ok {
				if err != nil {
					return err
				}
				continue
			}

			return errors.New("only struct allowed")
		}

		name := ret.Type().PkgPath() + nameGlue + ret.Type().Name()
		c.dependencies[name] = ret.Interface()
	}

	return nil
}
