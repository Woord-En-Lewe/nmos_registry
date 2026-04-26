package registry

import (
	"errors"
	"fmt"
)

var ErrResourceNotFound = errors.New("resource not found")

type ValidationError struct {
    Message string
    Cause   error
}

func (e *ValidationError) Error() string {
    if e.Cause != nil {
        return e.Message + ": " + e.Cause.Error()
    }
    return e.Message
}

func (e *ValidationError) Unwrap() error {
    return e.Cause
}

type ConflictError struct {
    Message string
    Cause   error
}

func (e *ConflictError) Error() string {
    if e.Cause != nil {
        return e.Message + ": " + e.Cause.Error()
    }
    return e.Message
}

func (e *ConflictError) Unwrap() error {
    return e.Cause
}

func NewValidationError(msg string) *ValidationError {
    return &ValidationError{Message: msg}
}

func NewValidationErrorf(format string, args ...interface{}) *ValidationError {
    return &ValidationError{Message: fmt.Sprintf(format, args...)}
}

func WrapValidationError(err error, msg string) *ValidationError {
    return &ValidationError{Message: msg, Cause: err}
}

func NewConflictError(msg string) *ConflictError {
    return &ConflictError{Message: msg}
}

func NewConflictErrorf(format string, args ...interface{}) *ConflictError {
    return &ConflictError{Message: fmt.Sprintf(format, args...)}
}

func WrapConflictError(err error, msg string) *ConflictError {
    return &ConflictError{Message: msg, Cause: err}
}