package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type Slice []interface{}

func (a Slice) Validate() (err error) {
	if a == nil {
		return nil
	}
	return nil
}

func (a Slice) Match(slice Slice, null string) (err error) {
	if a == nil || len(a) == 0 {
		return nil
	}
	if len(slice) != len(a) {
		return fmt.Errorf("wrong length")
	}
	if err = slice.Validate(); err != nil {
		return safe.Wrap(err, "examine")
	}
	for i, value := range a {
		if !same(value, slice[i], null) {
			return fmt.Errorf("wrong value at (%d)", i)
		}

	}
	return nil
}
