package args

import (
	"github.com/spyzhov/safe"
)

type Cell struct {
	Row
	Column
}

func (a *Cell) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.Row.Validate(); err != nil {
		return safe.Wrap(err, "row")
	}
	if err = a.Column.Validate(); err != nil {
		return safe.Wrap(err, "column")
	}
	return nil
}

func (a *Cell) Match(table Table) (err error) {
	if a == nil {
		return nil
	}
	if err = a.Row.Match(table); err != nil {
		return err
	}
	if err = a.Column.Match(table); err != nil {
		return err
	}
	return nil
}

func (a *Cell) get(table Table) interface{} {
	return table[a.Row.value()][a.Column.value()]
}
