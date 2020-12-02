package models

import (
	"database/sql"

	"fmt"
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
		panic(err.Error())
	}
	defer row.Close()

	for row.Next() {
		err = row.Scan(&user.UserID, &user.UserName, &user.Password)

		if err != nil {
			panic(err.Error())
		}
	}
	return user, nil
}

//InsertUserToDB Insert one user to Db
func (m UserModel) InsertUserToDB(user User) (string, error) {

	row, err := m.DB.Query("SELECT MAX(UserID) FROM MsUser")

	if err != nil {
		panic(err.Error())
	}
	defer row.Close()

	var maxUserID string
	for row.Next() {
		err = row.Scan(&maxUserID)

		if err != nil {
			panic(err.Error())
		}
	}
	maxUserID = maxUserID[1:4]
	nextUserIDInt, err := strconv.Atoi(maxUserID)
	newUserID := "U"
	nextUserIDInt = nextUserIDInt + 1
	maxUserID = strconv.Itoa(nextUserIDInt)

	for i := 0; i < 3-len(maxUserID); i++ {
		newUserID = newUserID + "0"
	}
	newUserID = newUserID + maxUserID

	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	hashedPass := string(bytes)

	insert, err := m.DB.Query("INSERT INTO MsUser VALUES('" + newUserID + "','" + user.UserName + "','" + hashedPass + "')")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	return "Successfully inserted data with User ID: " + newUserID, nil
}

//UpdateUserToDB Update specific user to db
func (m UserModel) UpdateUserToDB(user User) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	hashedPass := string(bytes)

	update, err := m.DB.Query("UPDATE MsUser SET Username = '" + user.UserName + "', Password = '" + hashedPass + "' WHERE UserID = '" + user.UserID + "'")

	if err != nil {
		panic(err.Error())
	}
	defer update.Close()

	return "Successfully updated data with User ID: " + user.UserID, nil
}

//DeleteUserFromDB Hard deletes a user from db
func (m UserModel) DeleteUserFromDB(userID string) (string, error) {

	delete, err := m.DB.Query("DELETE FROM MsUser WHERE UserID = '" + userID + "'")

	fmt.Println(delete)
	if err != nil {
		panic(err.Error())
	}
	defer delete.Close()

	return "Successfully deleted user with User ID: " + userID, nil
}
