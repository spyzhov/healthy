package args

import "github.com/spyzhov/safe"

type SqlArgsRequire struct {
	Count  RequireNumeric        `json:"count"`
	Rows   SqlArgsRequireTable   `json:"rows"`
	Row    SqlArgsRequireRows    `json:"row"`
	Column SqlArgsRequireColumns `json:"column"`
	Cell   SqlArgsRequireCells   `json:"cell"`
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
	if err = a.Row.Validate(); err != nil {
		return err
	}
	if err = a.Column.Validate(); err != nil {
		return err
	}
	if err = a.Cell.Validate(); err != nil {
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
	if err = a.Rows.Match(rows); err != nil {
		return safe.Wrap(err, "rows")
	}
	if err = a.Row.Match(rows); err != nil {
		return err
	}
	if err = a.Column.Match(rows); err != nil {
		return err
	}
	if err = a.Cell.Match(rows); err != nil {
		return err
	}
	return nil
}
