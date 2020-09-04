package executor

import (
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	. "github.com/spyzhov/healthy/executor/internal/args"
	http2 "github.com/spyzhov/healthy/executor/internal/net/http"
	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
)

type AppConsulArgs struct {
	// region Request
	Config *ConfigAppConsulArgs `json:"config"`
	// endregion
	// region Require
	Require AppConsulArgsRequire `json:"require"`
	// endregion
}

func (e *Executor) AppConsul(args *AppConsulArgs) (step.Function, error) {
	scope := "app/consul"
	if err := args.Validate(); err != nil {
		return nil, safe.Wrap(err, scope)
	}

	return func() (*step.Result, error) {
		// region Prepare
		client, err := api.NewClient(args.Config.config(30*time.Second, e.version))
		if err != nil {
			return nil, safe.Wrap(err, scope+": connect")
		}
		// endregion
		// region Match
		if err = args.Require.Match(client); err != nil {
			return nil, safe.Wrap(err, scope)
		}
		// endregion
		return step.NewResultSuccess("OK"), nil
	}, nil
}

func (a *AppConsulArgs) Validate() (err error) {
	if err = a.Require.Validate(); err != nil {
		return err
	}
	return
}

type ConfigAppConsulArgs struct {
	Address    string                       `json:"address"`
	Scheme     string                       `json:"scheme"`
	Datacenter string                       `json:"datacenter"`
	HttpAuth   *HttpAuthConfigAppConsulArgs `json:"http_auth"`
	WaitTime   time.Duration                `json:"wait_time"`
	Token      string                       `json:"token"`
	TokenFile  string                       `json:"token_file"`
	Namespace  string                       `json:"namespace"`
	TLSConfig  TLSConfigConfigAppConsulArgs `json:"tls_config"`
}

type TLSConfigConfigAppConsulArgs struct {
	Address            string `json:"address"`
	CAFile             string `json:"ca_file"`
	CAPath             string `json:"ca_path"`
	CAPem              []byte `json:"ca_pem"`
	CertFile           string `json:"cert_file"`
	CertPEM            []byte `json:"cert_pem"`
	KeyFile            string `json:"key_file"`
	KeyPEM             []byte `json:"key_pem"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
}

type HttpAuthConfigAppConsulArgs struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *ConfigAppConsulArgs) config(timeout time.Duration, version string) *api.Config {
	if a == nil {
		return api.DefaultNonPooledConfig()
	}
	return &api.Config{
		Address:    a.Address,
		Scheme:     a.Scheme,
		Datacenter: a.Datacenter,
		Transport:  cleanhttp.DefaultTransport(),
		HttpClient: http2.GetClient(timeout, version),
		HttpAuth:   a.HttpAuth.auth(),
		WaitTime:   a.WaitTime,
		Token:      a.Token,
		TokenFile:  a.TokenFile,
		Namespace:  a.Namespace,
		TLSConfig: api.TLSConfig{
			Address:            a.TLSConfig.Address,
			CAFile:             a.TLSConfig.CAFile,
			CAPath:             a.TLSConfig.CAPath,
			CAPem:              a.TLSConfig.CAPem,
			CertFile:           a.TLSConfig.CertFile,
			CertPEM:            a.TLSConfig.CertPEM,
			KeyFile:            a.TLSConfig.KeyFile,
			KeyPEM:             a.TLSConfig.KeyPEM,
			InsecureSkipVerify: a.TLSConfig.InsecureSkipVerify,
		},
	}
}

func (a *HttpAuthConfigAppConsulArgs) auth() *api.HttpBasicAuth {
	if a == nil || (a.Password+a.Username) == "" {
		return nil
	}
	return &api.HttpBasicAuth{
		Username: a.Username,
		Password: a.Password,
	}
}
