package Controllers

import (
	"encoding/json"
	"net/http"
	"os"
	auth "plc-backend/Auth"
	db "plc-backend/Db"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func UserRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the data from the request
	user := new(auth.User)
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
			result := db.DatabaseError{
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

// Wrap sessions handler around request handler
func UserLogin(s *auth.Sessions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Parse username and password from the login request
		request := new(auth.LoginRequest)
		_ = json.NewDecoder(r.Body).Decode(&request)

		// Validate the data received from the user
		validationError := request.Validate()

		if validationError != nil {
			resp, _ := json.Marshal(validationError)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}

		// Data received from the user has a valid format proceed

		// Set the path to the database and crete request object
		pwd, _ := os.Getwd()
		path := pwd + "/app.db"

		user := new(auth.User)
		user.Username = request.Username

		// Query the database
		usr, err := db.FindUser(path, user)

		// If db cannot be accessed, return internal server error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// Database was queried successfully, check the returned object

		// No such user was found, db returned nil pointer
		if usr == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Username was found, check the password
		compare := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(request.Password))

		// Passwords do not match, return status unauthorized
		if compare != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Passwords match, generate the jwt token and set the session
		expirationTime := time.Now().Add(time.Minute * 15)

		claims := &auth.Claims{
			Username: request.Username,
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

		// Add the cookie to the store
		s.Set(tokenString, request.Username)

		// Could add a log here to save to a log file that the user has logged in
		// timestamp and log

		// Set the cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "session",
			Value:   tokenString,
			Expires: expirationTime,
		})
	})
}

// Wrap sessions handler around request handler
func UserLogout(s *auth.Sessions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Not checking error for cookie
		// If the request got past the jwt verification middleware it means
		// that the cookie is present and valid
		cookie, _ := r.Cookie("session")
		tokenString := cookie.Value

		// Does it make any sense to check the user, if cookie passes verification
		// this means that the token has to be in the map

		// Remove the cookie
		s.Remove(tokenString)

		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		})
	})
}
