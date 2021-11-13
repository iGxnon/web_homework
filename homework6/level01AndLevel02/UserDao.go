package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

var UserDaoMap = map[string]*UserDao{}

type UserDao struct {
	Name    string `json:"name"`
	PwdLock string `json:"pwd_lock"`
	Gender  Gender `json:"gender"`
	Age     int    `json:"age"`
	NpyName string `json:"npy"`
}

func LoadALlUserDao() {
	var str = ""
	// 拼接起来?
	var errStr = ""
	b := make([]byte, 1024)
	file, err := os.Open("/Users/igxnon/个人项目/Golang/web_homework/homework6/level01AndLevel02/user_data.data")
	if err != nil {
		errStr += err.Error()
	}
	num, err := file.Read(b)
	for err != io.EOF {
		str = str + string(b[:num])
		num, err = file.Read(b)
	}
	data := strings.Split(str, ";")
	for _, entry := range data {
		dao := UserDao{}
		err := json.Unmarshal([]byte(entry), &dao)
		if err != nil {
			errStr += err.Error()
		}
		UserDaoMap[dao.Name] = &dao
	}
	err3 := file.Close()
	if err3 != nil {
		errStr += err.Error()
	}
	if errStr != "" {
		fmt.Println(errStr)
	}
}
