package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func TokenValid(r *http.Request) error {
	tokenStr := ExtractToken(r)
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}
	return nil
}

func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}
