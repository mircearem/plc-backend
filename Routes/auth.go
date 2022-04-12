package Routes

import (
	"encoding/json"
	"net/http"
	"os"
	db "plc-backend/DB"
	util "plc-backend/Utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the data from the request
	user := new(util.User)
	_ = json.NewDecoder(r.Body).Decode(&user)

	// Validate user data
	validationErr := user.Validate()

	// Send response if data sent by the user does not pass validation
	if validationErr != nil {
		resp, _ := json.Marshal(validationErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	// Get the path to the db file
	pwd, _ := os.Getwd()
	path := pwd + "/app.db"

	// Check to see wether the username provided is already taken
	usr, dbErr := db.FindUser(path, user)

	// If db cannot be accessed, return internal server error
	if dbErr != nil {
		resp, _ := json.Marshal(dbErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	// If username or email is already in the database, return error
	if usr != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	// User is not already in the database, insert the user
	// Step 1. Hash the user's password

	w.WriteHeader(http.StatusOK)
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
