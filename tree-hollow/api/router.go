package api

import "github.com/gin-gonic/gin"

func RegisterRouter() {
	engine := gin.Default()

	engine.POST("/login", login)
	engine.POST("/register", register)

	accountGroup := engine.Group("/account")
	{
		accountGroup.Use(auth())
		accountGroup.DELETE("/", logOffForever)

	}

	engine.Run(":80")
}
