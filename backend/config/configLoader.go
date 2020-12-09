package config

import (
	"encoding/json"
	"fmt"
	"os"
)

//Config struct for config.json
type Config struct {
	Secret           string `json:"secret"`
	ConnectionString string `json:"connection_string"`
	UserPasswordDB   string `json:"user_password_db"`
	DBName           string `json:"dbname"`
	Port             string `json:"port"`
}

//LoadConfiguration load from config.json
func LoadConfiguration() Config {
	var config Config
	configFile, err := os.Open("config/config.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

//GetSecret get secret key for auth
func GetSecret() string {
	config := LoadConfiguration()
	return config.Secret
}

//GetConnectionString get db connection string
func GetConnectionString() string {
	config := LoadConfiguration()
	return config.ConnectionString
}

//GetUserPasswordDB user password for login
func GetUserPasswordDB() string {
	config := LoadConfiguration()
	return config.UserPasswordDB
}

//GetDBName db name
func GetDBName() string {
	config := LoadConfiguration()
	return config.DBName
}

//GetPort port
func GetPort() string {
	config := LoadConfiguration()
	return config.Port
}
