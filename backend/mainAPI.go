package main

import (
	"database/sql"
	"encoding/json"
	"log"
	m "ml_website_project/backend/models"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Env Interface for Models
type Env struct {
	user interface {
		GetUserFromDB(userID string) (m.User, error)
		InsertUserToDB(user m.User) (string, error)
		UpdateUserToDB(user m.User) (string, error)
		DeleteUserFromDB(userID string) (string, error)
	}
}

func main() {
	r := mux.NewRouter()

	db, err := sql.Open("mysql", "root:Badger-85@tcp(127.0.0.1:3306)/MLDB")

	if err != nil {
		panic(err.Error())
	}

	env := &Env{
		user: m.UserModel{DB: db},
	}

	r.HandleFunc("/api/getUser/{userid}", env.GetUser).Methods("GET")
	r.HandleFunc("/api/createUser", env.CreateUser).Methods("POST")
	r.HandleFunc("/api/updateUser", env.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/deleteUser/{userid}", env.DeleteUser).Methods("DELETE")
	http.ListenAndServe(":8000", r)
}

//GetUser Send request to model to GetUser
func (env *Env) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	var user m.User

	user, err := env.user.GetUserFromDB(params["userid"])

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
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
		return
	}

	response, err := env.user.InsertUserToDB(user)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	json.NewEncoder(w).Encode(response)
}

//UpdateUser send request to Update existing user
func (env *Env) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user m.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	response, err := env.user.UpdateUserToDB(user)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	json.NewEncoder(w).Encode(response)
}

//DeleteUser send request to hard delete user
func (env *Env) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	response, err := env.user.DeleteUserFromDB(params["userid"])

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	json.NewEncoder(w).Encode(response)
}
