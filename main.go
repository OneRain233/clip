package main

import (
	"clipboard/config"
	"clipboard/db"
	"clipboard/server"
	"log"
)

func main() {
	log.Default().Print("Using config file: ", config.GetConfig().GetString("db.filepath"))
	_, err := db.InitDb(config.GetConfig().GetString("db.filepath"))
	if err != nil {
		panic(err)
	}

	go func() {
		server.RunTcp()
	}()

	go func() {
		server.RunWeb()
	}()
	select {}
}

//func main() {
//	r := gin.Default()
//	r.POST("/", func(c *gin.Context) {
//		wechat := c.PostForm("wechat")
//		c.String(200, wechat)
//	})
//
//	r.Run(":8080")
//}
