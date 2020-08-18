package internal

import (
	"fmt"
	"regexp"
)

type RequireFieldMatch []string

func (match RequireFieldMatch) Validate() error {
	for _, value := range match {
		if _, err := regexp.Compile(value); err != nil {
			return fmt.Errorf("MATCH: regexp compile(`%s`) error: %w", value, err)
		}
	}
	return nil
}

func (match RequireFieldMatch) Match(name string, input []byte) error {
	for _, value := range match {
		if !regexp.MustCompile(value).Match(input) {
			return fmt.Errorf("%s: value is not match `%s`", name, value)
		}
	}
	return nil
}

func (match RequireFieldMatch) MatchStrings(name string, input []string) error {
	for _, value := range match {
		if !match.any(regexp.MustCompile(value), input) {
			return fmt.Errorf("%s: value is not match `%s`", name, value)
		}
	}
	return nil
}

func (RequireFieldMatch) any(re *regexp.Regexp, values []string) bool {
	for _, value := range values {
		if re.MatchString(value) {
			return true
		}
	}
	return false
}
