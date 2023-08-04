package controllers

import (
	"clipboard/db"
	"clipboard/forms"
	"github.com/gin-gonic/gin"
)

//var db *gorm.DB

func AddClipBoard(c *gin.Context) {
	var req forms.AddClipBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	//db, err := OpenDb("clipboard.db")
	//if err != nil {
	//	c.JSON(500, gin.H{"error": err.Error()})
	//}
	if err := db.AddClipBoard(req.Content); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	var resp forms.AddClipBoardResponse
	resp.Code = 0
	c.JSON(200, resp)
	//return nil
}

func GetClipBoardList(c *gin.Context) {
	var req forms.GetClipBoardListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		list, err := db.GetClipBoard(0, 10)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var resp forms.GetClipBoardListResponse
		resp.Code = 0
		resp.Data = list
		c.JSON(200, resp)
		return
	}

	list, err := db.GetClipBoard(req.Offset, req.Limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var resp forms.GetClipBoardListResponse
	resp.Code = 0
	resp.Data = list
	c.JSON(200, resp)
	return
}
