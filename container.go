package refdep

import (
	"errors"
	"fmt"
	"reflect"
)

type dependencies map[string]reflect.Value

type Container struct {
	dependencies dependencies
}

func New() *Container {
	return &Container{dependencies: make(dependencies)}
}

func (c *Container) Inject(name string, injector interface{}, deps ...string) error {
	if _, ok := c.dependencies[name]; ok {
		return fmt.Errorf("dependency name %q is already reserved", name)
	}

	ref := reflect.TypeOf(injector)
	val := reflect.ValueOf(injector)

	if ref.Kind() != reflect.Func {
		c.dependencies[name] = val
		return nil
	}

	if ref.NumOut() < 1 || ref.NumOut() > 2 {
		return fmt.Errorf("function signature can only consist of 1 or 2 return values")
	}

	ins, err := c.in(injector, deps...)
	if err != nil {
		return err
	}

	rets := val.Call(ins)

	c.dependencies[name] = rets[0]
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

func (c *Container) Eject(name string) (interface{}, error) {
	val, ok := c.dependencies[name]
	if !ok {
		return nil, fmt.Errorf("dependency do not exist with name %q", name)
	}

	return val.Interface(), nil
}

func (c *Container) in(fn interface{}, deps ...string) ([]reflect.Value, error) {
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

		switch p.Kind() {
		case reflect.Pointer, reflect.Interface:
			p = p.Elem()
		}

		p.Set(val)
		ins[i] = p
	}

	return ins, nil
}
