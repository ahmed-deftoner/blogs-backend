package middleware

import (
	"errors"
	"net/http"

	"github.com/ahmed-deftoner/blogs-backend/api/auth"
	"github.com/ahmed-deftoner/blogs-backend/api/response"
)

func MiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		next(w, r)
	}
}

func MiddlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			response.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}
		next(w, r)
	}
}
