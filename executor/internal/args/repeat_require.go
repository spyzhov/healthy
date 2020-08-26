package args

import (
	"github.com/spyzhov/safe"
)

type RepeatArgsRequire struct {
	Success Uint `json:"success"`
	Warning Uint `json:"warning"`
}

func (a *RepeatArgsRequire) Validate() (err error) {
	if err = a.Success.Validate(); err != nil {
		return safe.Wrap(err, "success")
	}
	if err = a.Warning.Validate(); err != nil {
		return safe.Wrap(err, "warning")
	}
	return nil
}
