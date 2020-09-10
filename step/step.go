package step

import (
	"fmt"
	"runtime/debug"
	"strings"

	"go.uber.org/zap"
)

const mask = "***mask***"

type Step struct {
	Name     string
	Func     Function
	Replacer *strings.Replacer
}

type Function func() (*Result, error)

func NewStep(name string, fn Function, masked []string) *Step {
	replace := make([]string, 0, len(masked)*2)
	for _, str := range masked {
		replace = append(replace, str, mask)
	}

	return &Step{
		Name:     name,
		Func:     fn,
		Replacer: strings.NewReplacer(replace...),
	}
}

func (step *Step) Call() (res *Result) {
	defer func() {
		res.Message = step.Replacer.Replace(res.Message)
	}()
	defer func() {
		if res == nil {
			res = NewResultError("Validation result is not set")
		}
	}()
	defer func() {
		if rec := recover(); rec != nil {
			zap.L().Error("recover panic", zap.Any("recover", rec), zap.ByteString("stack", debug.Stack()))
			res = NewResultError(fmt.Sprintf("recover: %v", rec))
		}
	}()
	var err error

	res, err = step.Func()
	if err != nil {
		return NewResultError("error: " + err.Error())
	}
	return res
}
