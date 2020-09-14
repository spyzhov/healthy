package args

import "github.com/spyzhov/safe"

type AppRedisArgsRequire struct {
	Value *RequireValue             `json:"value"`
	Array *AppRedisArgsRequireArray `json:"array"`
	Map   *AppRedisArgsRequireMap   `json:"map"`
}

func (a *AppRedisArgsRequire) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.Value.Validate(); err != nil {
		return safe.Wrap(err, "value")
	}
	if err = a.Array.Validate(); err != nil {
		return safe.Wrap(err, "array")
	}
	if err = a.Map.Validate(); err != nil {
		return safe.Wrap(err, "map")
	}

	return nil
}

func (a *AppRedisArgsRequire) Match(value interface{}) (err error) {
	if a == nil {
		return nil
	}
	if err = a.Value.Match(value); err != nil {
		return safe.Wrap(err, "value")
	}
	if err = a.Array.Match(value); err != nil {
		return safe.Wrap(err, "array")
	}
	if err = a.Map.Match(value); err != nil {
		return safe.Wrap(err, "map")
	}

	return nil
}
