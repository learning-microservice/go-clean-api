package validate

import (
	"fmt"
	"strings"
)

type Errors []error

func (e *Errors) Error() string {
	if len(*e) == 0 {
		return "unknown error"
	}
	var builder strings.Builder
	for i := range *e {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString((*e)[i].Error())
	}
	return builder.String()
}

func (e *Errors) Unwrap() []error {
	if len(*e) == 0 {
		var zero []error
		return zero
	}
	return *e
}

type FieldError struct {
	field   string
	message string
}

func (e *FieldError) Field() string {
	return e.field
}

func (e *FieldError) Message() string {
	return e.message
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("field %s: %s", e.field, e.message)
}
