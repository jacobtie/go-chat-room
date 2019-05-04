// Modified from tutorial point

package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("inbrightestday")

// GenerateJWT generates and returns a JWT
func GenerateJWT(user string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = user
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// MustAuth is middleware that will force an authentication
func MustAuth(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Authenticating")
		cookie, err := r.Cookie("jwt")
		if err != nil {
			log.Printf("No cookie found")
			http.Error(w, "Unauthorized", 401)
			return
		}
		tokenVal := cookie.Value
		if tokenVal == "" {
			log.Printf("Cookie value is empty")
			http.Error(w, "Unauthorized", 401)
			return
		}
		token, err := jwt.Parse(tokenVal, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Incorrect JWT")
			}
			return signingKey, nil
		})
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, "Unauthorized", 401)
			return
		}
		if token.Valid {
			log.Printf("Authentication successful")
			fn(w, r)
		} else {
			log.Printf("Invalid token")
			http.Error(w, "Invalid authentication", 401)
			return
		}
	}
}
