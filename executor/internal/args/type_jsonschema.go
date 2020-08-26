package args

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spyzhov/safe"
	"github.com/xeipuuv/gojsonschema"
)

type JSONSchema string

func (a *JSONSchema) Validate() error {
	if a == nil {
		return nil
	}
	if !json.Valid([]byte(a.value())) {
		return fmt.Errorf("JSON is invalid")
	}
	return nil
}

func (a *JSONSchema) Match(value []byte) error {
	if a == nil {
		return nil
	}
	return a.match(gojsonschema.NewBytesLoader(value))
}

func (a *JSONSchema) MatchString(value string) error {
	if a == nil {
		return nil
	}
	return a.match(gojsonschema.NewStringLoader(value))
}

func (a *JSONSchema) match(documentLoader gojsonschema.JSONLoader) error {
	schemaLoader := gojsonschema.NewStringLoader(a.value())

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return safe.Wrap(err, "JSONSchema")
	}

	if !result.Valid() {
		messages := make([]string, 0)
		for _, desc := range result.Errors() {
			messages = append(messages, desc.String())
		}
		return fmt.Errorf("JSONSchema: %s", strings.Join(messages, "\n"))
	}

	return nil
}

func (a *JSONSchema) value() string {
	if a == nil {
		return ""
	}
	return string(*a)
}
