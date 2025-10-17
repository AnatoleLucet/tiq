package tiq

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/AnatoleLucet/as"
)

type Field struct {
	reflect.Value
	reflect.StructField
}

// Tags parses and returns every tag of the field as a map.
func (f *Field) Tags() (map[string]string, error) {
	var tags = make(map[string]string)

	structTag := string(f.StructField.Tag)
	if structTag == "" {
		return tags, nil
	}

	for segment := range strings.FieldsSeq(structTag) {
		colonIdx := strings.Index(segment, ":")
		if colonIdx == -1 {
			continue
		}

		key := segment[:colonIdx]
		value, ok := f.StructField.Tag.Lookup(key)
		if !ok {
			continue
		}

		tags[key] = value
	}

	return tags, nil

}

// Tag returns the tag value of the given name and whether it was found or not.
func (f *Field) Tag(name string) (string, bool) {
	return f.StructField.Tag.Lookup(name)
}

// Set updates the field's value to the provided value.
func (f *Field) Set(value any) error {
	if !f.Value.CanSet() {
		return ErrFieldNotSettable
	}

	v := reflect.ValueOf(value)
	if !v.CanConvert(f.Value.Type()) {
		return fmt.Errorf("%w: cannot convert %s to %s", ErrCannotConvert, v.Type(), f.Value.Type())
	}

	f.Value.Set(v.Convert(f.Value.Type()))
	return nil
}

// SetFrom updates the field's value to the provided value after converting it to the appropriate type.
// See as.Type for supported conversions.
func (f *Field) SetFrom(value any) error {
	typ := f.Value.Type()
	isPtr := typ.Kind() == reflect.Pointer

	if isPtr {
		typ = typ.Elem()
	}

	v, err := as.Type(typ, value)
	if err != nil {
		return fmt.Errorf("%w: cannot convert %T to %s: %v", ErrCannotConvert, value, f.Value.Type(), err)
	}

	if isPtr {
		ptr := reflect.New(typ)
		ptr.Elem().Set(reflect.ValueOf(v).Convert(typ))
		v = ptr.Interface()
	}

	return f.Set(v)
}
