package errors

type ErrorType[T any] struct {
	code int
	name string
}

func Type[T any](code int, name string) ErrorType[T] {
	return ErrorType[T]{
		code: code,
		name: name,
	}
}

// New は message と details を持つ業務エラーを生成する。
func (t ErrorType[T]) New(message string, details ...T) *Error[T] {
	return &Error[T]{
		errType: t,
		message: message,
		details: details,
		// skip=3: runtime.Callers → captureStack → Type.New の次（New の呼び出し元）から記録
		stack: captureStack(3),
	}
}

// Wrap は cause を包み、呼び出し元の ErrorType を付与する。
// HTTP ステータス・TypeName・Type.Is の判定は外側の Type が優先される（cause の TypeCode は httperror では使われない）。
func (t ErrorType[T]) Wrap(cause error, message string, details ...T) *Error[T] {
	return &Error[T]{
		errType: t,
		cause:   cause,
		message: message,
		details: details,
		// skip=3: runtime.Callers → captureStack → Type.Wrap の次（Wrap の呼び出し元）から記録
		stack: captureStack(3),
	}
}

func (t ErrorType[T]) Is(err error) bool {
	if e := AsError[T](err); e != nil {
		return t.equals(e.errType)
	}
	return false
}

func (t ErrorType[T]) Name() string {
	return t.name
}

func (t ErrorType[T]) equals(target ErrorType[T]) bool {
	return t.code == target.code && t.name == target.name
}
