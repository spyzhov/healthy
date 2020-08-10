package executor

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
)

type HttpArgs struct {
	// region Request
	Method    string              `json:"method"`
	Url       string              `json:"url"`
	Payload   *string             `json:"payload"`
	PostForm  map[string][]string `json:"post_form"`
	Headers   map[string]string   `json:"headers"`
	Timeout   Duration            `json:"timeout"`
	Redirect  bool                `json:"redirect"`
	BasicAuth *struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"basic_auth"`
	// endregion
	// region Response
	Require struct {
		Status  *HttpArgsRequireStatus  `json:"status"`
		Content *HttpArgsRequireContent `json:"content"`
		Header  *HttpArgsRequireHeader  `json:"header"`
	}
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
				Timeout: timeout,
			}
		}
		return &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}
	return func() (*step.Result, error) {
		// region Request
		request, err := http.NewRequest(args.method(), args.Url, args.body())
		if err != nil {
			return nil, fmt.Errorf("http: %w", err)
		}
		for key, values := range args.PostForm {
			for i, value := range values {
				if i == 0 {
					request.Header.Set(key, value)
				} else {
					request.Header.Add(key, value)
				}
			}
		}
		for key, value := range args.Headers {
			request.Header.Set(key, value)
		}
		if args.BasicAuth != nil {
			request.SetBasicAuth(args.BasicAuth.Username, args.BasicAuth.Password)
		}
		// endregion
		// region Response
		response, err := client(args.Timeout.Duration).Do(request.WithContext(e.ctx))
		if err != nil {
			return nil, fmt.Errorf("http: %w", err)
		}
		defer safe.Close(response.Body, "http:response.body")
		// endregion
		// region Match
		if args.Require.Status != nil {
			err = args.Require.Status.Match(response.StatusCode)
			if err != nil {
				return nil, fmt.Errorf("http: %w", err)
			}
		}
		if args.Require.Header != nil {
			err = args.Require.Header.Match(response.Header)
			if err != nil {
				return nil, fmt.Errorf("http: %w", err)
			}
		}
		if args.Require.Content != nil {
			var content []byte
			content, err = ioutil.ReadAll(response.Body)
			err = args.Require.Content.Match(content)
			if err != nil {
				return nil, fmt.Errorf("http: %w", err)
			}
		}
		// endregion
		return step.NewResultSuccess("OK"), nil
	}, nil
}

func (a *HttpArgs) Validate() (err error) {
	if err = a.Timeout.Validate(); err != nil {
		return err
	}
	if a.Require.Content != nil {
		if err = a.Require.Content.Validate(); err != nil {
			return safe.Wrap(err, "require.content")
		}
	}
	return
}

func (a *HttpArgs) method() string {
	if a.Method == "" {
		return "GET"
	}
	return a.Method
}

func (a *HttpArgs) body() io.Reader {
	if a.Payload == nil {
		return nil
	}
	return strings.NewReader(*a.Payload)
}
