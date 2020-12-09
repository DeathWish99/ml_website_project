package main

import (
	"database/sql"
	"fmt"
	"mime"
	c "ml_website_project/backend/config"
	h "ml_website_project/backend/handlers"
	m "ml_website_project/backend/models"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	passwordDB := c.GetUserPasswordDB()
	connectionString := c.GetConnectionString()
	dbName := c.GetDBName()
	port := c.GetPort()
	db, err := sql.Open("mysql", passwordDB+"@tcp("+connectionString+")/"+dbName+"")

	fmt.Println(passwordDB + "@tcp(" + connectionString + ")/" + dbName)
	if err != nil {
		panic(err.Error())
	}

	env := &h.Env{
		User: m.UserModel{DB: db},
		ML:   m.MLModel{DB: db},
	}

	api := r.PathPrefix("/api").Subrouter()
	api.Use(enforceJSONHandler, authorize)

	user := r.PathPrefix("/user").Subrouter()
	user.Use(enforceJSONHandler)

	api.HandleFunc("/getUser/{userid}", env.GetUser).Methods("GET")
	user.HandleFunc("/createUser", env.CreateUser).Methods("POST")
	api.HandleFunc("/updateUser", env.UpdateUser).Methods("PUT")
	api.HandleFunc("/deleteUser/{userid}", env.DeleteUser).Methods("DELETE")
	user.HandleFunc("/login", env.Login).Methods("POST")
	api.HandleFunc("/getAllModelsByUser/{userid}", env.GetAllModelsByUser).Methods("GET")
	api.HandleFunc("/createModel", env.CreateModel).Methods("POST")
	api.HandleFunc("/updateModelDescription/{modelid}/{modeldescription}", env.UpdateModelDescription).Methods("PUT")
	api.HandleFunc("/deleteModel/{modelid}", env.DeleteModel).Methods("DELETE")
	api.HandleFunc("/insertTrainingData", env.InsertTrainingData).Methods("POST")
	user.HandleFunc("/token", env.GetNewToken).Methods("GET")

	http.ListenAndServe(port, r)
}

//Middlewares

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
		if len(r.Header.Get("Authorization")) == 0 {
			http.Error(w, "Request must have bearer token", http.StatusBadRequest)
			return
		}
		req := r.Header.Get("Authorization")
		split := strings.Split(req, "Bearer ")
		req = split[1]
		secret := c.GetSecret()

		token, err := jwt.Parse(req, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(secret), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Signature Invalid", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Unauthorized Token", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})

}
