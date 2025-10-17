package tiq

import (
	"errors"
	"fmt"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

func parseTags[Schema any](tags map[string]string) (*Schema, error) {
	tag := new(Schema)
	inspector, err := Inspect(tag)
	if err != nil {
		return nil, err
	}

	for _, f := range inspector.Fields() {
		expression, ok := f.Tag("tag")
		if !ok {
			continue
		}

		program, err := compile(expression)
		if err != nil {
			return nil, err
		}

		output, err := expr.Run(program, tags)
		if err != nil || output == nil {
			continue
		}

		err = f.SetFrom(output)
		if err != nil {
			return nil, err
		}
	}

	return tag, nil
}

func compile(expression string) (*vm.Program, error) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),

		expr.DisableAllBuiltins(),
		expr.Function("get", fnGet, new(func(string, string) (string, error))),
		expr.Function("first", fnFirst, new(func(string) (string, error))),
		expr.Function("last", fnLast, new(func(string) (string, error))),
		expr.Function("nth", fnNth, new(func(string, int) (string, error))),
		expr.Function("has", fnHas, new(func(string, string) (bool, error))),
		expr.Function("split", fnSplit, new(func(string, string) ([]string, error))),
		expr.Function("default", fnDefault, new(func(any, any) (any, error))),

		expr.AsAny(),
	}

	program, err := expr.Compile(expression, opts...)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to compile expression %q: %w", ErrCompileTag, expression, err)
	}

	return program, nil
}

func kv(pair string) (string, string) {
	parts := strings.SplitN(pair, "=", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}

	return strings.TrimSpace(parts[0]), ""
}

func fnGet(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, errors.New("get() requires exactly 2 arguments")
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, errors.New("get() first argument must be a string")
	}

	key, ok := args[1].(string)
	if !ok {
		return nil, errors.New("get() second argument must be a string")
	}

	for pair := range strings.SplitSeq(str, ",") {
		k, v := kv(pair)
		if k == key {
			return v, nil
		}
	}

	return nil, nil
}

func fnFirst(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, errors.New("first() requires exactly 1 argument")
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, errors.New("first() argument must be a string")
	}

	pairs := strings.Split(str, ",")
	if len(pairs) == 0 {
		return nil, nil
	}

	k, v := kv(pairs[0])
	if v == "" {
		return k, nil
	}

	return v, nil
}

func fnLast(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, errors.New("last() requires exactly 1 argument")
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, errors.New("last() argument must be a string")
	}

	pairs := strings.Split(str, ",")
	if len(pairs) == 0 {
		return nil, nil
	}

	k, v := kv(pairs[len(pairs)-1])
	if v == "" {
		return k, nil
	}

	return v, nil
}

func fnNth(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, errors.New("nth() requires exactly 2 arguments")
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, errors.New("nth() first argument must be a string")
	}

	index, ok := args[1].(int)
	if !ok {
		return nil, errors.New("nth() second argument must be an integer")
	}

	pairs := strings.Split(str, ",")
	if index < 0 || index >= len(pairs) {
		return nil, nil
	}

	k, v := kv(pairs[index])
	if v == "" {
		return k, nil
	}

	return v, nil
}

func fnHas(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, errors.New("has() requires exactly 2 arguments")
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, errors.New("has() first argument must be a string")
	}

	key, ok := args[1].(string)
	if !ok {
		return nil, errors.New("has() second argument must be a string")
	}

	for pair := range strings.SplitSeq(str, ",") {
		k, _ := kv(pair)
		if k == key {
			return true, nil
		}
	}

	return false, nil
}

func fnSplit(args ...any) (any, error) {
	if len(args) < 2 {
		return nil, errors.New("split() requires at least 2 arguments")
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, errors.New("split() first argument must be a string")
	}

	sep, ok := args[1].(string)
	if !ok {
		return nil, errors.New("split() second argument must be a string")
	}

	parts := strings.Split(str, sep)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return parts, nil
}

func fnDefault(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, errors.New("default() requires exactly 2 arguments")
	}

	if args[0] == nil {
		return args[1], nil
	}

	return args[0], nil
}
