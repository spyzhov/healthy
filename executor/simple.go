package executor

import (
	"fmt"
	"time"

	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
)

type SimpleArgs struct {
	Sleep   Duration    `json:"sleep"`
	Status  step.Status `json:"status"`
	Message string      `json:"message"`
}

// Simple will just return value from args
func (e *Executor) Simple(args *SimpleArgs) (step.Function, error) {
	if safe.IsNil(args) {
		return nil, fmt.Errorf("arguments should be set: `status` and `message`")
	}
	if err := args.Validate(); err != nil {
		return nil, safe.Wrap(err, "simple")
	}

	return func() (*step.Result, error) {
		if args.Sleep.Duration != 0 {
			time.Sleep(args.Sleep.Duration)
		}
		message := args.Message
		if message == "" {
			message = string(args.Status)
		}
		return step.NewResult(args.Status, message), nil
	}, nil
}

func (a *SimpleArgs) Validate() (err error) {
	if safe.IsNil(a) {
		return fmt.Errorf("arguments should be set: `status` and `message`")
	}
	if err = a.Sleep.Validate(); err != nil {
		return err
	}
	if err = a.Status.Validate(); err != nil {
		return err
	}
	return nil
}
