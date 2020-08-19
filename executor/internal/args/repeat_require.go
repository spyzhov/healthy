package args

import "fmt"

type RepeatArgsRequire struct {
	Success int `json:"success"`
	Warning int `json:"warning"`
}

func (a *RepeatArgsRequire) Validate() (err error) {
	if a.Success < 0 {
		return fmt.Errorf("success: should be greater or equal than zero")
	}
	if a.Warning < 0 {
		return fmt.Errorf("warning: should be greater or equal than zero")
	}
	return nil
}
