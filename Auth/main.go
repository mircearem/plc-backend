package auth

import (
	"encoding/json"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the data from the request
	user := new(User)
	_ = json.NewDecoder(r.Body).Decode(&user)

	// Validate user data
	validationErr := validate(*user)

	// Send response if data sent by the user does not pass validation
	if validationErr != nil {
		resp, _ := json.Marshal(validationErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	// Check to see wether the username provided is already taken
	// resp, _ := json.Marshal(*user)
	w.WriteHeader(http.StatusOK)
	// w.Write(resp)
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
