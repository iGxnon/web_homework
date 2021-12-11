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
	err := service.CheckPassword(username, password)
	if err != nil {
		utils.RespErrorWithDate(ctx, err)
		return
	}

	token, refreshToken, err := utils.GenerateTokenPairWithUserName(username)
	if err != nil {
		utils.RespInternalError(ctx)
		return
	}
	ctx.SetCookie("token", token, 0, "/", "", false, true)
	ctx.SetCookie("refresh_token", refreshToken, 0, "/", "", false, true)

	utils.RespSuccessful(ctx)
}

func register(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	if !utils.CheckPwdSafe(password) {
		utils.RespErrorWithDate(ctx, "密码过于简单!")
		return
	}

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

// 以下为账户服务

func logOffForever(ctx *gin.Context) {
	confirmPwd := ctx.PostForm("confirm_password")
	username, ok := ctx.Get("username")
	if !ok {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
	name, ok := username.(string)
	if !ok {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
	err := service.CheckPassword(name, confirmPwd)
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
	err = service.LogOffForever(name)
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
	utils.RespSuccessfulWithDate(ctx, "注销成功!")
}
