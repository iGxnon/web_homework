package main

import "github.com/gin-gonic/gin"

func main() {
	engine := gin.Default()
	engine.Any("/hi", func(context *gin.Context) {
		context.JSON(200, "Hello World!")
	})
	err := engine.Run(":80")
	if err != nil {
		return
	}
}