package auth

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Username string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required,password"`
	Admin    *bool  `json:"admin" validate:"required"`
}

type errorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	var regex *regexp.Regexp

	regex = regexp.MustCompile(`^[a-z0-9_-]{5,10}$`)

	return regex.MatchString(username)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Variables for validation
	var number, upper, symbol bool
	var len int

	for _, char := range password {
		switch {

		// Check if char is number
		case unicode.IsNumber(char):
			number = true
			len++

		// Check if char is upper case
		case unicode.IsUpper(char):
			upper = true
			len++

		// Check if char is letter
		case unicode.IsLetter(char):
			len++

			// Check if char is special char
		case int(char) >= 33 && int(char) <= 47:
			symbol = true
			len++

		case int(char) >= 58 && int(char) <= 64:
			symbol = true
			len++
		}
	}

	return len >= 8 && number && symbol && upper
}

func validate(usr User) []*errorResponse {
	validate := validator.New()
	var errors []*errorResponse

	// Register validation functions
	validate.RegisterValidation("username", validateUsername)
	validate.RegisterValidation("password", validatePassword)

	err := validate.Struct(usr)

	// Append errors if they exist
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element errorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}
