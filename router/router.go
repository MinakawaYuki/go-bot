package router

import (
	"github.com/gin-gonic/gin"
	"go-bot/controller"
)

func InitRouter() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	api := r.Group("/api")
	{
		api.GET("/event", controller.GetEventMessage)
		//api.GET("/client")
		api.POST("/upload", controller.SavePartyPic)
	}
	r.Run("localhost:10940")
}
