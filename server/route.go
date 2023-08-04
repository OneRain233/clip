package server

import (
	"clipboard/controllers"
	"github.com/gin-gonic/gin"
)

func RunWeb() {
	route := gin.Default()
	route.POST("/clipboard/add", controllers.AddClipBoard)
	route.GET("/clipboard/list", controllers.GetClipBoardList)

	route.Run(":8080")
}
