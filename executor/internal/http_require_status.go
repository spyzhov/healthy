package internal

type HttpArgsRequireStatus struct {
	RequireNumeric
}

func (a *HttpArgsRequireStatus) Match(status int) error {
	return a.RequireNumeric.Match("status", float64(status))
}
