package tiq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKv(t *testing.T) {
	t.Run("parses key=value pair", func(t *testing.T) {
		k, v := kv("key=value")
		assert.Equal(t, "key", k)
		assert.Equal(t, "value", v)
	})

	t.Run("parses key without value", func(t *testing.T) {
		k, v := kv("key")
		assert.Equal(t, "key", k)
		assert.Empty(t, v)
	})

	t.Run("trims whitespace", func(t *testing.T) {
		k, v := kv("  key  =  value  ")
		assert.Equal(t, "key", k)
		assert.Equal(t, "value", v)
	})

	t.Run("handles value with equals sign", func(t *testing.T) {
		k, v := kv("key=value=extra")
		assert.Equal(t, "key", k)
		assert.Equal(t, "value=extra", v)
	})
}

func TestFnGet(t *testing.T) {
	t.Run("returns value for existing key", func(t *testing.T) {
		result, err := fnGet("key1=value1,key2=value2", "key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", result)
	})

	t.Run("returns nil for non-existing key", func(t *testing.T) {
		result, err := fnGet("key1=value1", "key2")
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("returns error with wrong number of arguments", func(t *testing.T) {
		_, err := fnGet("key1=value1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exactly 2 arguments")
	})

	t.Run("returns error when first argument is not string", func(t *testing.T) {
		_, err := fnGet(123, "key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "first argument must be a string")
	})
}

func TestFnFirst(t *testing.T) {
	t.Run("returns first value from key=value pairs", func(t *testing.T) {
		result, err := fnFirst("key1=value1,key2=value2")
		assert.NoError(t, err)
		assert.Equal(t, "value1", result)
	})

	t.Run("returns key when no value", func(t *testing.T) {
		result, err := fnFirst("key1,key2")
		assert.NoError(t, err)
		assert.Equal(t, "key1", result)
	})

	t.Run("returns error when argument is not string", func(t *testing.T) {
		_, err := fnFirst(123)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "argument must be a string")
	})
}

func TestFnLast(t *testing.T) {
	t.Run("returns last value from key=value pairs", func(t *testing.T) {
		result, err := fnLast("key1=value1,key2=value2")
		assert.NoError(t, err)
		assert.Equal(t, "value2", result)
	})

	t.Run("returns key when last has no value", func(t *testing.T) {
		result, err := fnLast("key1=value1,key2")
		assert.NoError(t, err)
		assert.Equal(t, "key2", result)
	})

	t.Run("returns error when argument is not string", func(t *testing.T) {
		_, err := fnLast(123)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "argument must be a string")
	})
}

func TestFnNth(t *testing.T) {
	t.Run("returns nth value", func(t *testing.T) {
		result, err := fnNth("key1=value1,key2=value2,key3=value3", 1)
		assert.NoError(t, err)
		assert.Equal(t, "value2", result)
	})

	t.Run("returns key when nth has no value", func(t *testing.T) {
		result, err := fnNth("key1,key2=value2,key3", 2)
		assert.NoError(t, err)
		assert.Equal(t, "key3", result)
	})

	t.Run("returns nil for out of bounds index", func(t *testing.T) {
		result, err := fnNth("key1=value1", 5)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("returns error when second argument is not int", func(t *testing.T) {
		_, err := fnNth("key=value", "0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "second argument must be an integer")
	})
}

func TestFnHas(t *testing.T) {
	t.Run("returns true when key exists", func(t *testing.T) {
		result, err := fnHas("key1=value1,key2=value2", "key1")
		assert.NoError(t, err)
		assert.True(t, result.(bool))
	})

	t.Run("returns false when key does not exist", func(t *testing.T) {
		result, err := fnHas("key1=value1,key2=value2", "key3")
		assert.NoError(t, err)
		assert.False(t, result.(bool))
	})

	t.Run("returns error when second argument is not string", func(t *testing.T) {
		_, err := fnHas("key=value", 123)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "second argument must be a string")
	})
}

func TestFnSplit(t *testing.T) {
	t.Run("splits string by separator", func(t *testing.T) {
		result, err := fnSplit("a,b,c", ",")
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, result)
	})

	t.Run("trims whitespace from parts", func(t *testing.T) {
		result, err := fnSplit(" a , b , c ", ",")
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, result)
	})

	t.Run("returns error when first argument is not string", func(t *testing.T) {
		_, err := fnSplit(123, ",")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "first argument must be a string")
	})
}

