package args

import (
	"fmt"

	"github.com/spyzhov/ajson"
	"github.com/spyzhov/safe"
)

type RequireJSONPath struct {
	JSONPath string `json:"jsonpath"`
	RequireMatch
	RequireJSONSchema
}

func (a *RequireJSONPath) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.RequireMatch.Validate(); err != nil {
		return err
	}
	if err = a.RequireJSONSchema.Validate(); err != nil {
		return err
	}
	return nil
}

func (a *RequireJSONPath) Match(name string, content []byte) error {
	if a == nil {
		return nil
	}
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
	return nil
}
