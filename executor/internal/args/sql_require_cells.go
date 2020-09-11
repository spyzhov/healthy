package args

import (
	"fmt"
)

type SqlArgsRequireCells []SqlArgsRequireCell

func (a SqlArgsRequireCells) Validate() (err error) {
	for _, value := range a {
		if err = value.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (a SqlArgsRequireCells) Match(rows Table) (err error) {
	for i, value := range a {
		if err = value.Match(rows); err != nil {
			return fmt.Errorf("value %d: %w", i, err)
		}
	}
	return nil
}
