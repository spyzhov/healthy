package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type SqlArgsRequireTable struct {
	Value Table `json:"value"`
	RequireValue
}

func (a *SqlArgsRequireTable) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.Value.Validate(); err != nil {
		return safe.Wrap(err, "value")
	}
	if err = a.RequireValue.Validate(); err != nil {
		return err
	}
	return nil
}

func (a *SqlArgsRequireTable) Match(table Table) (err error) {
	if a == nil {
		return nil
	}
	if err = a.Value.Match(table, "NULL"); err != nil {
		return err
	}
	for row, values := range table {
		for column, value := range values {
			if err = a.RequireValue.Match(value); err != nil {
				return fmt.Errorf("cell at (%d, %d): %w", row, column, err)
			}
		}
	}
	return nil
}
