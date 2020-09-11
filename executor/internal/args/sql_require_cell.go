package args

type SqlArgsRequireCell struct {
	Cell
	SqlArgsRequireValue
}

func (a *SqlArgsRequireCell) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.SqlArgsRequireValue.Validate(); err != nil {
		return err
	}
	return nil
}

func (a *SqlArgsRequireCell) Match(rows Table) (err error) {
	if a == nil {
		return nil
	}
	if err = a.Cell.Match(rows); err != nil {
		return err
	}
	if err = a.SqlArgsRequireValue.Match(a.Cell.get(rows)); err != nil {
		return err
	}
	return nil
}
