package auth

import (
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func CreateToken(user_id uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = user_id
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func ExtractToken(r *http.Request) string {
	key := r.URL.Query()
	token := key.Get("token")
	if token != "" {
		return token
	}
	bearertoken := r.Header.Get("Authorization")
	if len(strings.Split(bearertoken, " ")) == 2 {
		return strings.Split(bearertoken, " ")[1]
	}
	return ""
}
