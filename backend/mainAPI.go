package main

import (
	"database/sql"
	"log"
	h "ml_website_project/backend/handlers"
	m "ml_website_project/backend/models"
	"net/http"

	"github.com/gorilla/mux"

	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

func main() {

	manager := manage.NewDefaultManager()

	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

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
	r.HandleFunc("/api/getAllModelsByUser/{userid}", env.GetAllModelsByUser).Methods("GET")
	r.HandleFunc("/api/createModel", env.CreateModel).Methods("POST")
	r.HandleFunc("/api/updateModelDescription/{modelid}/{modeldescription}", env.UpdateModelDescription).Methods("PUT")
	r.HandleFunc("/api/deleteModel/{modelid}", env.DeleteModel).Methods("DELETE")
	r.HandleFunc("/api/insertTrainingData", env.InsertTrainingData).Methods("POST")
	http.ListenAndServe(":8000", r)
}
