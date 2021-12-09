package main

import (
	"homework7And8/api"
	"homework7And8/dao"
)

func main() {
	dao.InitializeDefault()
	api.RegisterRouter()
}
