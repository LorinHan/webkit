package router

import (
	"github.com/gin-gonic/gin"
	"webkit/controller"
)

// Init 路由初始化
func Init(engine *gin.Engine) {
	engine.Use(Cors)
	api := engine.Group("/api/v1")
	{
		api.GET("/hello", controller.Hello.SayHi)
		api.POST("/validator", controller.Hello.TestValidator)
	}
}

// Cors 跨域设置
func Cors(c *gin.Context) {
	requestOrigin := c.Request.Header.Get("Origin")
	c.Writer.Header().Set("Access-Control-Allow-Origin", requestOrigin)
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}
