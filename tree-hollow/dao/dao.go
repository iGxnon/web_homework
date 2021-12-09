package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"tree-hollow/config"
)

var dB *sql.DB

func InitializeDefault() {
	configuration := config.Config
	Initialize(
		configuration.DefaultDBName,
		configuration.DefaultRoot,
		configuration.DefaultPassword,
		configuration.DefaultIPAndPort,
		configuration.DefaultCharset,
	)
}

func Initialize(dbName, root, pwd, ipAndPort, charset string) {
	daraSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True", root, pwd, ipAndPort, dbName, charset)
	db, err := sql.Open("mysql", daraSourceName)
	if err != nil {
		log.Fatal(err)
	}
	dB = db
}
