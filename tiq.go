package tiq

func Parse[Schema any](field *Field) (*Schema, error) {
	tags, err := field.Tags()
	if err != nil {
		return nil, err
	}

	return parseTags[Schema](tags)
}

func Get(value any, field, tag string) (string, bool) {
	inspector, err := Inspect(value)
	if err != nil {
		return "", false
	}

	f, ok := inspector.Field(field)
	if !ok {
		return "", false
	}

	return f.Tag(tag)
}

func Set(value any, field string, newValue any) error {
	inspector, err := Inspect(value)
	if err != nil {
		return err
	}

	f, ok := inspector.Field(field)
	if !ok {
		return ErrFieldNotFound
	}

	return f.Set(newValue)
}
