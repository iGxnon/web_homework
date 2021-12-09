package api

import "github.com/gin-gonic/gin"

func RegisterRouter() {
	engine := gin.Default()

	engine.POST("/login", login)
	engine.POST("/register", register)

	engine.Run(":80")
}
