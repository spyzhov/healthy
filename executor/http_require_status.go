package executor

type HttpArgsRequireStatus struct {
	RequireInteger
}

func (a *HttpArgsRequireStatus) Match(status int) error {
	return a.RequireInteger.Match("status", status)
}
