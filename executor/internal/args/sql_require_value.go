package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type SqlArgsRequireValue struct {
	IsNull  *Bool           `json:"is_null"`
	Numeric *RequireNumeric `json:"numeric"`
	Content *RequireContent `json:"content"`
}

func (a *SqlArgsRequireValue) Validate() (err error) {
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
	return nil
}

func (a *SqlArgsRequireValue) Match(value interface{}) (err error) {
	if a == nil {
		return nil
	}
	isNil := safe.IsNil(value)
	str := fmt.Sprintf("%v", value)
	if isNil {
		str = "NULL"
	}
	if err = a.IsNull.Match(isNil, "NULL", "NOT NULL"); err != nil {
		return err
	}
	if err = a.Numeric.MatchString("numeric", str); err != nil {
		return err
	}
	if err = a.Content.Match("content", []byte(str)); err != nil {
		return err
	}
	return nil
}
