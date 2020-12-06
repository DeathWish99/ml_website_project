package main

import (
	"database/sql"
	h "ml_website_project/backend/handlers"
	m "ml_website_project/backend/models"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	db, err := sql.Open("mysql", "root:Badger-85@tcp(127.0.0.1:3306)/MLDB")

	if err != nil {
		panic(err.Error())
	}

	env := &h.Env{
		User: m.UserModel{DB: db},
		ML:   m.MLModel{DB: db},
	}

	r.HandleFunc("/api/getUser/{userid}", env.GetUser).Methods("GET")
	r.HandleFunc("/api/createUser", env.CreateUser).Methods("POST")
	r.HandleFunc("/api/updateUser", env.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/deleteUser/{userid}", env.DeleteUser).Methods("DELETE")
	r.HandleFunc("/api/login", env.Login).Methods("POST")
	r.HandleFunc("/api/getAllModelsByUser/{userid}", env.GetAllModelsByUser).Methods("GET")
	r.HandleFunc("/api/createModel", env.CreateModel).Methods("POST")
	r.HandleFunc("/api/updateModelDescription/{modelid}/{modeldescription}", env.UpdateModelDescription).Methods("PUT")
	r.HandleFunc("/api/deleteModel/{modelid}", env.DeleteModel).Methods("DELETE")
	r.HandleFunc("/api/insertTrainingData", env.InsertTrainingData).Methods("POST")
	r.HandleFunc("/api/token", env.GetNewToken).Methods("GET")

	http.ListenAndServe(":8000", r)
}
