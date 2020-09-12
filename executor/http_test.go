package executor

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	. "github.com/spyzhov/healthy/executor/internal/args"
	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
)

type contextKey string

func TestExecutor_Http(t *testing.T) {
	// region helpers
	newPayload := func(s string) *ReachText {
		return &ReachText{
			Text: &s,
		}
	}
	newFloat64 := func(f float64) *float64 {
		return &f
	}
	isValue := func(e *Executor, key, value string) {
		if actual, ok := e.ctx.Value(contextKey(key)).(string); !ok {
			t.Errorf("Value is not set for: %s", key)
		} else if actual != value {
			t.Errorf("Value is not valid for: %s (%s != %s)", key, value, actual)
		}
	}
	isNotValue := func(e *Executor, key, value string) {
		actual, ok := e.ctx.Value(contextKey(key)).(string)
		if ok && actual == value {
			t.Errorf("Value is not valid for: %s (%s != %s)", key, value, actual)
		}
	}
	setMethod := func(e *Executor, request *http.Request) {
		e.ctx = context.WithValue(e.ctx, contextKey("method"), request.Method)
	}
	setBody := func(e *Executor, request *http.Request) {
		defer safe.Close(request.Body, "request.Body")
		data, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
		e.ctx = context.WithValue(e.ctx, contextKey("body"), string(data))
	}
	setFormValues := func(e *Executor, request *http.Request) {
		if request.MultipartForm != nil {
			for key, values := range request.MultipartForm.Value {
				for i, value := range values {
					e.ctx = context.WithValue(e.ctx, contextKey(fmt.Sprintf("form.values.%s.%d", key, i)), value)
				}
			}
		}
	}
	setHeaders := func(e *Executor, request *http.Request) {
		for key, values := range request.Header {
			for i, value := range values {
				e.ctx = context.WithValue(e.ctx, contextKey(fmt.Sprintf("header.%s.%d", key, i)), value)
			}
		}
	}
	setContent := func(t *testing.T, w http.ResponseWriter, value string) {
		n, err := w.Write([]byte(value))
		if err != nil {
			t.Errorf("http.ResponseWriter.Write() error: %s", err)
		}
		if n != len(value) {
			t.Errorf("http.ResponseWriter.Write() length: %v != %v", n, len(value))
		}
	}
	// endregion
	tests := []struct {
		name     string
		e        *Executor
		getArgs  func(e *Executor) (args *HttpArgs, server *httptest.Server, err error)
		validate func(e *Executor)
		want     step.Function
		wantErr  bool
	}{
		{
			name: "nil",
			e:    NewExecutor(context.Background(), ""),
			getArgs: func(e *Executor) (args *HttpArgs, server *httptest.Server, err error) {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					setMethod(e, r)
				}))
				args = &HttpArgs{
					Method:    "",
					Url:       server.URL,
					Payload:   nil,
					Form:      HttpArgsForm{},
					Headers:   nil,
					Timeout:   Duration{},
					Redirect:  false,
					BasicAuth: nil,
					Require:   HttArgsRequire{},
				}
				return args, server, err
			},
			validate: func(e *Executor) {
				isValue(e, "method", "GET")
			},
			want: func() (*step.Result, error) {
				return step.NewResultSuccess("OK"), nil
			},
			wantErr: false,
		},
		{
			name: "simple_post_200_valid",
			e:    NewExecutor(context.Background(), ""),
			getArgs: func(e *Executor) (args *HttpArgs, server *httptest.Server, err error) {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					setMethod(e, r)
					setBody(e, r)
				}))
				args = &HttpArgs{
					Method:    "POST",
					Url:       server.URL,
					Payload:   newPayload(`value`),
					Form:      HttpArgsForm{},
					Headers:   nil,
					Timeout:   Duration{},
					Redirect:  false,
					BasicAuth: nil,
					Require: HttArgsRequire{
						Status: &HttpArgsRequireStatus{
							RequireNumeric: RequireNumeric{
								In: []float64{200},
							},
						},
					},
				}
				return args, server, err
			},
			validate: func(e *Executor) {
				isValue(e, "method", "POST")
				isValue(e, "body", "value")
			},
			want: func() (*step.Result, error) {
				return step.NewResultSuccess("OK"), nil
			},
			wantErr: false,
		},
		{
			name: "simple_post_201_invalid",
			e:    NewExecutor(context.Background(), ""),
			getArgs: func(e *Executor) (args *HttpArgs, server *httptest.Server, err error) {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					setMethod(e, r)
					setBody(e, r)
				}))
				args = &HttpArgs{
					Method:    "POST",
					Url:       server.URL,
					Payload:   newPayload(`value`),
					Form:      HttpArgsForm{},
					Headers:   nil,
					Timeout:   Duration{},
					Redirect:  false,
					BasicAuth: nil,
					Require: HttArgsRequire{
						Status: &HttpArgsRequireStatus{
							RequireNumeric: RequireNumeric{
								In: []float64{201},
							},
						},
					},
				}
				return args, server, err
			},
			validate: func(e *Executor) {
				isValue(e, "method", "POST")
				isValue(e, "body", "value")
			},
			want: func() (*step.Result, error) {
				return nil, fmt.Errorf("http: status: value 200 is not IN list: [201]")
			},
			wantErr: false,
		},
		{
			name: "simple_post_timeout_invalid",
			e:    NewExecutor(context.Background(), ""),
			getArgs: func(e *Executor) (args *HttpArgs, server *httptest.Server, err error) {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					setMethod(e, r)
				}))
				args = &HttpArgs{
					Method:  "POST",
					Url:     server.URL,
					Payload: newPayload(`value`),
					Form:    HttpArgsForm{},
					Headers: nil,
					Timeout: Duration{
						Duration: -1,
					},
					Redirect:  false,
					BasicAuth: nil,
					Require:   HttArgsRequire{},
				}
				return args, server, err
			},
			validate: func(e *Executor) {
				isNotValue(e, "method", "POST")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "simple_post_payload_form_invalid",
			e:    NewExecutor(context.Background(), ""),
			getArgs: func(e *Executor) (args *HttpArgs, server *httptest.Server, err error) {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					setMethod(e, r)
				}))
				args = &HttpArgs{
					Method:  "POST",
					Url:     server.URL,
					Payload: newPayload(`value`),
					Form: HttpArgsForm{
						Values: map[string]string{
							"key": "value",
						},
					},
					Headers:   nil,
					Timeout:   Duration{},
					Redirect:  false,
					BasicAuth: nil,
					Require:   HttArgsRequire{},
				}
				return args, server, err
			},
			validate: func(e *Executor) {
				isNotValue(e, "method", "POST")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "simple_post_form_invalid",
			e:    NewExecutor(context.Background(), ""),
			getArgs: func(e *Executor) (args *HttpArgs, server *httptest.Server, err error) {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					setMethod(e, r)
				}))
				args = &HttpArgs{
					Method:  "POST",
					Url:     server.URL,
					Payload: nil,
					Form: HttpArgsForm{
						Files: map[string]string{
							"key": "wrong/file\nname.zzz",
						},
					},
					Headers:   nil,
					Timeout:   Duration{},
					Redirect:  false,
					BasicAuth: nil,
					Require:   HttArgsRequire{},
				}
				return args, server, err
			},
			validate: func(e *Executor) {
				isNotValue(e, "method", "POST")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "simple_post_require_invalid",
			e:    NewExecutor(context.Background(), ""),
			getArgs: func(e *Executor) (args *HttpArgs, server *httptest.Server, err error) {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					setMethod(e, r)
				}))
				args = &HttpArgs{
					Method:    "POST",
					Url:       server.URL,
					Payload:   nil,
					Form:      HttpArgsForm{},
					Headers:   nil,
					Timeout:   Duration{},
					Redirect:  false,
					BasicAuth: nil,
					Require: HttArgsRequire{
						Content: &RequireContent{
							Type: "WRONG",
						},
					},
				}
				return args, server, err
			},
			validate: func(e *Executor) {
				isNotValue(e, "method", "POST")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "simple_post_200_invalid",
			e:    NewExecutor(context.Background(), ""),
			getArgs: func(e *Executor) (args *HttpArgs, server *httptest.Server, err error) {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					setMethod(e, r)
					setBody(e, r)
				}))
				args = &HttpArgs{
					Method:    "POST",
					Url:       server.URL,
					Payload:   newPayload(`value`),
					Form:      HttpArgsForm{},
					Headers:   nil,
					Timeout:   Duration{},
					Redirect:  false,
					BasicAuth: nil,
					Require: HttArgsRequire{
						Status: &HttpArgsRequireStatus{
							RequireNumeric: RequireNumeric{
								Eq: newFloat64(500),
							},
						},
					},
				}
				return args, server, err
			},
			validate: func(e *Executor) {
				isValue(e, "method", "POST")
				isValue(e, "body", "value")
			},
			want: func() (*step.Result, error) {
				return nil, fmt.Errorf("http: status: value 200 is not EQ: 500")
			},
			wantErr: false,
		},
		{
			name: "simple_post_form_valid",
			e:    NewExecutor(context.Background(), ""),
			getArgs: func(e *Executor) (args *HttpArgs, server *httptest.Server, err error) {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					err := r.ParseMultipartForm(1024 * 1024)
					if err != nil {
						t.Errorf("ParseMultipartForm() error: %v", err)
					}
					setMethod(e, r)
					setFormValues(e, r)
					setContent(t, w, strings.Repeat(".", 10))
				}))
				args = &HttpArgs{
					Method:  "POST",
					Url:     server.URL,
					Payload: nil,
					Form: HttpArgsForm{
						Values: map[string]string{
							"foo": "bar",
							"key": "value",
						},
						Files: nil,
					},
					Headers: nil,
					Timeout: Duration{
						Duration: time.Hour,
					},
					Redirect:  false,
					BasicAuth: nil,
					Require: HttArgsRequire{
						Content: &RequireContent{
							Length: &RequireNumeric{
								Eq: newFloat64(10),
							},
						},
					},
				}
				return args, server, err
			},
			validate: func(e *Executor) {
				isValue(e, "method", "POST")
				isValue(e, "form.values.foo.0", "bar")
				isValue(e, "form.values.key.0", "value")
			},
			want: func() (*step.Result, error) {
				return step.NewResultSuccess("OK"), nil
			},
			wantErr: false,
		},
		{
			name: "simple_put_header_valid",
			e:    NewExecutor(context.Background(), ""),
			getArgs: func(e *Executor) (args *HttpArgs, server *httptest.Server, err error) {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					setMethod(e, r)
					setHeaders(e, r)
					for key, values := range r.Header {
						for i, value := range values {
							if i == 0 {
								w.Header().Set(key, value)
							} else {
								w.Header().Add(key, value)
							}
						}
					}
				}))
				args = &HttpArgs{
					Method:  "PUT",
					Url:     server.URL,
					Payload: nil,
					Form:    HttpArgsForm{},
					Headers: map[string]string{
						"foo": "bar",
						"key": "value",
					},
					Timeout: Duration{
						Duration: time.Second,
					},
					Redirect:  true,
					BasicAuth: nil,
					Require: HttArgsRequire{
						Header: &HttpArgsRequireHeader{
							Exists:    []string{"foo", "Key"},
							NotExists: []string{"Bar"},
							Regexp: map[string]RequireFieldMatch{
								"foo": {"bar"},
							},
							NotRegexp: map[string]RequireFieldMatchNot{
								"key": {"baz", "biz"},
							},
							Eq: map[string][]string{
								"key": {"value"},
							},
						},
					},
				}
				return args, server, err
			},
			validate: func(e *Executor) {
				isValue(e, "method", "PUT")
				isValue(e, "header.Foo.0", "bar")
				isValue(e, "header.Key.0", "value")
			},
			want: func() (*step.Result, error) {
				return step.NewResultSuccess("OK"), nil
			},
			wantErr: false,
		},
		{
			name: "simple_put_header_invalid",
			e:    NewExecutor(context.Background(), ""),
			getArgs: func(e *Executor) (args *HttpArgs, server *httptest.Server, err error) {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					setMethod(e, r)
					setHeaders(e, r)
					for key, values := range r.Header {
						for i, value := range values {
							if i == 0 {
								w.Header().Set(key, value)
							} else {
								w.Header().Add(key, value)
							}
						}
					}
				}))
				args = &HttpArgs{
					Method:  "PUT",
					Url:     server.URL,
					Payload: nil,
					Form:    HttpArgsForm{},
					Headers: map[string]string{
						"foo": "bar",
						"key": "value",
					},
					Timeout: Duration{
						Duration: time.Second,
					},
					Redirect:  true,
					BasicAuth: nil,
					Require: HttArgsRequire{
						Header: &HttpArgsRequireHeader{
							NotExists: []string{"Foo"},
						},
					},
				}
				return args, server, err
			},
			validate: func(e *Executor) {
				isValue(e, "method", "PUT")
				isValue(e, "header.Foo.0", "bar")
				isValue(e, "header.Key.0", "value")
			},
			want: func() (*step.Result, error) {
				return nil, fmt.Errorf("http: header: EXISTS `Foo`")
			},
			wantErr: false,
		},
		{
			name: "simple_get_auth_valid",
			e:    NewExecutor(context.Background(), ""),
			getArgs: func(e *Executor) (args *HttpArgs, server *httptest.Server, err error) {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					setMethod(e, r)
					setHeaders(e, r)
					if r.Header.Get("Authorization") != "Basic Zmlkbzpjb2NhY29sYQ==" {
						t.Errorf("wrong Authorization")
					}
				}))
				args = &HttpArgs{
					Method:  "GET",
					Url:     server.URL,
					Payload: nil,
					Form:    HttpArgsForm{},
					Headers: map[string]string{
						"foo": "bar",
						"key": "value",
					},
					Timeout:  Duration{},
					Redirect: false,
					BasicAuth: &HttpArgsBasicAuth{
						Username: "fido",
						Password: "cocacola",
					},
					Require: HttArgsRequire{
						Status: &HttpArgsRequireStatus{
							RequireNumeric: RequireNumeric{
								Eq: newFloat64(200),
							},
						},
					},
				}
				return args, server, err
			},
			validate: func(e *Executor) {
				isValue(e, "method", "GET")
				isValue(e, "header.Authorization.0", "Basic Zmlkbzpjb2NhY29sYQ==")
			},
			want: func() (*step.Result, error) {
				return step.NewResultSuccess("OK"), nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, server, err := tt.getArgs(tt.e)
			if server != nil {
				defer server.Close()
			}
			if err != nil {
				t.Errorf("getArgs() error = %v", err)
				return
			}
			got, err := tt.e.Http(args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Http() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			validate(t, got, tt.want)
			tt.validate(tt.e)
		})
	}
}
