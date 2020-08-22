package args

import "fmt"

type Uint int

func (a Uint) Validate() (err error) {
	if a < 0 {
		return fmt.Errorf("should be greater or equal than zero")
	}
	return nil
}
