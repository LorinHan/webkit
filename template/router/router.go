package router

import (
	"github.com/gin-gonic/gin"
	"webkit/controller"
)

func Init(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	{
		api.GET("/hello", controller.Hello.SayHi)
		api.POST("/validator", controller.Hello.TestValidator)
	}
}
