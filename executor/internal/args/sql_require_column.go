package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type SqlArgsRequireColumn struct {
	Column
	RequireValue
	Value Slice `json:"value"`
}

func (a *SqlArgsRequireColumn) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.Column.Validate(); err != nil {
		return safe.Wrap(err, "column")
	}
	if err = a.Value.Validate(); err != nil {
		return safe.Wrap(err, "value")
	}
	if err = a.RequireValue.Validate(); err != nil {
		return err
	}
	return nil
}

func (a *SqlArgsRequireColumn) Match(table Table) (err error) {
	if a == nil {
		return nil
	}
	if err = a.Column.Match(table); err != nil {
		return err
	}
	column := table.Column(a.Column.value())
	if err = a.Value.Match(column, "NULL"); err != nil {
		return err
	}
	for i, value := range column {
		if err = a.RequireValue.Match(value); err != nil {
			return fmt.Errorf("row (%d): %w", i, err)
		}
	}
	return nil
}
