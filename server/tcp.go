package server

import (
	"clipboard/config"
	"clipboard/db"
	"clipboard/models"
	"log"
	"net"
	"strings"
)

var clients []models.Client

func handleConnection(client *models.Client) {
	conn := client.Conn
	defer conn.Close()
	// send to client
	go func() {
		for {
			select {
			case message := <-client.Write:
				conn.Write(message)
			}
		}
	}()

	// receive from client
	go func() {
		for {
			message := make([]byte, 1024)
			_, err := conn.Read(message)
			if err != nil {
				log.Default().Println("Connection error: ", conn.RemoteAddr(), err)
				return
			}
			models.WriteChan <- message
		}
	}()

	// wait goroutine
	select {}
}

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
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Default().Println("Receive connection from " + conn.RemoteAddr().String())
		client := models.Client{Conn: conn, Write: make(chan []byte), Read: make(chan []byte)}
		clients = append(clients, client)
		go handleConnection(&client)
	}
}

func init() {
	//WriteChan = make(chan []byte)
	go func() {
		for {
			select {
			case message := <-models.WriteChan:
				// send the message to all clients
				s := string(message)
				log.Default().Println("Receive message: ", strings.TrimSpace(s))
				err := db.AddClipBoard(s)
				if err != nil {
					log.Fatal(err)
				}
				for _, client := range clients {
					client.Write <- message
				}
			}
		}
	}()
}
