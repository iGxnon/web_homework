package main

import (
	"tree-hollow/api"
	"tree-hollow/config"
	"tree-hollow/dao"
)

func main() {
	config.InitializeConfiguration()
	dao.InitializeDefault()
	api.RegisterRouter()
}
