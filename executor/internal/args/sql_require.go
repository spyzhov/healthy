package args

import "github.com/spyzhov/safe"

type SqlArgsRequire struct {
	Count RequireNumeric       `json:"count"`
	Rows  RequireTable         `json:"rows"`
	Value SqlArgsRequireValues `json:"value"`
}

func (a *SqlArgsRequire) Validate() (err error) {
	if err = a.Count.Validate(); err != nil {
		return safe.Wrap(err, "count")
	}
	if err = a.Rows.Validate(); err != nil {
		return safe.Wrap(err, "rows")
	}
	if err = a.Value.Validate(); err != nil {
		return safe.Wrap(err, "value")
	}
	return nil
}

func (a *SqlArgsRequire) Match(rows [][]interface{}) (err error) {
	if err = a.Count.Match("count", float64(len(rows))); err != nil {
		return err
	}
	if err = a.Rows.Match(rows); err != nil {
		return safe.Wrap(err, "rows")
	}
	if err = a.Value.Match(rows); err != nil {
		return err
	}
	return nil
}
