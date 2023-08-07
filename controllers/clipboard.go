package controllers

import (
	"clipboard/db"
	"clipboard/forms"
	"clipboard/models"
	"github.com/gin-gonic/gin"
)

//var db *gorm.DB

func AddClipBoard(c *gin.Context) {
	var content string
	content, _ = c.GetPostForm("content")
	//content = utils.GetBase64([]byte(content))
	models.WriteChan <- []byte(content)
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
