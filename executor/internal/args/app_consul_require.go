package args

import (
	"github.com/hashicorp/consul/api"
	"github.com/spyzhov/safe"
)

type AppConsulArgsRequire struct {
	Key     map[string]*RequireContent              `json:"key"`
	Service map[string]*AppConsulArgsRequireService `json:"service"`
}

func (a *AppConsulArgsRequire) Validate() (err error) {
	if a == nil {
		return nil
	}
	for key, require := range a.Key {
		if err = require.Validate(); err != nil {
			return safe.Wrap(err, "require: key: "+key)
		}
	}
	return
}

func (a *AppConsulArgsRequire) Match(client *api.Client) (err error) {
	if a == nil {
		return nil
	}
	for key, require := range a.Key {
		scope := "require: key: " + key
		value := make([]byte, 0)
		pair, _, err := client.KV().Get(key, nil)
		if err != nil {
			return safe.Wrap(err, scope+": get")
		}
		if pair != nil {
			value = pair.Value
		}
		if err = require.Match(scope, value); err != nil {
			return err
		}
	}
	for key, require := range a.Service {
		scope := "require: service: " + key
		out, _, err := client.Health().ServiceMultipleTags(key, require.Tags, false, nil)
		if err != nil {
			return safe.Wrap(err, scope+": get")
		}
		if err = require.Match(out); err != nil {
			return safe.Wrap(err, scope)
		}
	}

	return
}
