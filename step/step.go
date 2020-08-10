package step

type Step struct {
	Name string
	Func Function
}

type Function func() (*Result, error)
