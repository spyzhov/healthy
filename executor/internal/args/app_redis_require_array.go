package args

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/spyzhov/safe"
)

type AppRedisArgsRequireArray struct {
	Count RequireNumeric                 `json:"count"`
	All   *RequireValue                  `json:"all"`
	Value AppRedisArgsRequireArrayValues `json:"value"`
}

func (a *AppRedisArgsRequireArray) Validate() (err error) {
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

func (a *AppRedisArgsRequireArray) Match(value interface{}) error {
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
		for i, val := range values {
			if err = a.All.Match(val); err != nil {
				return fmt.Errorf("all: row (%d): %w", i, err)
			}
		}
	}
	if err = a.Value.Match(values); err != nil {
		return safe.Wrap(err, "value")
	}
	return nil
}

func (a *AppRedisArgsRequireArray) get(value interface{}) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if values, err := redis.Strings(value, nil); err == nil {
		for _, val := range values {
			result = append(result, val)
		}
		return result, nil
	}
	if values, err := redis.Ints(value, nil); err == nil {
		for _, val := range values {
			result = append(result, val)
		}
		return result, nil
	}
	if values, err := redis.Int64s(value, nil); err == nil {
		for _, val := range values {
			result = append(result, val)
		}
		return result, nil
	}
	if values, err := redis.Uint64s(value, nil); err == nil {
		for _, val := range values {
			result = append(result, val)
		}
		return result, nil
	}
	if values, err := redis.Float64s(value, nil); err == nil {
		for _, val := range values {
			result = append(result, val)
		}
		return result, nil
	}
	if values, err := redis.ByteSlices(value, nil); err == nil {
		for _, val := range values {
			result = append(result, string(val))
		}
		return result, nil
	}
	if values, err := redis.Values(value, nil); err == nil {
		result = append(result, values...)
		return result, nil
	}

	return nil, fmt.Errorf("value is not an array")
}
