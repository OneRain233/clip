package server

import (
	"clipboard/config"
	"clipboard/controllers"
	"github.com/gin-gonic/gin"
)

func RunWeb() {
	route := gin.Default()

	route.POST("/clipboard/add", controllers.AddClipBoard)
	route.GET("/clipboard/list", controllers.GetClipBoardList)
	route.GET("/clipboard/latest", controllers.GetLatestClipBoard)

	port := config.GetConfig().GetString("web.port")
	if port == "" {
		port = "8080"
	}
	route.Run(":" + port)
}
