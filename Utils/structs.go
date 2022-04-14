package Utils

import "github.com/go-playground/validator/v10"

type updateRequest struct {
	jwt   string
	param string
}

type errorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (request *LoginRequest) Validate() []*errorResponse {
	validate := validator.New()
	var errors []*errorResponse

	err := validate.Struct(request)

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
