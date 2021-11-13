package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

var PaperMap = map[string]*PaperDao{}

type PaperDao struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func LoadAllPapers() {
	var str = ""
	// 拼接起来?
	var errStr = ""
	b := make([]byte, 1024)
	file, err := os.Open("/Users/igxnon/个人项目/Golang/web_homework/homework6/level01AndLevel02/paper.data")
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
		dao := PaperDao{}
		err := json.Unmarshal([]byte(entry), &dao)
		if err != nil {
			errStr += err.Error()
		}
		PaperMap[dao.Path] = &dao
	}
	err3 := file.Close()
	if err3 != nil {
		errStr += err.Error()
	}
	if errStr != "" {
		fmt.Println(errStr)
	}
}
