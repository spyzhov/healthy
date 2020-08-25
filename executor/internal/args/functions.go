package args

import (
	"fmt"
	"reflect"

	"github.com/spyzhov/safe"
)

func same(x, y interface{}, null string) bool {
	if reflect.DeepEqual(x, y) {
		return true
	}
	if str(x, null) == str(y, null) {
		return true
	}
	return false
}

func str(x interface{}, null string) string {
	if safe.IsNil(x) {
		return null
	}
	return fmt.Sprintf("%v", x)
}
