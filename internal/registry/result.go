package registry

import "fmt"

type Error struct {
    Code    int
    Message string
    Cause   error
}

func (e *Error) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

func (e *Error) Unwrap() error {
    return e.Cause
}

type Result[T any] struct {
    value     T
    err       *Error
    isCreated bool
}

func Success[T any](value T) Result[T] {
    return Result[T]{value: value, err: nil, isCreated: false}
}

func Created[T any](value T) Result[T] {
    return Result[T]{value: value, err: nil, isCreated: true}
}

func Failure[T any](code int, message string, cause error) Result[T] {
    return Result[T]{err: &Error{Code: code, Message: message, Cause: cause}}
}

func Failuref[T any](code int, format string, args ...interface{}) Result[T] {
    return Failure[T](code, fmt.Sprintf(format, args...), nil)
}

func (r Result[T]) IsSuccess() bool {
    return r.err == nil
}

func (r Result[T]) IsFailure() bool {
    return r.err != nil
}

func (r Result[T]) WasCreated() bool {
    return r.isCreated
}

func (r Result[T]) Value() (T, bool) {
    return r.value, r.err == nil
}

func (r Result[T]) Error() (*Error, bool) {
    return r.err, r.err != nil
}

func (r Result[T]) Unwrap() T {
    if r.err != nil {
        panic("unwrap called on Failure Result")
    }
    return r.value
}

func (r Result[T]) UnwrapOr(defaultValue T) T {
    if r.err != nil {
        return defaultValue
    }
    return r.value
}

func (r Result[T]) Map(f func(T) T) Result[T] {
    if r.err != nil {
        return Failure[T](r.err.Code, r.err.Message, r.err.Cause)
    }
    return Success[T](f(r.value))
}

func (r Result[T]) Bind(f func(T) Result[T]) Result[T] {
    if r.err != nil {
        return r
    }
    return f(r.value)
}

func Bind[T, U any](r Result[T], f func(T) Result[U]) Result[U] {
    if r.err != nil {
        return Failure[U](r.err.Code, r.err.Message, r.err.Cause)
    }
    return f(r.value)
}

func (r Result[T]) MapError(f func(*Error) *Error) Result[T] {
    if r.err == nil {
        return r
    }
    newErr := f(r.err)
    return Failure[T](newErr.Code, newErr.Message, newErr.Cause)
}

func (r Result[T]) Recover(f func(*Error) T) Result[T] {
    if r.err == nil {
        return r
    }
    return Success[T](f(r.err))
}

func (r Result[T]) OrElse(other Result[T]) Result[T] {
    if r.err != nil {
        return other
    }
    return r
}

func (r Result[T]) AndThen(f func(T)) Result[T] {
    if r.err != nil {
        return r
    }
    f(r.value)
    return r
}

func (r Result[T]) AndThenError(f func(T) error) Result[T] {
    if r.err != nil {
        return r
    }
    if err := f(r.value); err != nil {
        return Failure[T](400, err.Error(), err)
    }
    return r
}

func (r Result[T]) Fold(success func(T) T, failure func(*Error) T) T {
    if r.err != nil {
        return failure(r.err)
    }
    return success(r.value)
}

func (r Result[T]) ToVoid() Result[struct{}] {
    if r.err != nil {
        return Failure[struct{}](r.err.Code, r.err.Message, r.err.Cause)
    }
    if r.isCreated {
        return Created(struct{}{})
    }
    return Success(struct{}{})
}