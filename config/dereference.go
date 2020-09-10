package config

import (
	"os"
)

//func dereference(content []byte) []byte {
//	search := regexp.MustCompile(`(?m)(env\([^)]+\))`)
//	return search.ReplaceAllFunc(content, func(bytes []byte) []byte {
//		key := string(bytes[4 : len(bytes)-1])
//		value := os.Getenv(key)
//		return []byte(value)
//	})
//}

func expand(content string, vars map[string]string) string {
	return os.Expand(content, func(key string) string {
		if value, ok := vars[key]; ok {
			return value
		}
		return os.Getenv(key)
	})
}

func dereference(value interface{}, fn func(string) string) interface{} {
	switch result := value.(type) {
	case map[string]interface{}:
		for key, iVal := range result {
			result[key] = dereference(iVal, fn)
		}
		return result
	case []interface{}:
		for i, iVal := range result {
			result[i] = dereference(iVal, fn)
		}
		return result
	case string:
		return fn(result)
	default:
		return result
	}
}
