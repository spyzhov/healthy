package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type SqlArgsRequireRows []SqlArgsRequireRow

func (a SqlArgsRequireRows) Validate() (err error) {
	if a == nil {
		return nil
	}
	for i, row := range a {
		if err = row.Validate(); err != nil {
			return safe.Wrap(err, fmt.Sprintf("row (%d)", i))
		}
	}
	return nil
}

func (a SqlArgsRequireRows) Match(table Table) (err error) {
	if a == nil {
		return nil
	}
	for i, row := range a {
		if err = row.Match(table); err != nil {
			return safe.Wrap(err, fmt.Sprintf("row (%d)", i))
		}
	}
	return nil
}
