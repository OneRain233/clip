package controllers

import (
	"clipboard/db"
	"clipboard/forms"
	"clipboard/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

//var db *gorm.DB

func AddClipBoard(c *gin.Context) {
	var content string
	var deviceId string
	var deviceType string

	content, _ = c.GetPostForm("content")
	deviceId, _ = c.GetPostForm("device_id")
	deviceType, _ = c.GetPostForm("device_type")
	//content = utils.GetBase64([]byte(content))
	messageEntity := models.TCPMessage{
		DeviceId:   deviceId,
		DeviceType: deviceType,
		Timestamp:  time.Now().Unix(),
		Data:       content,
	}
	message, err := json.Marshal(messageEntity)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	models.WriteChan <- message
}

func GetLatestClipBoard(c *gin.Context) {
	content, err := db.GetLatestClipBoard()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var resp forms.GetLatestClipBoardResponse
	resp.Code = 0
	resp.Data = content
	c.JSON(200, resp)
	return
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
