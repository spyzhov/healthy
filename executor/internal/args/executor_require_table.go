package args

import (
	"fmt"
	"reflect"

	"github.com/spyzhov/safe"
)

type RequireTable [][]interface{}

func (a RequireTable) Validate() (err error) {
	if a == nil {
		return nil
	}
	var width int
	for row, values := range a {
		if row == 0 {
			width = len(values)
		} else if width != len(values) {
			return fmt.Errorf("row %d has wrong count of elements", row)
		}
	}
	return nil
}

func (a RequireTable) Match(table [][]interface{}) (err error) {
	if a == nil || len(a) == 0 {
		return nil
	}
	if len(table) != len(a) {
		return fmt.Errorf("wrong length")
	}
	if len(table[0]) != len(a[0]) {
		return fmt.Errorf("wrong width")
	}
	if err = RequireTable(table).Validate(); err != nil {
		return safe.Wrap(err, "examine")
	}
	for row, values := range a {
		for col, value := range values {
			if !a.same(value, table[row][col]) {
				return fmt.Errorf("wrong value at (%d, %d)", row, col)
			}
		}
	}
	return nil
}

func (a RequireTable) same(x, y interface{}) bool {
	if reflect.DeepEqual(x, y) {
		return true
	}
	if fmt.Sprintf("%v", x) == fmt.Sprintf("%v", y) {
		return true
	}
	return false
}
