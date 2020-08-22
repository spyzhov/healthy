package args

import (
	"fmt"
)

type SqlArgsRequireValues []SqlArgsRequireValue

func (a SqlArgsRequireValues) Validate() (err error) {
	for _, value := range a {
		if err = value.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (a SqlArgsRequireValues) Match(rows [][]interface{}) (err error) {
	for i, value := range a {
		if err = value.Match(rows); err != nil {
			return fmt.Errorf("value %d: %w", i, err)
		}
	}
	return nil
}
