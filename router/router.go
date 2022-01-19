package router

import (
	"github.com/gin-gonic/gin"
	"go-bot/controller"
	"go-bot/utils/tools"
	"net/http"
)

func InitRouter() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	r.StaticFS("/image", http.Dir(tools.GetImageFullPath()))
	api := r.Group("/api")
	{
		api.GET("/event", controller.GetEventMessage)
		//api.GET("/client")
		api.POST("/upload", controller.SavePartyPic)
	}
	r.Run("localhost:10940")
}
