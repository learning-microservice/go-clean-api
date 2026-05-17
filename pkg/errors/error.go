package errors

import (
	"encoding/json"
	"fmt"
)

type Error[T any] struct {
	errType ErrorType[T]
	message string
	details []T
	attrs   []any
	cause   error
	stack   []uintptr
}

func (e *Error[T]) Error() string {
	if len(e.message) == 0 {
		e.message = "unknown error"
	}
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.errType.name, e.message, e.cause)
	}
	return fmt.Sprintf("[%s] %s", e.errType.name, e.message)
}

func (e *Error[T]) Unwrap() error {
	// Go標準の errors.Is/As で再帰的にUnwrapを呼ぶため
	// ここではそのままcauseを返却する
	return e.cause
}

func (e *Error[T]) TypeCode() int {
	return e.errType.code
}

func (e *Error[T]) TypeName() string {
	return e.errType.name
}

func (e *Error[T]) Message() string {
	return e.message
}

func (e *Error[T]) Details() []T {
	out := make([]T, len(e.details))
	copy(out, e.details)
	if child := AsError[T](e.cause); child != nil {
		out = append(out, child.Details()...)
	}
	return out
}

func (e *Error[T]) Attrs() []any {
	return e.attrs
}

func (e *Error[T]) Is(target error) bool {
	if err := AsError[T](target); err != nil {
		return e.errType.equals(err.errType)
	}
	return false
}

func (e *Error[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Details []T    `json:"details,omitempty"`
	}{
		Type:    e.errType.name,
		Message: e.message,
		Details: e.Details(),
	})
}
