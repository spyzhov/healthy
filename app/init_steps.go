package app

import (
	"fmt"

	"github.com/spyzhov/healthy/executor"
	"github.com/spyzhov/healthy/step"
)

func (app *Application) setSteps() (err error) {
	for _, group := range app.StepConfig.Groups {
		for _, value := range group.Validate {
			var fn step.Function
			fn, err = executor.Get(app.Executor, value.Type, value.Args)
			if err != nil {
				return fmt.Errorf("can't get validate function: %w", err)
			}
			app.StepGroups.Get(group.Name).Add(value.Name, fn)
		}
	}
	// Example:
	//app.StepGroups.Get("Custom steps").
	//	Add("Random: <randomly>", func() (*step.Result, error) {
	//		res := rand.Float64()
	//		if res < .25 {
	//			panic(fmt.Sprintf("got %0.2f < .25 => PANIC", res))
	//		}
	//		if res < .50 {
	//			return step.NewResultError(fmt.Sprintf("got %0.2f < .50 => ERROR", res)), nil
	//		}
	//		if res < .75 {
	//			return step.NewResultWarning(fmt.Sprintf("got %0.2f < .75 => WARNING", res)), nil
	//		}
	//		return step.NewResultSuccess(fmt.Sprintf("got %0.2f < 1.0 => SUCCESS", res)), nil
	//	})
	return nil
}
