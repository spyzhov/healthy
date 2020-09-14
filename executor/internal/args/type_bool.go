package args

import (
	"fmt"
	"strconv"
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

func (a *Bool) MatchInterface(value interface{}, True, False string) error {
	if a == nil {
		return nil
	}
	if val, ok := value.(bool); ok {
		return a.Match(val, True, False)
	}
	var str string
	if val, ok := value.([]byte); ok {
		str = string(val)
	} else {
		str = fmt.Sprintf("%v", value)
	}
	if str == "1" {
		return a.Match(true, True, False)
	} else if str == "0" {
		return a.Match(false, True, False)
	}
	if val, err := strconv.ParseBool(str); err != nil {
		return a.Match(val, True, False)
	}

	return fmt.Errorf("value is not boolean")
}

func (a *Bool) value() bool {
	return bool(*a)
}
