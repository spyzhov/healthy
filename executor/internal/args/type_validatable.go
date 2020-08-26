package args

type Validatable interface {
	Validate() (err error)
}
