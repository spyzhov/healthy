package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type SqlArgsRequireValue struct {
	Cell
	IsNull  *Bool            `json:"is_null"`
	Numeric *RequireNumeric  `json:"numeric"`
	Text    *RequireMatch    `json:"text"`
	JSON    *RequireJSONPath `json:"json"`
}

func (a *SqlArgsRequireValue) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.Cell.Validate(); err != nil {
		return safe.Wrap(err, "value")
	}
	if err = a.IsNull.Validate(); err != nil {
		return safe.Wrap(err, "is_null")
	}
	if err = a.Numeric.Validate(); err != nil {
		return safe.Wrap(err, "numeric")
	}
	if err = a.Text.Validate(); err != nil {
		return safe.Wrap(err, "text")
	}
	if err = a.JSON.Validate(); err != nil {
		return safe.Wrap(err, "json")
	}
	return nil
}

func (a *SqlArgsRequireValue) Match(rows Table) (err error) {
	if a == nil {
		return nil
	}
	if err = a.Cell.Match(rows); err != nil {
		return err
	}
	isNil := safe.IsNil(a.Cell.get(rows))
	str := fmt.Sprintf("%v", a.Cell.get(rows))
	if isNil {
		str = "NULL"
	}
	if err = a.IsNull.Match(isNil, "NULL", "NOT NULL"); err != nil {
		return err
	}
	if err = a.Numeric.MatchString("numeric", str); err != nil {
		return err
	}
	if err = a.Text.Match("text", []byte(str)); err != nil {
		return err
	}
	if err = a.JSON.Match("json", []byte(str)); err != nil {
		return err
	}
	return nil
}
