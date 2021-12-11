package utils

import "github.com/gin-gonic/gin"

func SetSnitchName(ctx *gin.Context, name *string) {
	username, ok := ctx.Get("username")
	if !ok {
		RespErrorWithDate(ctx, "你没有登录!")
		ctx.Abort()
		return
	}
	*name, ok = username.(string)
	if !ok {
		RespInternalError(ctx)
		ctx.Abort()
		return
	}
}
