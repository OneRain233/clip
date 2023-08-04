package server

import (
	"clipboard/config"
	"log"
	"net"
)

func RunTcp() {
	port := config.GetConfig().GetString("tcp.port")
	if port == "" {
		port = "8081"
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Default().Print("Listening on port: ", port)

	for {
		_, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		//go handleConnection(conn)
	}
}
