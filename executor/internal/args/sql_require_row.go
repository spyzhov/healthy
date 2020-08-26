package args

import (
	"github.com/spyzhov/safe"
)

type SqlArgsRequireRow struct {
	Row
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
	return nil
}

func (a *SqlArgsRequireRow) Match(table Table) (err error) {
	if a == nil {
		return nil
	}
	if err = a.Row.Match(table); err != nil {
		return err
	}
	if err = a.Value.Match(table[a.Row.value()], "NULL"); err != nil {
		return err
	}
	return nil
}
