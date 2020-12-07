package main

import (
	"fmt"
	"mime"
	c "ml_website_project/backend/config"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func enforceJSONHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt != "application/json" {
				http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := r.Header.Get("Authorization")
		split := strings.Split(req, "Bearer ")
		req = split[1]

		token, err := jwt.Parse(req, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(c.GetSecret()), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Signature Invalid", http.StatusUnauthorized)
			}
			http.Error(w, "Something went wrong", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized token", http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})

}
