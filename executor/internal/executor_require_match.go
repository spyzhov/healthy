package internal

type RequireMatch struct {
	Regexp    RequireFieldMatch    `json:"match"`
	NotRegexp RequireFieldMatchNot `json:"not_match"`
}

func (a *RequireMatch) Validate() (err error) {
	if err = a.Regexp.Validate(); err != nil {
		return err
	}
	if err = a.NotRegexp.Validate(); err != nil {
		return err
	}
	return nil
}

func (a *RequireMatch) Match(name string, input []byte) (err error) {
	if err = a.Regexp.Match(name, input); err != nil {
		return err
	}
	if err = a.NotRegexp.Match(name, input); err != nil {
		return err
	}
	return nil
}

func (a *RequireMatch) MatchStrings(name string, input []string) (err error) {
	if err = a.Regexp.MatchStrings(name, input); err != nil {
		return err
	}
	if err = a.NotRegexp.MatchStrings(name, input); err != nil {
		return err
	}
	return nil
}
