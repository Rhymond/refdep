package refdep

import (
	"errors"
	"fmt"
	"reflect"
)

const tagName = "refdep"

func (c *Container) Populate(v interface{}, deps ...string) error {
	tref := reflect.TypeOf(v)
	if tref.Kind() != reflect.Pointer {
		return errors.New("given parameter must be a pointer")
	}

	tref = tref.Elem()
	vref := reflect.ValueOf(v).Elem()
	if tref.Kind() == reflect.Struct {
		return c.populateStruct(vref, tref)
	}

	if tref.Kind() == reflect.Func {
		return errors.New("func not yet implemented")
	}

	val, ok := c.dependencies[refval]
	if !ok {
		return fmt.Errorf("dependency %q is not injected", refval)
	}

	if !vref.CanSet() {
		return fmt.Errorf("unable to set field %q", tref.Field(i).Name)
	}

	vref.Set(val)

	return nil
}

func (c *Container) populateStruct(vref reflect.Value, tref reflect.Type) error {
	for i := 0; i < tref.NumField(); i++ {
		refval := tref.Field(i).Tag.Get(tagName)
		if refval == "" {
			continue
		}

		val, ok := c.dependencies[refval]
		if !ok {
			return fmt.Errorf("dependency %q is not injected", refval)
		}

		if !vref.Field(i).CanSet() {
			return fmt.Errorf("unable to set field %q", tref.Field(i).Name)
		}

		vref.Field(i).Set(val)
	}

	return nil
}
