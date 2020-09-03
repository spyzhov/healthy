package args

import "fmt"

type String string

func (a String) Validate() (err error) {
	if len(a) == 0 {
		return fmt.Errorf("string shouldn't be blank")
	}
	return nil
}

func (a String) Value() string {
	return string(a)
}
