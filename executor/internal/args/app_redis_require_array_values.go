package args

import (
	"fmt"
)

type AppRedisArgsRequireArrayValues []*AppRedisArgsRequireArrayValue

func (a AppRedisArgsRequireArrayValues) Validate() (err error) {
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

func (a AppRedisArgsRequireArrayValues) Match(values []interface{}) (err error) {
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
