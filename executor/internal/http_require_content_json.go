package internal

import (
	"fmt"

	"github.com/spyzhov/ajson"
)

type HttpArgsRequireContentJSON struct {
	JSONPath string `json:"jsonpath"`
	RequireXPath
}

func (a HttpArgsRequireContentJSON) Validate() (err error) {
	if err = a.RequireXPath.Validate(); err != nil {
		return err
	}
	return nil
}

func (a HttpArgsRequireContentJSON) Match(name string, content []byte) (err error) {
	if a.JSONPath != "" {
		nodes, err := ajson.JSONPath(content, a.JSONPath)
		if err != nil {
			return fmt.Errorf("%s: JSONPath(`%s`) compile error: %w", name, a.JSONPath, err)
		}
		result, err := ajson.Marshal(ajson.ArrayNode("", nodes))
		if err != nil {
			return fmt.Errorf("%s: JSONPath(`%s`) marshal error: %w", name, a.JSONPath, err)
		}
		return a.RequireMatch.Match(name, result)
	}
	if err = a.RequireXPath.Match(name, HttpArgsRequireContentTypeJSON, content); err != nil {
		return err
	}
	return nil
}
