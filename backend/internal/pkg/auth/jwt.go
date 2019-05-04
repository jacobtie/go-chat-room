package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = "inbrightestday"

func generateJWT(user string) (string, error) {
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

func AuthMiddle(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {

		} else {
			http.Error(w, "Unauthorized", 401)
		}
	}
}
