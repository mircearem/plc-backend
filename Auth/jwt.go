// Package that handles jwt authentication and authorization
package auth

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

// Claim to be encoded by jwt library
type Claims struct {
	Username string `json:"username"`
	Admin    string `json:"admin"`
	jwt.StandardClaims
}

// Middleware to determine if the jwt is valid
func CheckJwtValidity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")

		// Check if there is an error parsing cookie or no coockie sent
		// with the request
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Cookie found, extract the token
		tokenString := cookie.Value

		claims := &Claims{}

		// Parse the token
		tkn, err := jwt.ParseWithClaims(tokenString, claims,
			func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)

		// Error parsing token, return
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Token is invalid, return
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Token is valid, go to next function
		next.ServeHTTP(w, r)
	})
}

// Middleware to determine if user is admin or not
func CheckAdminStatus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// No sense to check for error, cookie is clearly there
		cookie, _ := r.Cookie("session")
		tokenString := cookie.Value

		claims := &Claims{}

		// Should be no parsing error, cookie is clearly there
		// otherwise this middleware would not been reached
		jwt.ParseWithClaims(tokenString, claims,
			func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)

		// User is not admin
		if !(claims.Admin == "true") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
