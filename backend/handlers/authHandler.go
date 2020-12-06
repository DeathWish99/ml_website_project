package handlers

import (
	"encoding/json"
	"log"
	m "ml_website_project/backend/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	secret = "?#$!524451XASD"
)

var jwtKey = []byte(secret)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Response struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
	ExpiresAt   string `json:"expires_at"`
}

func (env *Env) GetNewToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user m.User

	err := json.NewDecoder(r.Body).Decode(&user)

	encryptedPassword, err := env.User.GetUserPasswordFromDB(user.UserName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	} else if len(encryptedPassword) <= 0 {
		log.Println(err)
		log.Println(user.UserName)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	encryptionErr := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(user.Password))
	if encryptionErr == nil {
		user.Password = encryptedPassword
	} else {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonStr, err := generateJWTToken(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(jsonStr)
}

func generateJWTToken(user m.User) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute)

	claims := &Claims{
		Username: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println(err)
		return "", err
	}

	expTime := expirationTime.Format("2006-01-02 15:04:05")

	res := &Response{
		AccessToken: tokenString,
		Username:    user.UserName,
		ExpiresAt:   expTime,
	}
	content, err := json.Marshal(res)

	return string(content), nil
}
