package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var dB *sql.DB

func Initialize() {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/message_board?charset=utf8mb4&parseTime=True")
	if err != nil {
		panic(err)
	}

	dB = db
}
