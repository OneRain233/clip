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
