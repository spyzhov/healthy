package args

type RequireJSONSchema struct {
	JSONSchema *JSONSchema `json:"jsonschema"`
}

func (a *RequireJSONSchema) Validate() error {
	if a == nil {
		return nil
	}
	return a.JSONSchema.Validate()
}

func (a *RequireJSONSchema) Match(content []byte) (err error) {
	if a == nil {
		return nil
	}
	if err = a.JSONSchema.Match(content); err != nil {
		return err
	}
	return nil
}
