package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type RequireValue struct {
	IsNull  *Bool           `json:"is_null"`
	Numeric *RequireNumeric `json:"numeric"`
	Content *RequireContent `json:"content"`
	Bool    *Bool           `json:"bool"`
}

func (a *RequireValue) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.IsNull.Validate(); err != nil {
		return safe.Wrap(err, "is_null")
	}
	if err = a.Numeric.Validate(); err != nil {
		return safe.Wrap(err, "numeric")
	}
	if err = a.Content.Validate(); err != nil {
		return safe.Wrap(err, "content")
	}
	if err = a.Bool.Validate(); err != nil {
		return safe.Wrap(err, "bool")
	}
	return nil
}

func (a *RequireValue) Match(value interface{}) (err error) {
	if a == nil {
		return nil
	}
	isNil := safe.IsNil(value)
	str := fmt.Sprintf("%v", value)
	if isNil {
		str = "NULL"
	}
	if err = a.IsNull.Match(isNil, "NULL", "NOT NULL"); err != nil {
		return safe.Wrap(err, "is_null")
	}
	if err = a.Numeric.MatchString("numeric", str); err != nil {
		return err
	}
	if err = a.Content.Match("content", []byte(str)); err != nil {
		return err
	}
	if err = a.Bool.MatchInterface(value, "True", "False"); err != nil {
		return safe.Wrap(err, "bool")
	}
	return nil
}
