package config

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spyzhov/safe"
	"gopkg.in/yaml.v2"
)

const (
	version1 int = 1
)

func NewConfig(content []byte) (*Config, error) {
	content = dereference(content)
	values := make(map[interface{}]interface{})
	err := yaml.Unmarshal(content, &values)
	version, ok := values["version"]
	if err != nil || !ok {
		return nil, fmt.Errorf("cannot read version of config: %w", err)
	}

	switch intVersion(version) {
	case version1:
		config := new(Config)
		return config, safe.Wrap(read(values, config), "cannot read config")
	default:
		return nil, fmt.Errorf("unknown config version: %v", version)
	}
}

func intVersion(version interface{}) int {
	switch version.(type) {
	case int:
		return version.(int)
	case int8:
		return int(version.(int8))
	case int16:
		return int(version.(int16))
	case int32:
		return int(version.(int32))
	case int64:
		return int(version.(int64))
	default:
		return -1
	}
}

func read(values interface{}, config interface{}) (err error) {
	buf := bytes.NewBuffer(make([]byte, 0, 64))
	if err = json.NewEncoder(buf).Encode(normalize(values)); err != nil {
		return err
	}
	if err = json.NewDecoder(buf).Decode(&config); err != nil {
		return err
	}
	return nil
}

func normalize(value interface{}) interface{} {
	if mmap, ok := value.(map[interface{}]interface{}); ok {
		result := make(map[string]interface{}, len(mmap))
		for key, iVal := range mmap {
			result[fmt.Sprint(key)] = normalize(iVal)
		}
		return result
	} else if mmap, ok := value.([]interface{}); ok {
		result := make([]interface{}, len(mmap))
		for i, iVal := range mmap {
			result[i] = normalize(iVal)
		}
		return result
	} else {
		return value
	}
}
