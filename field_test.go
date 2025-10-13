package tiq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField_Tags(t *testing.T) {
	t.Run("returns all tags when field has multiple tags", func(t *testing.T) {
		type TestStruct struct {
			Field1 string `json:"field1" db:"field_1" validate:"required"`
		}

		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		tags, err := field.Tags()
		assert.NoError(t, err)
		assert.Len(t, tags, 3)
		assert.Equal(t, "field1", tags["json"])
		assert.Equal(t, "field_1", tags["db"])
		assert.Equal(t, "required", tags["validate"])
	})

	t.Run("returns empty map when field has no tags", func(t *testing.T) {
		type TestStruct struct {
			Field1 string
		}

		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		tags, err := field.Tags()
		assert.NoError(t, err)
		assert.Empty(t, tags)
	})

	t.Run("returns tags when field has single tag", func(t *testing.T) {
		type TestStruct struct {
			Field1 string `json:"field1"`
		}

		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		tags, err := field.Tags()
		assert.NoError(t, err)
		assert.Len(t, tags, 1)
		assert.Equal(t, "field1", tags["json"])
	})

	t.Run("handles complex tag values", func(t *testing.T) {
		type TestStruct struct {
			Field1 string `json:"field1,omitempty" validate:"min=5,max=10"`
		}

		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		tags, err := field.Tags()
		assert.NoError(t, err)
		assert.Equal(t, "field1,omitempty", tags["json"])
		assert.Equal(t, "min=5,max=10", tags["validate"])
	})
}

func TestField_Tag(t *testing.T) {
	t.Run("returns tag value when tag exists", func(t *testing.T) {
		type TestStruct struct {
			Field1 string `json:"field1" db:"field_1"`
		}

		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		value, found := field.Tag("json")
		assert.True(t, found)
		assert.Equal(t, "field1", value)

		value, found = field.Tag("db")
		assert.True(t, found)
		assert.Equal(t, "field_1", value)
	})

	t.Run("returns false when tag does not exist", func(t *testing.T) {
		type TestStruct struct {
			Field1 string `json:"field1"`
		}

		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		value, found := field.Tag("nonexistent")
		assert.False(t, found)
		assert.Empty(t, value)
	})

	t.Run("returns false when field has no tags", func(t *testing.T) {
		type TestStruct struct {
			Field1 string
		}

		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		value, found := field.Tag("json")
		assert.False(t, found)
		assert.Empty(t, value)
	})
}

func TestField_Set(t *testing.T) {
	t.Run("sets string field value", func(t *testing.T) {
		type TestStruct struct {
			Field1 string
		}

		testStruct := TestStruct{}
		inspector, err := Inspect(&testStruct)
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		err = field.Set("new value")
		assert.NoError(t, err)
		assert.Equal(t, "new value", testStruct.Field1)
	})

	t.Run("sets int field value", func(t *testing.T) {
		type TestStruct struct {
			Field1 int
		}

		testStruct := TestStruct{}
		inspector, err := Inspect(&testStruct)
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		err = field.Set(42)
		assert.NoError(t, err)
		assert.Equal(t, 42, testStruct.Field1)
	})

	t.Run("sets bool field value", func(t *testing.T) {
		type TestStruct struct {
			Field1 bool
		}

		testStruct := TestStruct{}
		inspector, err := Inspect(&testStruct)
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		err = field.Set(true)
		assert.NoError(t, err)
		assert.Equal(t, true, testStruct.Field1)
	})

	t.Run("sets float field value", func(t *testing.T) {
		type TestStruct struct {
			Field1 float64
		}

		testStruct := TestStruct{}
		inspector, err := Inspect(&testStruct)
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		err = field.Set(3.14)
		assert.NoError(t, err)
		assert.Equal(t, 3.14, testStruct.Field1)
	})

	t.Run("converts compatible types", func(t *testing.T) {
		type TestStruct struct {
			Field1 int64
		}

		testStruct := TestStruct{}
		inspector, err := Inspect(&testStruct)
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		err = field.Set(int(42))
		assert.NoError(t, err)
		assert.Equal(t, int64(42), testStruct.Field1)
	})

	t.Run("returns error when field is not settable", func(t *testing.T) {
		type TestStruct struct {
			Field1 string
		}

		testStruct := TestStruct{}
		inspector, err := Inspect(testStruct) // not a pointer
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		err = field.Set("new value")
		assert.ErrorIs(t, err, ErrFieldNotSettable)
	})

	t.Run("returns error when type conversion is not possible", func(t *testing.T) {
		type TestStruct struct {
			Field1 int
		}

		testStruct := TestStruct{}
		inspector, err := Inspect(&testStruct)
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		err = field.Set("not a number")
		assert.ErrorIs(t, err, ErrCannotConvert)
		assert.Contains(t, err.Error(), "cannot convert")
		assert.Contains(t, err.Error(), "string")
		assert.Contains(t, err.Error(), "int")
	})

	t.Run("returns error when converting incompatible struct types", func(t *testing.T) {
		type StructA struct {
			Value string
		}
		type StructB struct {
			Value int
		}
		type TestStruct struct {
			Field1 StructA
		}

		testStruct := TestStruct{}
		inspector, err := Inspect(&testStruct)
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		err = field.Set(StructB{Value: 42})
		assert.ErrorIs(t, err, ErrCannotConvert)
	})

	t.Run("sets pointer field value", func(t *testing.T) {
		type TestStruct struct {
			Field1 *string
		}

		testStruct := TestStruct{}
		inspector, err := Inspect(&testStruct)
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		value := "test"
		err = field.Set(&value)
		assert.NoError(t, err)
		assert.NotNil(t, testStruct.Field1)
		assert.Equal(t, "test", *testStruct.Field1)
	})

	t.Run("sets slice field value", func(t *testing.T) {
		type TestStruct struct {
			Field1 []int
		}

		testStruct := TestStruct{}
		inspector, err := Inspect(&testStruct)
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		err = field.Set([]int{1, 2, 3})
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, testStruct.Field1)
	})

	t.Run("sets struct field value", func(t *testing.T) {
		type InnerStruct struct {
			Value string
		}
		type TestStruct struct {
			Field1 InnerStruct
		}

		testStruct := TestStruct{}
		inspector, err := Inspect(&testStruct)
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		err = field.Set(InnerStruct{Value: "test"})
		assert.NoError(t, err)
		assert.Equal(t, "test", testStruct.Field1.Value)
	})
}
