package tiq

import "errors"

var (
	ErrNilValue      = errors.New("nil value provided")
	ErrNotAStruct    = errors.New("provided value is not a struct or pointer to struct")
	ErrCannotConvert = errors.New("cannot convert value")

	ErrFieldNotFound    = errors.New("field not found")
	ErrFieldNotSettable = errors.New("field is not settable")

	ErrCompileTag = errors.New("cannot compile tag")
)
