package executor

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	. "github.com/spyzhov/healthy/executor/internal/args"
	"github.com/spyzhov/healthy/executor/internal/net/http/transport"
	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
)

type HttpArgs struct {
	// region Request
	Method    string             `json:"method"`
	Url       string             `json:"url"`
	Payload   *string            `json:"payload"`
	Form      HttpArgsForm       `json:"form"`
	Headers   map[string]string  `json:"headers"`
	Timeout   Duration           `json:"timeout"`
	Redirect  bool               `json:"redirect"`
	BasicAuth *HttpArgsBasicAuth `json:"basic_auth"`
	// endregion
	// region Response
	Require HttArgsRequire `json:"require"`
	// endregion
}

func (e *Executor) Http(args *HttpArgs) (step.Function, error) {
	if err := args.Validate(); err != nil {
		return nil, safe.Wrap(err, "http")
	}
	client := func(timeout time.Duration) *http.Client {
		if timeout == 0 {
			timeout = 30 * time.Second
		}
		if args.Redirect {
			return &http.Client{
				Timeout:   timeout,
				Transport: transport.NewAgent("healthy/" + e.version),
			}
		}
		return &http.Client{
			Timeout:   timeout,
			Transport: transport.NewAgent("healthy/" + e.version),
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}
	return func() (*step.Result, error) {
		// region Request
		contentType, body, err := args.body()
		if err != nil {
			return nil, fmt.Errorf("http: %w", err)
		}
		request, err := http.NewRequest(args.method(), args.Url, body)
		if err != nil {
			return nil, fmt.Errorf("http: %w", err)
		}
		if contentType != "" {
			request.Header.Add("Content-Type", contentType)
		}
		for key, value := range args.Headers {
			request.Header.Set(key, value)
		}
		if args.BasicAuth != nil {
			request.SetBasicAuth(args.BasicAuth.Username, args.BasicAuth.Password)
		}
		// endregion
		// region Response
		var response *http.Response
		response, err = client(args.Timeout.Duration).Do(request.WithContext(e.ctx))
		if err != nil {
			return nil, fmt.Errorf("http: %w", err)
		}
		defer safe.Close(response.Body, "http:response.body")
		// endregion
		// region Match
		if err = args.Require.Match(response); err != nil {
			return nil, fmt.Errorf("http: %w", err)
		}
		// endregion
		return step.NewResultSuccess("OK"), nil
	}, nil
}

func (a *HttpArgs) Validate() (err error) {
	if err = a.Timeout.Validate(); err != nil {
		return err
	}
	if err = a.Form.Validate(); err != nil {
		return err
	}
	if err = a.Require.Validate(); err != nil {
		return err
	}
	if a.Payload != nil && (len(a.Form.Files)+len(a.Form.Values) != 0) {
		return fmt.Errorf("body: http.payload and http.form set in the same time")
	}
	return
}

func (a *HttpArgs) method() string {
	if a.Method == "" {
		return "GET"
	}
	return a.Method
}

func (a *HttpArgs) body() (contentType string, r io.Reader, err error) {
	var b bytes.Buffer
	if a.Payload != nil {
		_, err = b.WriteString(*a.Payload)
		if err != nil {
			return "", nil, fmt.Errorf("write buffer: %w", err)
		}
	} else {
		contentType, err = a.Form.SubmitForm(&b)
		if err != nil {
			return "", nil, err
		}
	}
	return contentType, &b, err
}
