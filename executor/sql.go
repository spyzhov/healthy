package executor

import (
	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
)

type SqlArgs struct {
}

func (e *Executor) Sql(args *SqlArgs) (step.Function, error) {
	if err := args.Validate(); err != nil {
		return nil, safe.Wrap(err, "http")
	}

	return func() (*step.Result, error) {
		return step.NewResultWarning("not implemented"), nil
	}, nil
}

func (a *SqlArgs) Validate() (err error) {
	panic("implement me")
}
