package errors

import (
	"errors"
	"runtime"
)

func AsError[T any](err error) *Error[T] {
	if err == nil {
		return nil
	}
	var e *Error[T]
	if errors.As(err, &e) {
		return e
	}
	return nil
}

func captureStack(skip int) []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	return pcs[0:n]
}
