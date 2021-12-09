package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	DefaultDBName    = "tree_hollows_db"
	DefaultIPAndPort = "localhost:3306"
	DefaultRoot      = "root"
	DefaultPassword  = "502508"
	DefaultCharset   = "utf8mb4"
)

var dB *sql.DB

func InitializeDefault() {
	Initialize(DefaultDBName, DefaultRoot, DefaultPassword, DefaultIPAndPort, DefaultCharset)
}

func Initialize(dbName, root, pwd, ipAndPort, charset string) {
	daraSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True", root, pwd, ipAndPort, dbName, charset)
	db, err := sql.Open("mysql", daraSourceName)
	if err != nil {
		log.Fatal(err)
	}
	dB = db
}
