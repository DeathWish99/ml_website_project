package models

import (
	"database/sql"
	"log"

	"strconv"

	"golang.org/x/crypto/bcrypt"
)

//User Model for table MsUser
type User struct {
	UserID   string `json:"userid"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

//UserModel Db connection
type UserModel struct {
	DB *sql.DB
}

//GetUserFromDB Gets specific user from DB
func (m UserModel) GetUserFromDB(userID string) (User, error) {
	row, err := m.DB.Query("SELECT * FROM MsUser WHERE UserID = '" + userID + "'")

	var user User
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	for row.Next() {
		err = row.Scan(&user.UserID, &user.UserName, &user.Password)

		if err != nil {
			log.Fatal(err)
		}
	}
	return user, nil
}

func (m UserModel) GetUserIDFromDB(username string, password string) (string, error) {
	stmt, err := m.DB.Prepare("SELECT UserID FROM MsUser WHERE Username = '?' AND Password = '?'")
	if err != nil {
		log.Fatal(err)
	}

	get := stmt.QueryRow(username, password)

	var userID string

	err = get.Scan(&userID)

	switch err {
	case nil:
		log.Fatal(err)
	case sql.ErrNoRows:
		return "", nil

	}

	return userID, nil
}

//InsertUserToDB Insert one user to Db
func (m UserModel) InsertUserToDB(user User) (string, error) {

	var maxUserID string
	var newUserID string
	err := m.DB.QueryRow("SELECT MAX(UserID) FROM MsUser").Scan(&maxUserID)

	switch {
	case err == sql.ErrNoRows:
		newUserID = "U001"
	case err != nil:
		log.Fatal(err)
	default:
		maxUserID = maxUserID[1:4]
		nextUserIDInt, err := strconv.Atoi(maxUserID)
		if err != nil {
			log.Fatal(err)
		}
		newUserID = "U"
		nextUserIDInt = nextUserIDInt + 1
		maxUserID = strconv.Itoa(nextUserIDInt)

		for i := 0; i < 3-len(maxUserID); i++ {
			newUserID = newUserID + "0"
		}
		newUserID = newUserID + maxUserID

	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	hashedPass := string(bytes)

	insert, err := m.DB.Exec("INSERT INTO MsUser VALUES('" + newUserID + "','" + user.UserName + "','" + hashedPass + "')")

	if err != nil {
		log.Fatal(err)
	}

	rows, err := insert.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		return "", nil
	}

	return newUserID, nil
}

//UpdateUserToDB Update specific user to db
func (m UserModel) UpdateUserToDB(user User) (bool, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	hashedPass := string(bytes)

	update, err := m.DB.Exec("UPDATE MsUser SET Username = '" + user.UserName + "', Password = '" + hashedPass + "' WHERE UserID = '" + user.UserID + "'")

	if err != nil {
		panic(err.Error())
	}

	rows, err := update.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		return false, nil
	}

	return true, nil
}

//DeleteUserFromDB Hard deletes a user from db
func (m UserModel) DeleteUserFromDB(userID string) (bool, error) {
	deleteModel, err := m.DB.Exec("DELETE FROM TrModels WHERE UserID = '" + userID + "'")
	if err != nil {
		log.Fatal(err)
	}
	rowsModel, err := deleteModel.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rowsModel != 1 {

	}

	delete, err := m.DB.Exec("DELETE FROM MsUser WHERE UserID = '" + userID + "'")

	if err != nil {
		log.Fatal(err)
	}

	rows, err := delete.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		return false, nil
	}

	return true, nil
}
