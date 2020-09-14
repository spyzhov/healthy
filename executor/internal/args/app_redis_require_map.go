package args

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/spyzhov/safe"
)

type AppRedisArgsRequireMap struct {
	Count RequireNumeric               `json:"count"`
	All   *RequireValue                `json:"all"`
	Value AppRedisArgsRequireMapValues `json:"value"`
}

func (a *AppRedisArgsRequireMap) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.Count.Validate(); err != nil {
		return safe.Wrap(err, "count")
	}
	if err = a.All.Validate(); err != nil {
		return safe.Wrap(err, "all")
	}
	if err = a.Value.Validate(); err != nil {
		return safe.Wrap(err, "value")
	}
	return nil
}

func (a *AppRedisArgsRequireMap) Match(value interface{}) error {
	if a == nil {
		return nil
	}
	values, err := a.get(value)
	if err != nil {
		return err
	}
	if err = a.Count.MatchInt("count", len(values)); err != nil {
		return err
	}
	if a.All != nil {
		for key, val := range values {
			if err = a.All.Match(val); err != nil {
				return fmt.Errorf("all: key (%s): %w", key, err)
			}
		}
	}
	if err = a.Value.Match(values); err != nil {
		return safe.Wrap(err, "value")
	}
	return nil
}

func (a *AppRedisArgsRequireMap) get(value interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if values, err := redis.StringMap(value, nil); err == nil {
		for key, val := range values {
			result[key] = val
		}
		return result, nil
	}
	if values, err := redis.IntMap(value, nil); err == nil {
		for key, val := range values {
			result[key] = val
		}
		return result, nil
	}
	if values, err := redis.Int64Map(value, nil); err == nil {
		for key, val := range values {
			result[key] = val
		}
		return result, nil
	}
	if values, err := redis.Uint64Map(value, nil); err == nil {
		for key, val := range values {
			result[key] = val
		}
		return result, nil
	}

	return nil, fmt.Errorf("value is not a map")
}
