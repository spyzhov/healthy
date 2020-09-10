package executor

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/spyzhov/healthy/config"
	. "github.com/spyzhov/healthy/executor/internal/args"
	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
)

type RepeatArgs struct {
	Count   int               `json:"count"`
	Delay   Duration          `json:"delay"`
	Require RepeatArgsRequire `json:"require"`
}

func (e *Executor) Repeat(args *RepeatArgs, cmd *config.Step) (step.Function, error) {
	if err := args.Validate(); err != nil {
		return nil, safe.Wrap(err, "repeat")
	}

	if args.Require.Success == 0 {
		args.Require.Success = Uint(args.Count)
	}
	length := int(math.Log10(float64(args.Count))) + 1
	format := fmt.Sprintf("%%0%dd/%%0%dd:[%%-7s]:%%s", length, length)

	fn, err := Get(e, cmd.Type, cmd.Args)
	if err != nil {
		return nil, safe.Wrap(err, "repeat")
	}
	test := step.NewStep("repeat", fn, cmd.Vars.Masked())

	return func() (*step.Result, error) {
		success := Uint(0)
		messages := make([]string, 0, args.Count)

		for i := 0; i < args.Count; i++ {
			if i != 0 && args.Delay.Duration != 0 {
				time.Sleep(args.Delay.Duration)
			}
			res := test.Call()
			messages = append(messages, fmt.Sprintf(format, i+1, args.Count, res.Status, res.Message))
			if res.Status == step.Success {
				success++
			}
		}

		message := strings.Join(messages, "\n")
		if success >= args.Require.Success {
			return step.NewResultSuccess(message), nil
		}
		if success >= args.Require.Warning {
			return step.NewResultWarning(message), nil
		}

		return step.NewResultError(message), nil
	}, nil
}

func (a *RepeatArgs) Validate() (err error) {
	if a.Count <= 0 {
		return fmt.Errorf("repeat: count: should be greater than zero")
	}
	if err = a.Delay.Validate(); err != nil {
		return safe.Wrap(err, "repeat: delay")
	}
	if err = a.Require.Validate(); err != nil {
		return safe.Wrap(err, "repeat: require")
	}
	if int(a.Require.Success) > a.Count {
		return fmt.Errorf("repeat: require: success: should be lesser or equal than %d", a.Count)
	}
	if int(a.Require.Warning) > a.Count {
		return fmt.Errorf("repeat: require: warning: should be lesser or equal than %d", a.Count)
	}
	return nil
}
