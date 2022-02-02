package refdep

import (
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
			if ret.IsNil() {
				continue
			}

			ret = ret.Elem()
		}

		if ret.Kind() != reflect.Struct {
			continue
		}

		name := ret.Type().PkgPath() + nameGlue + ret.Type().Name()
		c.dependencies[name] = ret.Interface()
	}

	return nil
}
