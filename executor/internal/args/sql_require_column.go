package args

import (
	"github.com/spyzhov/safe"
)

type SqlArgsRequireColumn struct {
	Column
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
	return nil
}

func (a *SqlArgsRequireColumn) Match(table Table) (err error) {
	if a == nil {
		return nil
	}
	if err = a.Column.Match(table); err != nil {
		return err
	}
	if err = a.Value.Match(table[a.Column.value()], "NULL"); err != nil {
		return err
	}
	return nil
}
