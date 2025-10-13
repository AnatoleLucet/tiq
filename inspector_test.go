package tiq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInspect(t *testing.T) {
	type TestStruct struct {
		Field1 string `json:"field1" db:"field_1"`
		Field2 int    `json:"field2" db:"field_2"`
	}

	t.Run("accepts struct value", func(t *testing.T) {
		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)
		assert.NotNil(t, inspector)
	})

	t.Run("accepts pointer to struct", func(t *testing.T) {
		inspector, err := Inspect(&TestStruct{})
		assert.NoError(t, err)
		assert.NotNil(t, inspector)
	})

	t.Run("returns error when passed nil", func(t *testing.T) {
		inspector, err := Inspect(nil)
		assert.ErrorIs(t, err, ErrNilValue)
		assert.Nil(t, inspector)
	})

	t.Run("returns error when passed non-struct", func(t *testing.T) {
		inspector, err := Inspect(123)
		assert.ErrorIs(t, err, ErrNotAStruct)
		assert.Nil(t, inspector)
	})

	t.Run("returns error when passed string", func(t *testing.T) {
		inspector, err := Inspect("not a struct")
		assert.ErrorIs(t, err, ErrNotAStruct)
		assert.Nil(t, inspector)
	})

	t.Run("returns error when passed slice", func(t *testing.T) {
		inspector, err := Inspect([]int{1, 2, 3})
		assert.ErrorIs(t, err, ErrNotAStruct)
		assert.Nil(t, inspector)
	})

	t.Run("returns error when passed map", func(t *testing.T) {
		inspector, err := Inspect(map[string]string{"key": "value"})
		assert.ErrorIs(t, err, ErrNotAStruct)
		assert.Nil(t, inspector)
	})
}

func TestInspector_Fields(t *testing.T) {
	t.Run("returns all fields from struct", func(t *testing.T) {
		type TestStruct struct {
			Field1 string `json:"field1" db:"field_1"`
			Field2 int    `json:"field2" db:"field_2"`
		}

		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		fields := inspector.Fields()
		assert.Len(t, fields, 2)
		assert.Equal(t, "Field1", fields[0].Name)
		assert.Equal(t, "Field2", fields[1].Name)
	})

	t.Run("returns all fields from pointer to struct", func(t *testing.T) {
		type TestStruct struct {
			Field1 string `json:"field1" db:"field_1"`
			Field2 int    `json:"field2" db:"field_2"`
		}

		inspector, err := Inspect(&TestStruct{})
		assert.NoError(t, err)

		fields := inspector.Fields()
		assert.Len(t, fields, 2)
		assert.Equal(t, "Field1", fields[0].Name)
		assert.Equal(t, "Field2", fields[1].Name)
	})

	t.Run("returns empty slice for struct with no fields", func(t *testing.T) {
		type EmptyStruct struct{}

		inspector, err := Inspect(EmptyStruct{})
		assert.NoError(t, err)

		fields := inspector.Fields()
		assert.Empty(t, fields)
	})

	t.Run("returns fields with correct types", func(t *testing.T) {
		type TestStruct struct {
			StringField string
			IntField    int
			BoolField   bool
			FloatField  float64
		}

		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		fields := inspector.Fields()
		assert.Len(t, fields, 4)
		assert.Equal(t, "string", fields[0].StructField.Type.Name())
		assert.Equal(t, "int", fields[1].StructField.Type.Name())
		assert.Equal(t, "bool", fields[2].StructField.Type.Name())
		assert.Equal(t, "float64", fields[3].StructField.Type.Name())
	})
}

func TestInspector_Field(t *testing.T) {
	type TestStruct struct {
		Field1 string `json:"field1" db:"field_1"`
		Field2 int    `json:"field2" db:"field_2"`
	}

	t.Run("returns field when it exists", func(t *testing.T) {
		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)
		assert.NotNil(t, field)
		assert.Equal(t, "Field1", field.Name)
	})

	t.Run("returns false when field does not exist", func(t *testing.T) {
		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("NonExistentField")
		assert.False(t, ok)
		assert.Nil(t, field)
	})

	t.Run("returns correct field from multiple fields", func(t *testing.T) {
		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("Field2")
		assert.True(t, ok)
		assert.NotNil(t, field)
		assert.Equal(t, "Field2", field.Name)
		assert.Equal(t, "int", field.StructField.Type.Name())
	})

	t.Run("field lookup is case sensitive", func(t *testing.T) {
		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("field1")
		assert.False(t, ok)
		assert.Nil(t, field)
	})
}

func TestInspector_FieldOperations(t *testing.T) {
	type TestStruct struct {
		Field1 string `json:"field1" db:"field_1"`
		Field2 int    `json:"field2" db:"field_2"`
	}

	t.Run("can retrieve and modify field value", func(t *testing.T) {
		testStruct := TestStruct{Field1: "initial"}
		inspector, err := Inspect(&testStruct)
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		err = field.Set("new value")
		assert.NoError(t, err)
		assert.Equal(t, "new value", testStruct.Field1)
	})

	t.Run("can retrieve field tags", func(t *testing.T) {
		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("Field1")
		assert.True(t, ok)

		tagValue, ok := field.Tag("json")
		assert.True(t, ok)
		assert.Equal(t, "field1", tagValue)

		tagValue, ok = field.Tag("db")
		assert.True(t, ok)
		assert.Equal(t, "field_1", tagValue)
	})

	t.Run("can retrieve all field tags", func(t *testing.T) {
		inspector, err := Inspect(TestStruct{})
		assert.NoError(t, err)

		field, ok := inspector.Field("Field2")
		assert.True(t, ok)

		tags, err := field.Tags()
		assert.NoError(t, err)
		assert.Len(t, tags, 2)
		assert.Equal(t, "field2", tags["json"])
		assert.Equal(t, "field_2", tags["db"])
	})
}
