package tiq

import (
	"reflect"
)

type Inspector struct {
	value reflect.Value
}

// Inspect takes a struct or pointer to struct and returns an Inspector
// that can be used to inspect the struct's fields and tags.
func Inspect(value any) (*Inspector, error) {
	if value == nil {
		return nil, ErrNilValue
	}

	if !isStruct(value) {
		return nil, ErrNotAStruct
	}

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	return &Inspector{v}, nil
}

// Fields returns every fields of the struct.
func (i *Inspector) Fields() []*Field {
	fields := []*Field{}

	v := i.value
	t := i.value.Type()

	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, &Field{
			v.Field(i),
			t.Field(i),
		})
	}

	return fields
}

// Field returns the field with the given name, or nil if it doesn't exist.
func (i *Inspector) Field(name string) (*Field, bool) {
	v := i.value
	t := i.value.Type()

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Name == name {
			return &Field{
				v.Field(i),
				t.Field(i),
			}, true
		}
	}

	return nil, false
}

func isStruct(v any) bool {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Struct {
		return true
	}
	if typ.Kind() == reflect.Pointer && typ.Elem().Kind() == reflect.Struct {
		return true
	}

	return false
}
