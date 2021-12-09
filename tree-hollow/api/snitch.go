package api

import (
	"github.com/gin-gonic/gin"
	"tree-hollow/model"
	"tree-hollow/service"
	"tree-hollow/utils"
)

func login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	checkPassword, err := service.CheckPassword(username, password)
	if err != nil {
		utils.RespErrorWithDate(ctx, err)
		return
	}
	if checkPassword {
		utils.RespSuccessful(ctx)
		// todo
	}
}

func register(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	err := service.RegisterSnitch(model.Snitch{
		Name:     username,
		Password: password,
	})
	if err != nil {
		utils.RespInternalError(ctx)
		return
	}
	utils.RespSuccessful(ctx)
}
