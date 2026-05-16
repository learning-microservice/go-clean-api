package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/go-playground/locales/ja"
	translations "github.com/go-playground/validator/v10/translations/ja"

	"go-clean-api/pkg/validate"
)

//go:embed resource.json
var resource []byte

func main() {
	// load resource
	var resourceMap map[string]string
	err := json.Unmarshal(resource, &resourceMap)
	if err != nil {
		panic(err)
	}

	// default locale
	locale := ja.New()

	// setup validator
	validator, err := validate.New(
		validate.WithTranslator(locale, translations.RegisterDefaultTranslations, true),
		validate.WithFieldNameMap(locale.Locale(), resourceMap),
	)
	if err != nil {
		panic(err)
	}

	user := User{
		FirstName: "",
		LastName:  "",
		Age:       -10,
		Email:     "test.taro@example.com",
	}

	if err = validator.ValidateCtx(context.Background(), user); err != nil {
		response := struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Details error  `json:"details,omitempty"`
		}{
			Code:    400,
			Message: "invalid parameters",
			Details: err,
		}
		data, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(data))
	}

	fmt.Println("validation successful")
}

// User contains user information
type User struct {
	FirstName      string     `validate:"required"`
	LastName       string     `validate:"required"`
	Age            int        `validate:"gte=0,lte=130"`
	Email          string     `validate:"required,email"`
	FavouriteColor string     `validate:"iscolor"`                // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}
