package main

import (
	"tree-hollow/api"
	"tree-hollow/dao"
)

func main() {
	dao.InitializeDefault()
	api.RegisterRouter()
}
