package internal

type Validatable interface {
	Validate() (err error)
}
