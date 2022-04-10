package auth

import (
	"encoding/json"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get the data from the request
	user := new(User)
	_ = json.NewDecoder(r.Body).Decode(&user)

	// server side, only minimal validation, to check for required fields, main validation on client

}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
