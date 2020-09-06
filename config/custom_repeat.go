package config

import "encoding/json"

func stepDereference(value interface{}, vars map[string]string) (result interface{}, err error) {
	var (
		data  []byte
		raw   = new(step)
		local = make(map[string]string, len(vars))
	)
	for key, val := range vars {
		local[key] = val
	}
	if data, err = json.Marshal(value); err != nil {
		return nil, err
	} else if err = raw.unmarshal(data, local); err != nil {
		return nil, err
	} else if err = read(raw, &result); err != nil {
		return nil, err
	}
	return result, nil
}
