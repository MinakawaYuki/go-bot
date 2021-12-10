package router

import (
	"github.com/gin-gonic/gin"
	"go-bot/controller"
)

func InitRouter() {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/event", controller.GetEventMessage)
	}
	r.Run("localhost:10940")
}
