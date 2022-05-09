package auth

import (
	"net/http"
	"os"
	util "plc-backend/Utils"

	"github.com/golang-jwt/jwt"
)

func JwtVerify(next http.Handler) http.Handler {
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

		claims := &util.Claims{}

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

func IsAdmin() {
	// Claims contains a property called admin, this has to be true
}
