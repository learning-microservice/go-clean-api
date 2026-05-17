package main

import (
	"encoding/json"
	goerrors "errors"
	"fmt"

	"go-clean-api/pkg/errors"
)

var (
	TypeValidation = errors.Type[FieldError]("validation error")
	TypeNotFound   = errors.Type[FieldError]("not found")
)

func main() {
	{
		println("--- standard error ---")
		err := TypeNotFound.New("user not found")
		v, _ := json.MarshalIndent(err, "", "  ")
		println(string(v))
	}
	{
		println("--- wrap error ---")
		validationError := goerrors.Join(
			fmt.Errorf("validation error 1"),
			fmt.Errorf("validation error 2"),
			fmt.Errorf("validation error 3"),
		)
		err := TypeValidation.Wrap(validationError, "user not found")
		v, _ := json.MarshalIndent(err, "", "  ")
		println(string(v))
	}
	{
		println("--- error with details ---")
		err := TypeValidation.New("user not found",
			FieldError{field: "email", message: "email is required"},
			FieldError{field: "password", message: "password is required"},
		)
		v, _ := json.MarshalIndent(err, "", "  ")
		println(string(v))

		// wap1
		println("--- wrap error with details ---")
		err = TypeValidation.Wrap(err, "wrap1 user not found",
			FieldError{field: "nowrap", message: "nowrap error"},
		)
		v, _ = json.MarshalIndent(err, "", "  ")
		println(string(v))
	}
}

type FieldError struct {
	field   string
	message string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("field %s: %s", e.field, e.message)
}

func (e *FieldError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Field   string `json:"field,omitempty"`
		Message string `json:"message"`
	}{
		Field:   e.field,
		Message: e.message,
	})
}
