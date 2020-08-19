package args

import (
	"fmt"

	"github.com/spyzhov/healthy/executor/internal"
)

type RequireNumeric struct {
	In    []float64 `json:"in"`
	NotIn []float64 `json:"not_in"`
	Eq    *float64  `json:"eq"`
	Le    *float64  `json:"le"`
	Leq   *float64  `json:"leq"`
	Ge    *float64  `json:"ge"`
	Geq   *float64  `json:"geq"`
	Not   *float64  `json:"not"`
}

func (a *RequireNumeric) Match(name string, numeric float64) error {
	if len(a.In) > 0 {
		if !internal.FloatInSlice(numeric, a.In) {
			return fmt.Errorf("%s: value %v is not IN list: %v", name, numeric, a.In)
		}
	}
	if len(a.NotIn) > 0 {
		if internal.FloatInSlice(numeric, a.NotIn) {
			return fmt.Errorf("%s: value %v is in NOT_IN list: %v", name, numeric, a.NotIn)
		}
	}
	if a.Eq != nil {
		if numeric != *a.Eq {
			return fmt.Errorf("%s: value %v is not EQ: %v", name, numeric, *a.Eq)
		}
	}
	if a.Not != nil {
		if numeric == *a.Not {
			return fmt.Errorf("%s: value %v is eq NOT: %v", name, numeric, *a.Not)
		}
	}
	if a.Geq != nil {
		if numeric < *a.Geq {
			return fmt.Errorf("%s: value %v is not GEQ: %v", name, numeric, *a.Geq)
		}
	}
	if a.Leq != nil {
		if numeric > *a.Leq {
			return fmt.Errorf("%s: value %v is not LEQ: %v", name, numeric, *a.Leq)
		}
	}
	if a.Ge != nil {
		if numeric <= *a.Ge {
			return fmt.Errorf("%s: value %v is not GE: %v", name, numeric, *a.Ge)
		}
	}
	if a.Le != nil {
		if numeric >= *a.Le {
			return fmt.Errorf("%s: value %v is not LE: %v", name, numeric, *a.Le)
		}
	}
	return nil
}
