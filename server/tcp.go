package server

import (
	"clipboard/config"
	"clipboard/db"
	"clipboard/models"
	"encoding/json"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var clients []models.Client

func handleConnection(client *models.Client) {

	var wg sync.WaitGroup

	conn := client.Conn
	defer conn.Close()

	latest, err := db.GetLatestClipBoard()
	if err != nil {
		latest = models.ClipBoardEntity{Content: ""}
	}
	// from time string
	timestamp, _ := time.Parse("2006-01-02 15:04:05", latest.Time)

	latestMessage := models.TCPMessage{
		DeviceId:   "server",
		DeviceType: "server",
		Timestamp:  timestamp.Unix(),
		Data:       latest.Content,
	}
	latestMessageJson, _ := json.Marshal(latestMessage)
	conn.Write(latestMessageJson)

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case message := <-client.Write:
				conn.Write(message)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			message := make([]byte, 1024)
			n, err := conn.Read(message)
			if n == 0 || err != nil {
				return
			}
			message = message[:n]
			models.WriteChan <- message
		}
	}()

	wg.Wait()

	// remove client
	for i, c := range clients {
		if c.Conn == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
		log.Default().Println("Remove client: ", conn.RemoteAddr().String())
	}
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

	go HandleClients()
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

func HandleClients() {
	for {
		select {
		case message := <-models.WriteChan: // json format message
			// send the message to all clients
			s := string(message)
			log.Default().Println("Receive message: ", strings.TrimSpace(s), " from client")
			var messageEntity models.TCPMessage
			err := json.Unmarshal(message, &messageEntity)
			if err != nil {
				log.Default().Println("Unmarshal message error: ", err)
			}

			err = db.AddClipBoard(messageEntity.Data, messageEntity.Timestamp)
			if err != nil {
				log.Fatal(err)
			}
			for _, client := range clients {
				if client.Conn == nil {
					continue
				}
				log.Default().Println("Sending to client: ", client.Conn.RemoteAddr().String())
				messageEntityJson, _ := json.Marshal(messageEntity)
				client.Write <- messageEntityJson
			}
		}
	}
}
