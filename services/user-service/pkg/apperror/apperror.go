package apperror

import (
	"errors"
	"fmt"
)

type AppError struct {
	Layer   string
	Op      string
	Message string
	Err     error
}

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
)

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s:%s] %s: %v", e.Layer, e.Op, e.Message, e.Err)
	}

	return fmt.Sprintf("[%s:%s] %s", e.Layer, e.Op, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func Wrap(layer, op, message string, err error) error {
	return &AppError{
		Layer:   layer,
		Op:      op,
		Message: message,
		Err:     err,
	}
}
