package validate

import (
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
