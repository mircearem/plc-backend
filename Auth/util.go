package auth

type User struct {
	username string `json:"username" validate:"required"`
	password string `json:"password" validate:"required"`
	admin    *bool  `json:"required" validate:"required"`
}

type errorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func validate(usr *User) []*errorResponse {

}
