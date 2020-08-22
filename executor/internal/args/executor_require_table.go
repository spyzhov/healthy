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

func (a RequireTable) Match(table [][]interface{}, null string) (err error) {
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
			if !a.same(value, table[row][col], null) {
				return fmt.Errorf("wrong value at (%d, %d)", row, col)
			}
		}
	}
	return nil
}

func (a RequireTable) same(x, y interface{}, null string) bool {
	if reflect.DeepEqual(x, y) {
		return true
	}
	if a.str(x, null) == a.str(y, null) {
		return true
	}
	return false
}

func (a RequireTable) str(x interface{}, null string) string {
	if safe.IsNil(x) {
		return null
	}
	return fmt.Sprintf("%v", x)
}
