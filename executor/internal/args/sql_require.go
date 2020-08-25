package args

import "github.com/spyzhov/safe"

type SqlArgsRequire struct {
	Count RequireNumeric       `json:"count"`
	Rows  Table                `json:"rows"`
	Value SqlArgsRequireValues `json:"value"`
}

func (a *SqlArgsRequire) Validate() (err error) {
	if a == nil {
		return nil
	}
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

func (a *SqlArgsRequire) Match(rows Table) (err error) {
	if a == nil {
		return nil
	}
	if err = rows.Validate(); err != nil {
		return safe.Wrap(err, "examine value")
	}
	if err = a.Count.Match("count", float64(len(rows))); err != nil {
		return err
	}
	if err = a.Rows.Match(rows, "NULL"); err != nil {
		return safe.Wrap(err, "rows")
	}
	if err = a.Value.Match(rows); err != nil {
		return err
	}
	return nil
}
