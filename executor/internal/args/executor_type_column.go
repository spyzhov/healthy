package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type Column struct {
	Column Uint `json:"column"`
}

func (a *Column) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.Column.Validate(); err != nil {
		return safe.Wrap(err, "column")
	}
	return nil
}

func (a *Column) Match(table Table) (err error) {
	if a == nil {
		return nil
	}
	if len(table) == 0 {
		return fmt.Errorf("table is blank")
	}
	if len(table[0]) <= a.Column.value() {
		return fmt.Errorf("column %d not found", a.Column)
	}
	return nil
}

func (a *Column) value() int {
	if a == nil {
		return -1
	}
	return a.Column.value()
}
