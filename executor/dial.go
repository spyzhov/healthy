package executor

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	. "github.com/spyzhov/healthy/executor/internal/args"
	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
)

type DialArgs struct {
	// region Request
	Type    String     `json:"type"`
	Address String     `json:"address"`
	Input   *ReachText `json:"input"`
	Until   string     `json:"until"`
	Timeout Duration   `json:"timeout"`
	// endregion
	// region Require
	Require *DialArgsRequire `json:"require"`
	// endregion
}

func (e *Executor) Dial(args *DialArgs) (step.Function, error) {
	if err := args.Validate(); err != nil {
		return nil, safe.Wrap(err, "dial")
	}
	// region Default
	if args.Timeout.Duration == 0 {
		args.Timeout.Duration = time.Second
	}
	// endregion
	return func() (*step.Result, error) {
		// region Request
		conn, err := net.DialTimeout(args.Type.Value(), args.Address.Value(), args.Timeout.Duration)
		if err != nil {
			return nil, safe.Wrap(err, "dial: connect")
		}
		err = conn.SetReadDeadline(time.Now().Add(args.Timeout.Duration))
		if err != nil {
			return nil, safe.Wrap(err, "dial: config")
		}
		defer safe.Close(conn, "dial: connection")
		var input io.ReadCloser
		input, err = args.Input.Value()
		if err != nil {
			return nil, err
		}
		defer safe.Close(input, "dial: input")
		_, err = io.Copy(conn, input)
		if err != nil {
			return nil, safe.Wrap(err, "dial: write")
		}
		// endregion
		// region Response
		response := make([]byte, 0)
		var message []byte
		reader := bufio.NewReader(conn)
		for {
			message, err = reader.ReadBytes('\n')
			if err != nil {
				if op, ok := err.(*net.OpError); ok && op.Timeout() {
					break
				}
				if errors.Is(err, context.DeadlineExceeded) {
					break
				}
				return nil, safe.Wrap(err, "dial: read")
			}
			if args.Until != "" {
				index := strings.Index(string(message), args.Until)
				if index >= 0 {
					response = append(response, message[:index+len(args.Until)]...)
					break
				}
			}
			response = append(response, message...)
		}
		// endregion
		// region Match
		if err = args.Require.Match(response); err != nil {
			return nil, fmt.Errorf("dial: %w", err)
		}
		// endregion
		return step.NewResultSuccess("OK"), nil
	}, nil
}

func (a *DialArgs) Validate() (err error) {
	if err = a.Type.Validate(); err != nil {
		return safe.Wrap(err, "type")
	}
	if err = a.Input.Validate(); err != nil {
		return safe.Wrap(err, "input")
	}
	if err = a.Address.Validate(); err != nil {
		return safe.Wrap(err, "address")
	}
	if err = a.Timeout.Validate(); err != nil {
		return safe.Wrap(err, "timeout")
	}
	if err = a.Require.Validate(); err != nil {
		return err
	}
	return
}
