package refdep

import (
	"errors"
	"reflect"
	"testing"
)

type foo struct {
	Bar string
}

func TestContainer_Inject(t *testing.T) {
	const fooName = "foo"
	var must = func(v interface{}, err error) interface{} {
		if err != nil {
			return nil
		}

		return v
	}

	t.Run("return value with pointer", func(t *testing.T) {
		c := New()
		f := foo{Bar: "baz"}
		err := c.Inject(fooName, func() *foo {
			return &f
		})

		if err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}

		if !reflect.DeepEqual(must(c.Eject(fooName)), &f) {
			t.Errorf("expected foo dependency to be the same as stored")
		}
	})

	t.Run("return value without pointer", func(t *testing.T) {
		c := New()
		f := foo{Bar: "baz"}
		err := c.Inject(fooName, func() foo {
			return f
		})

		if err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}

		if !reflect.DeepEqual(must(c.Eject(fooName)), f) {
			t.Errorf("expected foo dependency to be the same as stored")
		}
	})

	t.Run("return value is nil error", func(t *testing.T) {
		c := New()
		f := foo{Bar: "baz"}
		err := c.Inject(fooName, func() (foo, error) {
			return f, nil
		})

		if err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}

		if !reflect.DeepEqual(must(c.Eject(fooName)), f) {
			t.Errorf("expected foo dependency to be the same as stored")
		}
	})

	t.Run("return value consists an error", func(t *testing.T) {
		c := New()
		err := c.Inject(fooName, func() (*foo, error) {
			return nil, errors.New("unexpected error")
		})

		if err == nil {
			t.Error("expected error, but got nil")
		}
	})

	t.Run("return dependency without the error", func(t *testing.T) {
		c := New()
		f := foo{Bar: "baz"}
		err := c.Inject(fooName, func() *foo {
			return &f
		})

		if err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}

		if !reflect.DeepEqual(must(c.Eject(fooName)), &f) {
			t.Errorf("expected foo dependency to be the same as stored")
		}
	})

	t.Run("return dependency without the error", func(t *testing.T) {
		c := New()
		f := foo{Bar: "baz"}
		err := c.Inject(fooName, func() *foo {
			return &f
		})

		if err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}

		if !reflect.DeepEqual(must(c.Eject(fooName)), &f) {
			t.Errorf("expected foo dependency to be the same as stored")
		}
	})

	t.Run("dependency injection", func(t *testing.T) {
		c := New()
		f := foo{Bar: "baz"}
		err := c.Inject(fooName, func() *foo {
			return &f
		})

		if err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}

		const strDep = "I'm string dependency"
		const barName = "bar"
		if err := c.Inject(barName, func(f1 *foo) string {
			if f1.Bar != f.Bar {
				t.Fatalf("injected and received *foo is not equal")
			}

			return strDep
		}, fooName); err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}

		injected := must(c.Eject(barName)).(string)
		if injected != strDep {
			t.Errorf("expected string dependency to be %s, but got %s", strDep, injected)
		}
	})

	t.Run("dependency injection with same name", func(t *testing.T) {
		c := New()
		f := foo{Bar: "baz"}
		err := c.Inject(fooName, func() *foo {
			return &f
		})

		if err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}

		err = c.Inject(fooName, func() *foo {
			return &f
		})

		if err == nil {
			t.Error("error is expected, but got nil")
		}
	})

	t.Run("injecting a string dependency", func(t *testing.T) {
		c := New()
		const depVal = "just a string"
		if err := c.Inject("string", depVal); err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}

		if err := c.Inject(fooName, func(dep string) *foo {
			if dep != depVal {
				t.Errorf("expected string dependency to be %s, but got %s", depVal, dep)
			}
			return &foo{Bar: dep}
		}, "string"); err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}
	})
}

func TestContainer_in(t *testing.T) {
	const fooName = "foo"
	t.Run("replace single parameter", func(t *testing.T) {
		f := foo{Bar: "baz"}
		c := Container{dependencies: dependencies{
			fooName: reflect.ValueOf(f),
		}}

		ins, err := c.in(func(f foo) {}, fooName)
		if err != nil {
			t.Errorf("unexpected error, got %s", err)
		}

		if len(ins) != 1 {
			t.Errorf("expected 1 reflect.Value parameter in slice, but got %d", len(ins))
		}
	})

	t.Run("dependency do not exist", func(t *testing.T) {
		c := Container{dependencies: dependencies{}}

		_, err := c.in(func(f foo) {}, fooName)
		if err == nil {
			t.Errorf("expected to return an error, got nil")
		}
	})
}
