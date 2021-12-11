package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"tree-hollow/model"
	"tree-hollow/service"
	"tree-hollow/utils"
)

func setSnitchName(ctx *gin.Context, name *string) {
	username, ok := ctx.Get("username")
	if !ok {
		utils.RespErrorWithDate(ctx, "你没有登录!")
		ctx.Abort()
		return
	}
	*name, ok = username.(string)
	if !ok {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
}

func getSecret(ctx *gin.Context) {
	typeGet := ctx.PostForm("type")

	if typeGet == "brief" {
		getSecretBrief(ctx)
		return
	}

	if typeGet == "details" {
		getSecretDetails(ctx)
		return
	}

	utils.RespInternalError(ctx)
}

func getSecretBrief(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
	brief, err := service.GetSecretBrief(id)
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}

	utils.RespSuccessfulWithDate(ctx, brief)
}

func getSecretDetails(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
	details, err := service.GetSecretDetails(id)
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
	utils.RespSuccessfulWithDate(ctx, details)
}

func addSecret(ctx *gin.Context) {
	var name string
	setSnitchName(ctx, &name)
	content := ctx.PostForm("content")
	isOpen, err := strconv.ParseBool(ctx.PostForm("is_open"))
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
	secret := model.Secret{
		Content:    content,
		SnitchName: name,
		CommentCnt: 0,
		IsOpen:     isOpen,
		PostTime:   time.Now(),
		UpdateTime: time.Now(),
	}
	err = service.AddSecret(secret)
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
	utils.RespSuccessful(ctx)
}

func deleteSecret(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
	var name string
	setSnitchName(ctx, &name)

	ok, err := service.CheckIdMatchName(id, name)
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
	if ok {
		err = service.DeleteSecretFromId(id)
		if err != nil {
			utils.RespInternalError(ctx)
			ctx.Abort()
			return
		}
		utils.RespSuccessful(ctx)
		return
	}
	utils.RespErrorWithDate(ctx, "你不能删除别人的秘密!")
}

func updateSecret(ctx *gin.Context) {
	var name string
	setSnitchName(ctx, &name)
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}

	ok, err := service.CheckIdMatchName(id, name)
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}

	if !ok {
		utils.RespErrorWithDate(ctx, "你不能更新别人的秘密!")
		ctx.Abort()
		return
	}

	content := ctx.PostForm("content")
	isOpen, err := strconv.ParseBool(ctx.PostForm("is_open"))
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}
	cnt, err := service.GetCommentCnt(id)
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}

	postTime, err := service.GetPostTime(id)
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}

	secret := model.Secret{
		Id:         id,
		Content:    content,
		SnitchName: name,
		CommentCnt: cnt,
		IsOpen:     isOpen,
		PostTime:   postTime,
		UpdateTime: time.Now(),
	}

	err = service.UpdateSecret(secret)
	if err != nil {
		utils.RespInternalError(ctx)
		ctx.Abort()
		return
	}

	utils.RespSuccessful(ctx)
}
