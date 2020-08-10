package step

import "fmt"

type Status string

const (
	Error   Status = "error"
	Warning Status = "warning"
	Success Status = "success"
)

type Result struct {
	Status  Status
	Message string
}

func NewResult(status Status, message string) *Result {
	return &Result{
		Status:  status,
		Message: message,
	}
}

func NewResultSuccess(message string) *Result {
	return NewResult(Success, message)
}

func NewResultWarning(message string) *Result {
	return NewResult(Warning, message)
}

func NewResultError(message string) *Result {
	return NewResult(Error, message)
}

func (s Status) Validate() error {
	switch s {
	case Success, Warning, Error:
		return nil
	default:
		return fmt.Errorf("status is not valid, require one of: [%s, %s, %s]", Success, Warning, Error)
	}
}
