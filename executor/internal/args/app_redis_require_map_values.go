package args

import (
	"fmt"
)

type AppRedisArgsRequireMapValues []*AppRedisArgsRequireMapValue

func (a AppRedisArgsRequireMapValues) Validate() (err error) {
	if a == nil {
		return nil
	}
	for i, value := range a {
		if err = value.Validate(); err != nil {
			return fmt.Errorf("row (%d): %w", i, err)
		}
	}
	return nil
}

func (a AppRedisArgsRequireMapValues) Match(values map[string]interface{}) (err error) {
	if a == nil {
		return nil
	}
	for i, value := range a {
		if err = value.Match(values); err != nil {
			return fmt.Errorf("row (%d): %w", i, err)
		}
	}
	return nil
}
