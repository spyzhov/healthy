package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type AppRedisArgsRequireMapValue struct {
	Key   *String `json:"key"`
	Exist *Bool   `json:"exist"`
	RequireValue
}

func (a *AppRedisArgsRequireMapValue) Validate() (err error) {
	if a == nil {
		return nil
	}
	if a.Key == nil {
		return fmt.Errorf("index: not set")
	}
	if err = a.Key.Validate(); err != nil {
		return safe.Wrap(err, "key")
	}
	if err = a.Exist.Validate(); err != nil {
		return safe.Wrap(err, "exist")
	}
	if err = a.RequireValue.Validate(); err != nil {
		return err
	}
	return nil
}

func (a *AppRedisArgsRequireMapValue) Match(values map[string]interface{}) (err error) {
	if a == nil {
		return nil
	}
	value, ok := values[a.Key.Value()]
	if err = a.Exist.Match(ok, "Exist", "Not exist"); err != nil {
		return fmt.Errorf("key (%s): exist: %w", a.Key.Value(), err)
	}
	if err = a.RequireValue.Match(value); err != nil {
		return fmt.Errorf("key (%s): %w", a.Key.Value(), err)
	}
	return nil
}
