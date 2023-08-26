package main

import (
	"clipboard/config"
	"clipboard/db"
	"clipboard/server"
	"flag"
	"log"
)

func main() {
	mode := flag.String("mode", "HTTP", "TCP or HTTP or BOTH")
	flag.Parse()

	log.Default().Print("Using config file: ", config.GetConfig().GetString("db.filepath"))
	_, err := db.InitDb(config.GetConfig().GetString("db.filepath"))
	if err != nil {
		panic(err)
	}

	log.Default().Print("Using mode: ", *mode)

	switch *mode {
	case "TCP":
		go func() {
			server.RunTcp()
		}()
	case "HTTP":
		go func() {
			server.RunWeb()
		}()
	case "BOTH":
		go func() {
			server.RunTcp()
		}()
		go func() {
			server.RunWeb()
		}()
	}
	select {}
}
