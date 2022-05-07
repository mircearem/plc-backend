package Routes

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	db "plc-backend/Db"
	util "plc-backend/Utils"

	"github.com/golang-jwt/jwt"
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

	// Variables to handle the file descriptor of the data
	var pwd string
	var path string

	// Set the path for the db file
	pwd, _ = os.Getwd()
	path = pwd + "/app.db"

	// Check to see wether the username provided is already taken
	result, dbErr := db.FindUser(path, user)

	// If db cannot be accessed, return internal server error
	if dbErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(dbErr.Error()))
		return
	}

	// If username or email is already in the database, return error
	if result != nil {
		if (result.Username == user.Username) || (result.Email == user.Email) {
			result := util.DatabaseError{
				Message: "Username or email already in use",
			}
			json, _ := json.Marshal(result)
			w.WriteHeader(http.StatusConflict)
			w.Write(json)
			return
		}
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

	// Set the path for the db file
	pwd, _ := os.Getwd()
	path := pwd + "/app.db"

	user := util.User{}
	user.Username = loginRequest.Username

	// Check to see wether the username provided is already taken
	usr, dbErr := db.FindUser(path, &user)

	// If db cannot be accessed, return internal server error
	if dbErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(dbErr.Error()))
		return
	}

	// User not found, query returns nil
	if usr == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Compare password with hashed password from the database
	compare := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(loginRequest.Password))

	// Password is not correct
	if compare != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Password valid, generate jsonwebtoken
	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &util.Claims{
		Username: loginRequest.Username,
		Admin:    usr.Admin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var jwtKey = []byte(os.Getenv("JWT_SECRET"))

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Set the cookies
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

// This route will update user email or password
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

// This route will delete the user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
