package args

import (
	"fmt"

	"github.com/spyzhov/ajson"
)

type RequireJSONPath struct {
	JSONPath string `json:"jsonpath"`
	RequireMatch
}

func (a *RequireJSONPath) Validate() error {
	return a.RequireMatch.Validate()
}

func (a *RequireJSONPath) Match(name string, content []byte) error {
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
	return nil
}
