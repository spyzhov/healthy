package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type Row struct {
	Row Uint `json:"row"`
}

func (a *Row) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.Row.Validate(); err != nil {
		return safe.Wrap(err, "row")
	}
	return nil
}

func (a *Row) Match(table Table) (err error) {
	if a == nil {
		return nil
	}
	if len(table) <= a.Row.value() {
		return fmt.Errorf("row %d not found", a.Row)
	}
	return nil
}

func (a *Row) value() int {
	if a == nil {
		return -1
	}
	return a.Row.value()
}
