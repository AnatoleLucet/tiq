<h1 align="center"><code>tiq</code></h1>

<p align="center">Modular Golang Struct tags parser that's actually useful.</p>

```go
type Config struct {
	Url  string `env:"name=URL, type=string"`
	Port int    `env:"name=PORT, type=port"`
}

type EnvSchema struct {
	Name string `tag:"env | get('name')"`
	Type string `tag:"env | get('type')"`
}

conf := Config{}

inspector, err := tiq.Inspect(&conf)
for _, field := range inspector.Fields() {
    env, err := tiq.Parse[EnvSchema](field)

    value := os.Getenv(env.Name)
    field.Set(validate(env.Type, value))
}

conf.Url // now set to the value of the URL env var
conf.Port // now set to the value of the PORT env var
```

## Installation

```bash
go get github.com/AnatoleLucet/tiq
```

## Usage

`tiq` is a modular Golang Struct tags parser. It offers a very simple [DSL](https://en.wikipedia.org/wiki/Domain-specific_language) to extract what you need from tags, and multiple APIs to [inspect](#tiqinspect) and [update](#tiqset) user defined structs.

The DSL is designed to be as straightforward as possible for you to pick up and get what's happening at a glance, even if you've never used `tiq` before.

```go
import (
    "github.com/AnatoleLucet/tiq"
)

// User defined struct with tags that can by parsed by the `load()` function.
type Config struct {
	Url  string `env:"name=URL, type=string, optional"`
	Port int    `env:"name=PORT, type=port, oneof=8080|3000|5000"`
}

func main() {
    conf, err := load(&Config{})

    conf.Url // Now set to the value of the URL env var.
    conf.Port // Now set to the value of the PORT env var.
}


// Define your schema and how to parse tags
type EnvSchema struct {
    // Each field containing a `tag:""` will be evaluated by tiq's DSL to
    // extract what you need from the user defined tags.
    // See "The DSL" section of the README to learn more.
	Name     string   `tag:"env | get('name')"`
	Type     string   `tag:"env | get('type')"`
	Optional bool     `tag:"env | has('optional')"`
	Oneof    []string `tag:"env | get('oneof') | split('|')"`
}

func load[T any](conf *T) (*T, error) {
    inspector, err := tiq.Inspect(conf)

    for _, field := range inspector.Fields() {
        env, err := tiq.Parse[EnvSchema](field)

        value := os.Getenv(env.Name)

        // Ideally you'd call a function to validate if the value
        // is correct according what was parsed in the `env` variable:
        // validate(env, value)

        field.Set(value)
    }

    return conf
}
```

### The DSL

The DSL is based on [ExprLang](https://expr-lang.org/), a simple but powerful expression language.

> Don't try to use functions from the official ExprLang docs, they probably won't work. Instead, take a look at `tiq`'s [function set](#functions) to find what you need!

#### Basic syntax

```bash
# The most simple expression would look like this:
`tag:"123"`
# where `tag:"..."` is the Golang tag tiq will pick up for evaluation,
# and `123` is the DSL expression tiq will evaluate.

# To get a tag's value, simply use the name of the tag you want to get:
`tag:"mytag"`
# it will return `mytag`'s content unaltered (e.g. if given `mytag:"content"`, the expression above will return `content`).

# Once you have the value you want to parse, you can use tiq's function set to extract entries and values from it:
`tag:"get(mytag, 'foo')"`
# here we pass `mytag`'s content to the `get()` function, and try to get the `foo` entry's value from it.
# So when given `mytag:"foo=bar"`, the expression above will return `bar` (the value of the `foo` entry).

# To chain one or more functions together, you can use ExprLang's pipe operator:
`tag:"mytag | get("foo") | default("baz")"`
# the pipe operator will pass the left operand's value as the first parameter the right operand.
# What this means is that `"foo=bar" | get("foo")` is equivalent to `get("foo=bar", "foo")`.

# To learn more about tiq's syntax, check out ExprLang's docs at https://expr-lang.org/docs/getting-started.
# But remember most functions from ExprLang won't work because tiq uses its own functions set (described below).
```

#### Functions

| Name        | Description                                                                                      | Usage                                    |
| ----------- | ------------------------------------------------------------------------------------------------ | ---------------------------------------- |
| `get()`     | Gets an entry's value from a comma-separated key-value list.                                     | `get("foo=1, bar=2", "foo") -> 1`        |
| `first()`   | Gets the first entry's value (or key if there's no value) from a comma-separated key-value list. | `first("foo=1, bar=2") -> 1`             |
| `last()`    | Gets the last entry's value (or key if there's no value) from a comma-separated key-value list.  | `last("foo=1, bar=2") -> 2`              |
| `nth()`     | Gets the nth entry's value (or key if there's no value) from a comma-separated key-value list.   | `nth("foo=1, bar=2", 0) -> 1`            |
| `has()`     | Returns true or false if the entry is present in a comma-separated key-value list.               | `has("foo=1, bar=2", "bar") -> true`     |
| `split()`   | Splits a string with the given separator.                                                        | `split("1\|2\|3\|4", "\|") -> [1 2 3 4]` |
| `default()` | Returns a default value if the given value if `nil`.                                             | `default(nil, "foo") -> "foo"`           |

### `tiq.Inspect`

The inspector helps you crawl through a struct's fields, read tags from them, and update values accordingly.

```go
inspector, err := tiq.Inspect(&mystruct)

// Get a field by name
field, ok := inspector.Field("Name")

field.Set("value") // update the field's value
field.Tag("mytag") // returns the content of `mytag:"content"`
field.Tags() // returns every tags of the field in a map[string]string

// Alternatively you could loop through every field on the struct:
for _, field := range inspector.Fields() {
    // field.Set("value")
}
```

### `tiq.Parse`

The parser is how you retrieve what you want from tags with `tiq`. It takes a schema and a `tiq.Field` to parse tags on.

```go
type EnvSchema struct {
	Name     string   `tag:"env | get('name')"`
	Optional bool     `tag:"env | has('optional')"`
	Oneof    []string `tag:"env | get('oneof') | split('|')"`
}

env, err := tiq.Parse[EnvSchema](field) // field is usually retrieved via tiq.Inspect

env.Name // if `field` has a tag `env:"name=foo"`, this will be set to "foo", else ""
env.Optional // if `field` has a tag `env:"optional"`, this will be set to true, else false
env.Oneof // if `field` has a tag `env:"oneof=one|two|three"`, this will be set to [one two three], else []
```

### `tiq.Get`

A simple static function to get a tag's content from anywhere.

```go
type User struct {
    Name string `json:"name,omitempty"`
}

var user User
json, err := tiq.Get(&user, "Name", "json")

json // "name,omitempty"
```

### `tiq.Set`

A simple static function to set a field's content from anywhere.

```go
type User struct {
    Name string `json:"name,omitempty"`
}

var user User
err := tiq.Set(&user, "Name", "Bob")

user.Name // "Bob"
```
