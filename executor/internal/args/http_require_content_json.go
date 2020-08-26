package args

import (
	"fmt"

	"github.com/spyzhov/ajson"
	"github.com/spyzhov/safe"
)

type HttpArgsRequireContentJSON struct {
	JSONPath string `json:"jsonpath"`
	RequireXPath
	RequireJSONSchema
}

func (a HttpArgsRequireContentJSON) Validate() (err error) {
	if err = a.RequireXPath.Validate(); err != nil {
		return err
	}
	if err = a.RequireJSONSchema.Validate(); err != nil {
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
		if err = a.RequireMatch.Match(name, result); err != nil {
			return err
		}
		if err = a.RequireJSONSchema.Match(result); err != nil {
			return safe.Wrap(err, name)
		}
	}
	if a.XPath != "" {
		if err = a.RequireXPath.Match(name, HttpArgsRequireContentTypeJSON, content); err != nil {
			return err
		}
		result, err := a.RequireXPath.json(content)
		if err != nil {
			return safe.Wrap(err, name)
		}
		if err = a.RequireJSONSchema.Match(result); err != nil {
			return safe.Wrap(err, name)
		}
	}
	return nil
}
