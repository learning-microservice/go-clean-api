package errors

type Type[T any] string

func (t Type[T]) New(message string, details ...T) *Error[T] {
	return &Error[T]{
		errType: t,
		message: message,
		details: details,
		stack:   captureStack(4),
	}
}

func (t Type[T]) Wrap(cause error, message string, details ...T) *Error[T] {
	return &Error[T]{
		errType: t,
		cause:   cause,
		message: message,
		details: details,
		stack:   captureStack(4),
	}
}

func (t Type[T]) Is(err error) bool {
	if e := AsError[T](err); e != nil {
		return e.errType == t
	}
	return false
}

func (t Type[T]) Name() string {
	return string(t)
}
