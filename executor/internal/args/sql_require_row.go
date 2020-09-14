package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type SqlArgsRequireRow struct {
	Row
	RequireValue
	Value Slice `json:"value"`
}

func (a *SqlArgsRequireRow) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.Row.Validate(); err != nil {
		return safe.Wrap(err, "row")
	}
	if err = a.Value.Validate(); err != nil {
		return safe.Wrap(err, "value")
	}
	if err = a.RequireValue.Validate(); err != nil {
		return err
	}
	return nil
}

func (a *SqlArgsRequireRow) Match(table Table) (err error) {
	if a == nil {
		return nil
	}
	if err = a.Row.Match(table); err != nil {
		return err
	}
	row := table.Row(a.Row.value())
	if err = a.Value.Match(row, "NULL"); err != nil {
		return err
	}
	for i, value := range row {
		if err = a.RequireValue.Match(value); err != nil {
			return fmt.Errorf("column (%d): %w", i, err)
		}
	}
	return nil
}
