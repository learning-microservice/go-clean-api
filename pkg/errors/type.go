package errors

type Type string

func (t Type) New(messages ...string) error {
	return &Error{
		errType:  t,
		messages: messages,
		stack:    captureStack(4),
	}
}

func (t Type) Wrap(cause error, messages ...string) error {
	return &Error{
		errType:  t,
		cause:    cause,
		messages: messages,
		stack:    captureStack(4),
	}
}

func (t Type) Is(err error) bool {
	if e := AsError(err); e != nil {
		return e.errType == t
	}
	return false
}

func (t Type) Name() string {
	return string(t)
}
