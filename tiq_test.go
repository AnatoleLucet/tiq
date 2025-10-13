package tiq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	type TestStruct struct {
		Field1 string `json:"field1" db:"field_1"`
		Field2 int    `json:"field2"`
	}

	t.Run("returns tag value when field and tag exist", func(t *testing.T) {
		testStruct := TestStruct{}

		value, ok := Get(testStruct, "Field1", "json")
		assert.True(t, ok)
		assert.Equal(t, "field1", value)
	})

	t.Run("returns false when field does not exist", func(t *testing.T) {
		testStruct := TestStruct{}

		value, ok := Get(testStruct, "NonExistentField", "json")
		assert.False(t, ok)
		assert.Empty(t, value)
	})

	t.Run("returns false when Inspect fails", func(t *testing.T) {
		value, ok := Get(nil, "Field1", "json")
		assert.False(t, ok)
		assert.Empty(t, value)
	})
}

func TestSet(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
	}

	t.Run("sets field value", func(t *testing.T) {
		testStruct := &TestStruct{}

		err := Set(testStruct, "Field1", "new value")
		assert.NoError(t, err)
		assert.Equal(t, "new value", testStruct.Field1)
	})

	t.Run("returns ErrFieldNotFound when field does not exist", func(t *testing.T) {
		testStruct := &TestStruct{}

		err := Set(testStruct, "NonExistentField", "value")
		assert.ErrorIs(t, err, ErrFieldNotFound)
	})

	t.Run("returns error when Inspect fails", func(t *testing.T) {
		err := Set(nil, "Field1", "value")
		assert.ErrorIs(t, err, ErrNilValue)
	})
}
