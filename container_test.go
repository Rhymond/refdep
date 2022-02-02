package refdep

import (
	"reflect"
	"testing"
)

type Foo struct {
	Bar string
}

func TestContainer_Add(t *testing.T) {
	t.Run("return value with pointer", func(t *testing.T) {
		c := New()
		foo := Foo{Bar: "baz"}
		err := c.Add(func() *Foo {
			return &foo
		})

		if err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}

		if !reflect.DeepEqual(c.dependencies["refdep"+nameGlue+"Foo"], foo) {
			t.Errorf("expected Foo dependency to be the same as stored")
		}
	})

	t.Run("return value without pointer", func(t *testing.T) {
		c := New()
		foo := Foo{Bar: "baz"}
		err := c.Add(func() Foo {
			return foo
		})

		if err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}

		if !reflect.DeepEqual(c.dependencies["refdep"+nameGlue+"Foo"], foo) {
			t.Errorf("expected Foo dependency to be the same as stored")
		}
	})

	t.Run("return value has an error", func(t *testing.T) {
		c := New()
		foo := Foo{Bar: "baz"}
		err := c.Add(func() (Foo, error) {
			return foo, nil
		})

		if err != nil {
			t.Errorf("error is not expected, but got %s", err)
		}

		if !reflect.DeepEqual(c.dependencies["refdep"+nameGlue+"Foo"], foo) {
			t.Errorf("expected Foo dependency to be the same as stored")
		}
	})
}
