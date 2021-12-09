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
	JWTTimeOut       int64  `json:"jwt_time_out"`
}

/*
config.json 模板
{
  "default_db_name": "tree_hollows_db",
  "default_ip_and_port": "localhost:3306",
  "default_root": "root",
  "default_password": "******",
  "default_charset": "utf8mb4",
  "jwt_key": "伟神嘿嘿我的伟神嘿嘿",
  "jwt_time_out": 172800
}
*/

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
