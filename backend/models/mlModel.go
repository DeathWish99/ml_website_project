package models

import (
	"database/sql"
	"log"
)

type ML struct {
	ModelID          int    `json:"modelid"`
	UserID           string `json:"userid"`
	ModelDescription string `json:"modeldescription"`
	CategoryID       int    `json:"categoryid"`
	FolderName       string `json:"foldername"`
	Downloadable     bool   `json:"downloadable"`
	CategoryName     string `json:"categoryname"`
}

type Training struct {
	TrainingDataID int    `json:"trainingdataid"`
	CategoryID     int    `json:"categoryid"`
	FolderName     string `json:"foldername"`
	DataName       string `json:"dataname"`
}

type MLModel struct {
	DB *sql.DB
}

func (m MLModel) GetAllModelsByUserFromDB(userID string) ([]ML, error) {
	rows, err := m.DB.Query("SELECT ModelID, UserID, ModelDescription, CategoryID, CategoryName, FolderName, Downloadable FROM TrModels a JOIN LtCategory b ON a.CategoryID = b.CategoryID WHERE a.UserID = '" + userID + "' AND Downloadable = 1")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var mls []ML
	for rows.Next() {
		var ml ML

		err := rows.Scan(&ml.ModelID, &ml.UserID, &ml.ModelDescription, &ml.CategoryID, &ml.CategoryName, &ml.FolderName, &ml.Downloadable)
		if err != nil {
			log.Fatal(err)
		}

		mls = append(mls, ml)
	}

	return mls, nil
}

func (m MLModel) InsertModelToDB(ml ML) (bool, error) {
	stmt, err := m.DB.Prepare("INSERT INTO TrModels VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var maxModelID int
	var newModelID int
	err1 := m.DB.QueryRow("SELECT MAX(ModelID) FROM TrModels").Scan(&maxModelID)

	switch {
	case err1 == sql.ErrNoRows:
		newModelID = 1
	case err != nil:
		log.Fatal(err1)
	default:
		newModelID = maxModelID + 1
	}

	insert, err := stmt.Exec(newModelID, ml.UserID, ml.ModelDescription, ml.CategoryID, ml.FolderName, ml.Downloadable)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := insert.RowsAffected()

	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		return false, nil
	}

	return true, nil
}

func (m MLModel) UpdateModelDescriptionToDB(modelID int, modelDescription string) (bool, error) {
	stmt, err := m.DB.Prepare("UPDATE TrModels SET ModelDescription = ? WHERE ModelID = ?")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	update, err := stmt.Exec(modelDescription, modelID)
	if err != nil {
		log.Fatal(err)
	}

	row, err := update.RowsAffected()

	if err != nil {
		log.Fatal(err)
	}
	if row != 1 {
		return false, nil
	}
	return true, nil
}

func (m MLModel) DeleteModelFromDB(modelID int) (bool, error) {
	stmt, err := m.DB.Prepare("DELETE FROM TrModels WHERE ModelID = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	delete, err := stmt.Exec(modelID)
	if err != nil {
		log.Fatal(err)
	}

	row, err := delete.RowsAffected()

	if err != nil {
		log.Fatal(err)
	}
	if row != 1 {
		return false, nil
	}
	return true, nil
}

func (m MLModel) InsertTrainingDataToDB(tr Training) (bool, error) {
	stmt, err := m.DB.Prepare("INSERT INTO TrTrainingData VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	insert, err := stmt.Exec(tr.CategoryID, tr.FolderName, tr.DataName)

	row, err := insert.RowsAffected()

	if err != nil {
		log.Fatal(err)
	}
	if row != 1 {
		return false, nil
	}
	return true, nil
}
