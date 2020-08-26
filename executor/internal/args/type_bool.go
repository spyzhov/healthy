package args

import (
	"fmt"
)

type Bool bool

func (a *Bool) Validate() error {
	return nil
}

func (a *Bool) Match(value bool, True, False string) error {
	if a == nil {
		return nil
	}
	if a.value() && !value {
		return fmt.Errorf("value is %s", False)
	}
	if !a.value() && value {
		return fmt.Errorf("value is %s", True)
	}

	return nil
}

func (a *Bool) value() bool {
	return bool(*a)
}
