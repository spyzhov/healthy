package executor

import (
	"fmt"
	"reflect"
	"sort"
)

type HttpArgsRequireHeader struct {
	Exists    []string                        `json:"exists"`
	NotExists []string                        `json:"not_exists"`
	Regexp    map[string]RequireFieldMatch    `json:"match"`
	NotRegexp map[string]RequireFieldMatchNot `json:"not_match"`
	Eq        map[string][]string             `json:"eq"`
}

func (a *HttpArgsRequireHeader) Validate() (err error) {
	for _, match := range a.Regexp {
		if err = match.Validate(); err != nil {
			return err
		}
	}
	for _, match := range a.NotRegexp {
		if err = match.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (a *HttpArgsRequireHeader) Match(header map[string][]string) (err error) {
	var (
		ok     bool
		actual []string
	)
	for _, key := range a.Exists {
		if _, ok = header[key]; !ok {
			return fmt.Errorf("header: not EXISTS `%s`", key)
		}
	}
	for _, key := range a.NotExists {
		if _, ok = header[key]; ok {
			return fmt.Errorf("header: EXISTS `%s`", key)
		}
	}
	for key, value := range a.Eq {
		if actual, ok = header[key]; !ok {
			return fmt.Errorf("header: not EQ `%s`. Not exists", key)
		} else {
			sort.Strings(actual)
			sort.Strings(value)
			if !reflect.DeepEqual(actual, value) {
				return fmt.Errorf("header: not EQ `%s`. Expected: %v; Actual: %v", key, value, actual)
			}
		}
	}
	for key, match := range a.Regexp {
		if actual, ok = header[key]; !ok {
			return fmt.Errorf("header: MATCH `%s`. Not exists", key)
		} else {
			if err = match.MatchStrings("header", actual); err != nil {
				return err
			}
		}
	}
	for key, match := range a.NotRegexp {
		if actual, ok = header[key]; !ok {
			return fmt.Errorf("header: NOT_MATCH `%s`. Not exists", key)
		} else {
			if err = match.MatchStrings("header", actual); err != nil {
				return err
			}
		}
	}
	return nil
}
