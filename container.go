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

func (c *Container) Add(name string, fn interface{}, deps ...string) error {
	ref := reflect.TypeOf(fn)
	val := reflect.ValueOf(fn)
	if ref.Kind() != reflect.Func {
		return fmt.Errorf("only functions allowed")
	}

	if ref.NumOut() < 1 || ref.NumOut() > 2 {
		return fmt.Errorf("function signature can only consist of 1 or 2 return values")
	}

	ins, err := c.replaceIn(fn, deps...)
	if err != nil {
		return err
	}

	rets := val.Call(ins)
	c.dependencies[name] = rets[0].Interface()
	if ref.NumOut() == 1 {
		return nil
	}

	if rets[1].Kind() != reflect.Interface {
		return errors.New("second parameter can only be an error")
	}

	if err, ok := rets[1].Interface().(error); ok && err != nil {
		return err
	}

	return nil
}

func (c *Container) replaceIn(fn interface{}, deps ...string) ([]reflect.Value, error) {
	ref := reflect.TypeOf(fn)
	if ref.NumIn() != len(deps) {
		const msg = "function parameter count do not match dep count, parameters: %d, deps: %d"
		return []reflect.Value{}, fmt.Errorf(msg, ref.NumIn(), len(deps))
	}

	if ref.NumIn() == 0 {
		return []reflect.Value{}, nil
	}

	ins := make([]reflect.Value, ref.NumIn())
	for i := 0; i < ref.NumIn(); i++ {
		val, ok := c.dependencies[deps[i]]
		if !ok {
			return []reflect.Value{}, fmt.Errorf("dependecy %s is not initialized", deps[i])
		}

		param := ref.In(i)
		p := reflect.New(param)
		p.Elem().Set(reflect.ValueOf(val))
		ins[i] = p
	}

	return ins, nil
}
