package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
	"tree-hollow/utils"
)

func auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")
		if err != nil {
			utils.RespInternalError(ctx)
			ctx.Abort()
			return
		}
		refreshToken, err := ctx.Cookie("refresh_token")
		if err != nil {
			utils.RespInternalError(ctx)
			ctx.Abort()
			return
		}
		flagToken, err := utils.AuthorizeJWT(token)
		if err != nil {
			utils.RespInternalError(ctx)
			ctx.Abort()
			return
		}
		flagRefreshToken, err := utils.AuthorizeJWT(refreshToken)
		if err != nil {
			utils.RespInternalError(ctx)
			ctx.Abort()
			return
		}
		if !flagRefreshToken && !flagToken {
			utils.RespErrorWithDate(ctx, "认证失效，需要重新登录!")
			ctx.Abort()
			return
		}

		payload := strings.Split(token, ".")[1]
		var claim = utils.Claims{}
		bytes, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			utils.RespInternalError(ctx)
			ctx.Abort()
			return
		}
		err = json.Unmarshal(bytes, &claim)
		if err != nil {
			utils.RespInternalError(ctx)
			ctx.Abort()
			return
		}

		username := claim.UserName
		ctx.Set("username", username)

		// token 过期了，依据未过期的 refreshToken 生成新的 token 对
		if !flagToken {

			token, refreshToken, err := utils.GenerateTokenPairWithUserName(username)
			if err != nil {
				utils.RespInternalError(ctx)
				ctx.Abort()
				return
			}
			ctx.SetCookie("token", token, 0, "/", "", false, true)
			ctx.SetCookie("refresh_token", refreshToken, 0, "/", "", false, true)

		}
	}
}
