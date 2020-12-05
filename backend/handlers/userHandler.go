package handlers

import (
	"encoding/json"
	"log"
	m "ml_website_project/backend/models"
	"time"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-session/session"
	"github.com/gorilla/mux"
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

//GetUser Send request to model to GetUser
func (env *Env) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	var user m.User

	user, err := env.User.GetUserFromDB(params["userid"])

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

//CreateUser send request to Create new User
func (env *Env) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user m.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := env.User.InsertUserToDB(user)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
}

//UpdateUser send request to Update existing user
func (env *Env) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user m.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := env.User.UpdateUserToDB(user)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !response {
		http.Error(w, "Record does not exist", http.StatusInternalServerError)
	}
}

//DeleteUser send request to hard delete user
func (env *Env) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	response, err := env.User.DeleteUserFromDB(params["userid"])

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !response {
		http.Error(w, "Record does not exist", http.StatusInternalServerError)
	}
}

func (env *Env) Login(w http.ResponseWriter, r *http.Request) {
	var user m.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		w.WriteHeader(http.StatusBadRequest)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(bytes)

	userID, err := env.User.GetUserIDFromDB(user.UserName, user.Password)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else if len(userID) <= 0 {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err = env.User.GetUserFromDB(userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	expirationTime := time.Now().Add(5 * time.Minute)

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
		w.WriteHeader(http.StatusInternalServerError)
	}

	store, err := session.Start(nil, w, r)

	if r.Method == "POST" {
		store.Set(userID, tokenString)
		store.Save()
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
