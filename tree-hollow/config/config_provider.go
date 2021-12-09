package config

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	DefaultDBName    string `json:"default_db_name"`
	DefaultIPAndPort string `json:"default_ip_and_port"`
	DefaultRoot      string `json:"default_root"`
	DefaultPassword  string `json:"default_password"`
	DefaultCharset   string `json:"default_charset"`
	JWTKey           string `json:"jwt_key"`
}

var Config Configuration

func InitializeConfiguration() {
	bytes, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &Config)
	if err != nil {
		panic(err)
	}
}
