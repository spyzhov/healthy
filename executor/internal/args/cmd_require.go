package args

import "github.com/spyzhov/safe"

type CmdArgsRequire struct {
	ExitCode *RequireNumeric `json:"exit_code"`
	Output   *RequireContent `json:"output"`
	Stderr   *RequireContent `json:"stderr"`
}

func (a *CmdArgsRequire) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.ExitCode.Validate(); err != nil {
		return safe.Wrap(err, "status")
	}
	if err = a.Output.Validate(); err != nil {
		return safe.Wrap(err, "output")
	}
	if err = a.Stderr.Validate(); err != nil {
		return safe.Wrap(err, "stderr")
	}
	return nil
}

func (a *CmdArgsRequire) Match(exitCode int, output []byte, stderr []byte) (err error) {
	if a == nil {
		return nil
	}
	if err = a.ExitCode.Match("exit_code", float64(exitCode)); err != nil {
		return err
	}
	if err = a.Output.Match("output", output); err != nil {
		return safe.Wrap(err, "output")
	}
	if err = a.Stderr.Match("stderr", stderr); err != nil {
		return safe.Wrap(err, "stderr")
	}
	return nil
}
