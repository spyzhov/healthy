package args

type HttpArgsRequireStatus struct {
	RequireNumeric
}

func (a *HttpArgsRequireStatus) Validate() error {
	return a.RequireNumeric.Validate()
}

func (a *HttpArgsRequireStatus) Match(status int) error {
	return a.RequireNumeric.Match("status", float64(status))
}
