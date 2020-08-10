package config

import (
	"os"
	"regexp"
)

func dereference(content []byte) []byte {
	search := regexp.MustCompile(`(?m)(env\([^)]+\))`)
	return search.ReplaceAllFunc(content, func(bytes []byte) []byte {
		key := string(bytes[4 : len(bytes)-1])
		value := os.Getenv(key)
		return []byte(value)
	})
}
