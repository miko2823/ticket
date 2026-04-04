package errors

import "fmt"

// AppError is a custom error type that carries an HTTP status code and a user-facing message.
type AppError struct {
	Status  int
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func New(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}

func Wrap(status int, message string, cause error) *AppError {
	return &AppError{Status: status, Message: message, Cause: cause}
}
