package errors

import (
	"fmt"
	"strings"
)

type Error struct {
	errType  Type
	messages []string
	attrs    []any
	cause    error
	stack    []uintptr
}

func (e *Error) Error() string {
	message := strings.Join(e.messages, ",")
	if len(message) == 0 {
		message = "unknown error"
	}
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.errType, message, e.cause)
	}
	return fmt.Sprintf("[%s] %s", e.errType, message)
}

func (e *Error) Unwrap() error {
	// Go標準の errors.Is/As で再帰的にUnwrapを呼ぶため
	// ここではそのままcauseを返却する
	return e.cause
}

func (e *Error) Type() string {
	return string(e.errType)
}

func (e *Error) Message() []string {
	return e.messages
}

func (e *Error) Attrs() []any {
	return e.attrs
}

func (e *Error) Is(target error) bool {
	if err := AsError(target); err != nil {
		return e.errType == err.errType
	}
	return false
}

func (e *Error) StackTrace() []uintptr {
	return e.stack
}