func TestFnDefault(t *testing.T) {
	t.Run("returns first value when not nil", func(t *testing.T) {
		result, err := fnDefault("value", "default")
		assert.NoError(t, err)
		assert.Equal(t, "value", result)
	})

	t.Run("returns default when first is nil", func(t *testing.T) {
		result, err := fnDefault(nil, "default")
		assert.NoError(t, err)
		assert.Equal(t, "default", result)
	})
}

func TestCompile(t *testing.T) {
	t.Run("compiles simple variable reference", func(t *testing.T) {
		program, err := compile("json")
		assert.NoError(t, err)
		assert.NotNil(t, program)
	})

	t.Run("compiles function call", func(t *testing.T) {
		program, err := compile("get(db, 'table')")
		assert.NoError(t, err)
		assert.NotNil(t, program)
	})

	t.Run("compiles has() function", func(t *testing.T) {
		program, err := compile("has(validate, 'required')")
		assert.NoError(t, err)
		assert.NotNil(t, program)
	})

	t.Run("returns error for invalid expression", func(t *testing.T) {
		_, err := compile("invalid(((")
		assert.ErrorIs(t, err, ErrCompileTag)
		assert.Contains(t, err.Error(), "failed to compile")
	})
}

func TestParseTags(t *testing.T) {
	type Schema struct {
		Name     string `tag:"json"`
		Table    string `tag:"get(db, 'table')"`
		Required bool   `tag:"has(validate, 'required')"`
	}

	t.Run("parses simple tag reference", func(t *testing.T) {
		tags := map[string]string{
			"json": "user_name",
		}

		schema, err := parseTags[Schema](tags)
		assert.NoError(t, err)
		assert.NotNil(t, schema)
		assert.Equal(t, "user_name", schema.Name)
	})

	t.Run("parses get() function", func(t *testing.T) {
		tags := map[string]string{
			"db": "table=users,column=name",
		}

		schema, err := parseTags[Schema](tags)
		assert.NoError(t, err)
		assert.NotNil(t, schema)
		assert.Equal(t, "users", schema.Table)
	})

	t.Run("parses has() function", func(t *testing.T) {
		tags := map[string]string{
			"validate": "required,min=5",
		}

		schema, err := parseTags[Schema](tags)
		assert.NoError(t, err)
		assert.NotNil(t, schema)
		assert.True(t, schema.Required)
	})

	t.Run("parses multiple schema fields", func(t *testing.T) {
		tags := map[string]string{
			"json":     "name",
			"db":       "table=users",
			"validate": "required",
		}

		schema, err := parseTags[Schema](tags)
		assert.NoError(t, err)
		assert.NotNil(t, schema)
		assert.Equal(t, "name", schema.Name)
		assert.Equal(t, "users", schema.Table)
		assert.True(t, schema.Required)
	})

	t.Run("handles empty tags", func(t *testing.T) {
		tags := map[string]string{}

		schema, err := parseTags[Schema](tags)
		assert.NoError(t, err)
		assert.NotNil(t, schema)
		assert.Empty(t, schema.Name)
		assert.Empty(t, schema.Table)
		assert.False(t, schema.Required)
	})

	t.Run("returns error when compile fails", func(t *testing.T) {
		type BadSchema struct {
			Field string `tag:"invalid((("`
		}

		_, err := parseTags[BadSchema](map[string]string{"json": "value"})
		assert.ErrorIs(t, err, ErrCompileTag)
	})

	t.Run("returns error when conversion fails", func(t *testing.T) {
		type IntSchema struct {
			Port int `tag:"json"`
		}

		tags := map[string]string{
			"json": "not_a_number",
		}

		_, err := parseTags[IntSchema](tags)
		assert.ErrorIs(t, err, ErrCannotConvert)
	})
}
