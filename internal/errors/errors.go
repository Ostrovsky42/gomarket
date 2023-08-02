package errors

import "fmt"

const (
	FailedValidation    = "Failed validation"
	NotFound            = "Resource not found"
	InternalServerError = "Internal server error"
	Unauthorized        = "Unauthorized"
	UniquenessViolation = "Uniqueness Violation"
	InsufficientFunds   = "Insufficient funds"
)

type ErrorApp struct {
	description string
	details     any
}

func (e *ErrorApp) Error() string {
	return fmt.Sprintf("%#v", *e)
}

func (e *ErrorApp) Description() string {
	return e.description
}

func NewErrFailedValidation(details any) *ErrorApp {
	return NewError(FailedValidation, details)
}

func NewErrInternal(details any) *ErrorApp {
	return NewError(InternalServerError, details)
}

func NewErrUnauthorized() *ErrorApp {
	return NewError(Unauthorized, "")
}

func NewErrNotFound() *ErrorApp {
	return NewError(NotFound, "")
}
func NewErrInsufficientFunds() *ErrorApp {
	return NewError(InsufficientFunds, "")
}

func NewErrUniquenessViolation(details any) *ErrorApp {
	return NewError(UniquenessViolation, details)
}

func NewError(description string, details any) *ErrorApp {
	return &ErrorApp{description: description, details: details}
}
