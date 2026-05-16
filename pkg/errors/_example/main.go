package main

import (
	"encoding/json"
	goerrors "errors"
	"fmt"

	"go-clean-api/internal/domain/errors"
)

func main() {
	{
		err := errors.TypeNotFound.New("user not found")
		v, _ := json.MarshalIndent(err, "", "  ")
		println(string(v))
	}
	{
		validationError := goerrors.Join(
			fmt.Errorf("validation error 1"),
			fmt.Errorf("validation error 2"),
			fmt.Errorf("validation error 3"),
		)
		err := errors.TypeValidation.Wrap(validationError, "user not found")
		v, _ := json.MarshalIndent(err, "", "  ")
		println(string(v))
	}
}
