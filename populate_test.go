package refdep

import (
	"testing"
)

func TestContainer_Populate(t *testing.T) {
	type inject struct {
		name string
		val  interface{}
	}

	type foo struct {
		Bar string `refdep:"bar"`
	}

	type unexportedFoo struct {
		bar string `refdep:"bar"`
	}

	var str string

	tests := []struct {
		name   string
		inject []inject
		v      interface{}
		assert func(t *testing.T, v interface{}, err error)
	}{
		{
			name: "only pointer is allowed",
			v:    "just a string",
			assert: func(t *testing.T, v interface{}, err error) {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}
			},
		},
		{
			name: "no field with tag presented",
			v: &struct {
				Foo string `json:"foo"`
			}{},
			assert: func(t *testing.T, v interface{}, err error) {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
			},
		},
		{
			name: "missing dependecy",
			v: &struct {
				Foo string `refdep:"foo"`
			}{},
			assert: func(t *testing.T, v interface{}, err error) {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}
			},
		},
		{
			name: "populate injected dependency",
			inject: []inject{
				{
					name: "bar",
					val:  "baz",
				},
			},
			v: &foo{},
			assert: func(t *testing.T, v interface{}, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				if f := v.(*foo); f.Bar != "baz" {
					t.Errorf("expected foo to have fields with populated values, foo = %+v", f)
				}
			},
		},
		{
			name: "string can be used for v",
			inject: []inject{
				{
					name: "bar",
					val:  "baz",
				},
			},
			v: &unexportedFoo{},
			assert: func(t *testing.T, v interface{}, err error) {
				if err == nil {
					t.Errorf("expected error for struct with unexported fields")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			di := New()
			for _, i := range test.inject {
				if err := di.Inject(i.name, i.val); err != nil {
					t.Fatalf("unexpected inject error: %s", err)
				}
			}

			err := di.Populate(test.v)
			test.assert(t, test.v, err)
		})
	}
}
