package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type SqlArgsRequireColumns []SqlArgsRequireColumn

func (a SqlArgsRequireColumns) Validate() (err error) {
	if a == nil {
		return nil
	}
	for i, column := range a {
		if err = column.Validate(); err != nil {
			return safe.Wrap(err, fmt.Sprintf("column (%d)", i))
		}
	}
	return nil
}

func (a SqlArgsRequireColumns) Match(table Table) (err error) {
	if a == nil {
		return nil
	}
	for i, column := range a {
		if err = column.Match(table); err != nil {
			return safe.Wrap(err, fmt.Sprintf("column (%d)", i))
		}
	}
	return nil
}
