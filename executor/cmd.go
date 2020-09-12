package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	. "github.com/spyzhov/healthy/executor/internal/args"
	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
)

type CmdArgs struct {
	Command []string          `json:"command"`
	Dir     string            `json:"dir"`
	Env     map[string]string `json:"env"`
	Input   *ReachText        `json:"input"`
	Timeout Duration          `json:"timeout"`
	Require CmdArgsRequire    `json:"require"`
}

func (e *Executor) Cmd(args *CmdArgs) (step.Function, error) {
	if err := args.Validate(); err != nil {
		return nil, safe.Wrap(err, "cmd")
	}
	env := make([]string, 0, len(args.Env))
	for key, value := range args.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	timeout := time.Minute
	if args.Timeout.Duration > 0 {
		timeout = args.Timeout.Duration
	}

	return func() (*step.Result, error) {
		var (
			status     int
			output     bytes.Buffer
			stderr     bytes.Buffer
			stdin, err = args.Input.Value()
		)
		if err != nil {
			return nil, err
		}
		defer safe.Close(stdin, "cmd: input")
		ctx, cancel := context.WithTimeout(e.ctx, timeout)
		defer cancel()
		// #nosec G204
		// region Command
		cmd := exec.CommandContext(ctx, args.Command[0], args.Command[1:]...)
		cmd.Stdout = &output
		cmd.Stderr = &stderr
		cmd.Stdin = stdin
		cmd.Env = env
		if args.Dir != "" {
			cmd.Dir = args.Dir
		}
		err = cmd.Run()
		if err != nil {
			if res, ok := err.(*exec.ExitError); ok {
				status = res.ExitCode()
			} else {
				return nil, safe.Wrap(err, "can't run command")
			}
		} else {
			status = cmd.ProcessState.ExitCode()
		}
		// endregion
		// region Match
		if err = args.Require.Match(status, output.Bytes(), stderr.Bytes()); err != nil {
			return nil, fmt.Errorf("cmd: %w", err)
		}
		// endregion
		return step.NewResultSuccess("OK"), nil
	}, nil
}

func (a *CmdArgs) Validate() (err error) {
	if len(a.Command) == 0 {
		return fmt.Errorf("command: field `command` is required")
	}
	if a.Dir != "" {
		stat, err := os.Stat(a.Dir)
		if err != nil {
			return safe.Wrap(err, "dir")
		}
		if !stat.IsDir() {
			return fmt.Errorf("dir: the directory should exist and be a valid directory")
		}
	}
	if err = a.Input.Validate(); err != nil {
		return safe.Wrap(err, "input")
	}
	if err = a.Timeout.Validate(); err != nil {
		return safe.Wrap(err, "timeout")
	}
	if err = a.Require.Validate(); err != nil {
		return safe.Wrap(err, "require")
	}
	return nil
}
