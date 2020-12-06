package handlers

import (
	"encoding/json"
	"log"
	m "ml_website_project/backend/models"

	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

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

//Login login handler
func (env *Env) Login(w http.ResponseWriter, r *http.Request) {
	var user m.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		w.WriteHeader(http.StatusBadRequest)
	}

	encryptedPassword, err := env.User.GetUserPasswordFromDB(user.UserName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else if len(encryptedPassword) <= 0 {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	encryptionErr := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(user.Password))
	if encryptionErr == nil {
		user.Password = encryptedPassword
	} else {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	response, err := generateJWTToken(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(response)
}
