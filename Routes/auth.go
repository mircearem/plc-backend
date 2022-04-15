package Routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	db "plc-backend/Db"
	util "plc-backend/Utils"

	"golang.org/x/crypto/bcrypt"
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
	result, dbErr := db.FindUserForSignup(path, user)

	// If db cannot be accessed, return internal server error
	if dbErr != nil {
		resp, _ := json.Marshal(*dbErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	// If username or email is already in the database, return error
	if result {
		result := util.DatabaseError{
			Message: "Username or email already in use",
		}
		json, _ := json.Marshal(result)
		w.WriteHeader(http.StatusConflict)
		w.Write(json)
		return
	}

	// User is not already in the database, insert the user
	// Step 1. Hash the user's password
	password := []byte(user.Password)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	// Check for error while hashing
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Password = string(hash)

	// Step 2. Insert the user into the database
	insertErr := db.InsertUser(path, user)

	if insertErr != nil {
		resp, _ := json.Marshal(insertErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse username and password
	loginRequest := new(util.LoginRequest)
	_ = json.NewDecoder(r.Body).Decode(&loginRequest)

	validationErr := loginRequest.Validate()

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
	usr, dbErr := db.FindUserForLogin(path, loginRequest.Username)

	// If db cannot be accessed, return internal server error
	if dbErr != nil {
		resp, _ := json.Marshal(dbErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	result, _ := json.Marshal(usr)
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	keys, ok := r.URL.Query()["key"]

	// if there are no params in the request, send back error
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := keys[0]

	fmt.Println(username)

	w.WriteHeader(http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
