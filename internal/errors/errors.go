package errors

import "fmt"

const (
	BadRequestJSON      = "Bad request JSON"
	FailedValidation    = "Failed validation"
	NotFound            = "Resource not found"
	InternalServerError = "Internal server error"
	Unauthorized        = "Unauthorized"
	UniquenessViolation = "Uniqueness Violation"
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

func (e *ErrorApp) Details() any {
	return e.details
}

func NewErrBadRequestJSON(details any) *ErrorApp {
	return NewError(BadRequestJSON, details)
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

func NewErrUniquenessViolation(details any) *ErrorApp {
	return NewError(UniquenessViolation, details)
}

func NewError(description string, details any) *ErrorApp {
	return &ErrorApp{description: description, details: details}
}
