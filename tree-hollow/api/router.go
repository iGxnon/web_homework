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

	secretGroup := engine.Group("/secret")
	{
		secretGroup.Use(auth())
		secretGroup.GET("/", getSecret)
		secretGroup.POST("/", addSecret)
		secretGroup.PUT("/", updateSecret)
		secretGroup.DELETE("/", deleteSecret)
	}

	commentGroup := engine.Group("/comment")
	{
		commentGroup.Use(auth())
		commentGroup.GET("/", getComment)
		commentGroup.POST("/", addComment)
		commentGroup.PUT("/", updateComment)
		commentGroup.DELETE("/", deleteComment)
	}

	engine.Run(":80")
}
