package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
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
		params := mux.Vars(r)

		log.Println("Authenticating")
		if params["token"] != "" {
			token, err := jwt.Parse(params["token"], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("Incorrect JWT")
				}
				return signingKey, nil
			})
			if err != nil {
				log.Println("Wrong JWT: " + params["token"])
				http.Error(w, "Unauthorized", 401)
				return
			}
			if token.Valid {
				fn(w, r)
			}
		} else {
			log.Println("No JWT")
			http.Error(w, "Unauthorized", 401)
		}
	}
}
