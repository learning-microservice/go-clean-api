package validate

import (
	"fmt"
	"strings"
)

type FieldErrors []FieldError

func (e *FieldErrors) Error() string {
	var builder strings.Builder
	if e != nil {
		for i, err := range *e {
			if i != 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(err.Error())
		}
	}
	return builder.String()
}

type FieldError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("field %s: %s", e.Field, e.Message)
}
