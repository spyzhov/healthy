package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type AppRedisArgsRequireArrayValue struct {
	Index *Uint `json:"index"`
	RequireValue
}

func (a *AppRedisArgsRequireArrayValue) Validate() (err error) {
	if a == nil {
		return nil
	}
	if a.Index == nil {
		return fmt.Errorf("index: not set")
	}
	if err = a.Index.Validate(); err != nil {
		return safe.Wrap(err, "index")
	}
	if err = a.RequireValue.Validate(); err != nil {
		return err
	}
	return nil
}

func (a *AppRedisArgsRequireArrayValue) Match(value []interface{}) (err error) {
	if a == nil {
		return nil
	}
	if len(value) <= a.Index.value() {
		return fmt.Errorf("index: out of range")
	}
	if err = a.RequireValue.Match(value[a.Index.value()]); err != nil {
		return fmt.Errorf("index (%d): %w", a.Index.value(), err)
	}
	return nil
}
