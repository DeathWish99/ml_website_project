package handlers

import m "ml_website_project/backend/models"

//Env Interface for Models
type Env struct {
	User interface {
		GetUserFromDB(userID string) (m.User, error)
		InsertUserToDB(user m.User) (string, error)
		UpdateUserToDB(user m.User) (bool, error)
		DeleteUserFromDB(userID string) (bool, error)
		GetUserIDFromDB(username string, password string) (string, error)
	}
	ML interface {
		GetAllModelsByUserFromDB(userID string) ([]m.ML, error)
		InsertModelToDB(ml m.ML) (bool, error)
		UpdateModelDescriptionToDB(modelID int, modelDescription string) (bool, error)
		DeleteModelFromDB(modelID int) (bool, error)
		InsertTrainingDataToDB(tr m.Training) (bool, error)
	}
}
