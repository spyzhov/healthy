package executor

import "fmt"

type RequireInteger struct {
	In    []int `json:"in"`
	NotIn []int `json:"not_in"`
	Eq    *int  `json:"eq"`
	Le    *int  `json:"le"`
	Leq   *int  `json:"leq"`
	Ge    *int  `json:"ge"`
	Geq   *int  `json:"geq"`
	Not   *int  `json:"not"`
}

func (a *RequireInteger) Match(name string, integer int) error {
	if len(a.In) > 0 {
		if !intInSlice(integer, a.In) {
			return fmt.Errorf("%s: %d is not IN list: %v", name, integer, a.In)
		}
	}
	if len(a.NotIn) > 0 {
		if intInSlice(integer, a.NotIn) {
			return fmt.Errorf("%s: %d is in NOT_IN list: %v", name, integer, a.NotIn)
		}
	}
	if a.Eq != nil {
		if integer != *a.Eq {
			return fmt.Errorf("%s:  %d is not EQ: %d", name, integer, *a.Eq)
		}
	}
	if a.Not != nil {
		if integer == *a.Not {
			return fmt.Errorf("%s:  %d is eq NOT: %d", name, integer, *a.Not)
		}
	}
	if a.Geq != nil {
		if integer < *a.Geq {
			return fmt.Errorf("%s:  %d is not GEQ: %d", name, integer, *a.Geq)
		}
	}
	if a.Leq != nil {
		if integer > *a.Leq {
			return fmt.Errorf("%s:  %d is not LEQ: %d", name, integer, *a.Leq)
		}
	}
	if a.Ge != nil {
		if integer <= *a.Ge {
			return fmt.Errorf("%s:  %d is not GE: %d", name, integer, *a.Geq)
		}
	}
	if a.Le != nil {
		if integer >= *a.Le {
			return fmt.Errorf("%s:  %d is not LE: %d", name, integer, *a.Leq)
		}
	}
	return nil
}
