package handlers

import (
	"encoding/json"
	"log"
	m "ml_website_project/backend/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (env *Env) GetAllModelsByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	mls, err := env.ML.GetAllModelsByUserFromDB(params["userid"])

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	json.NewEncoder(w).Encode(mls)
}

func (env *Env) CreateModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ml m.ML

	err := json.NewDecoder(r.Body).Decode(&ml)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	status, err := env.ML.InsertModelToDB(ml)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if !status {
		http.Error(w, "An unexpected Error occured", http.StatusInternalServerError)
	}
}

func (env *Env) UpdateModelDescription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	modelid, err := strconv.Atoi(params["modelid"])
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	status, err := env.ML.UpdateModelDescriptionToDB(modelid, params["modeldescription"])
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if !status {
		http.Error(w, "Record does not exist", http.StatusInternalServerError)
	}
}

func (env *Env) DeleteModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	modelid, err := strconv.Atoi(params["modelid"])
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	status, err := env.ML.DeleteModelFromDB(modelid)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if !status {
		http.Error(w, "Record does not exist", http.StatusInternalServerError)
	}

}

func (env *Env) InsertTrainingData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tr m.Training

	err := json.NewDecoder(r.Body).Decode(&tr)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	status, err := env.ML.InsertTrainingDataToDB(tr)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if !status {
		http.Error(w, "An unexpected Error occured", http.StatusInternalServerError)
	}
}
