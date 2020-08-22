package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type SqlArgsRequireValue struct {
	Row     Uint             `json:"row"`
	Column  Uint             `json:"column"`
	Numeric *RequireNumeric  `json:"numeric"`
	Text    *RequireMatch    `json:"text"`
	JSON    *RequireJSONPath `json:"json"`
}

func (a *SqlArgsRequireValue) Validate() (err error) {
	if err = a.Row.Validate(); err != nil {
		return safe.Wrap(err, "value: row")
	}
	if err = a.Column.Validate(); err != nil {
		return safe.Wrap(err, "value: column")
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

func (a *SqlArgsRequireValue) Match(rows [][]interface{}) (err error) {
	if len(rows) <= int(a.Row) {
		return fmt.Errorf("row: not found")
	}
	if len(rows[a.Row]) <= int(a.Column) {
		return fmt.Errorf("column: not found")
	}
	str := fmt.Sprintf("%v", rows[a.Row][a.Column])
	if safe.IsNil(rows[a.Row][a.Column]) {
		str = "NULL"
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
