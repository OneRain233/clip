package models

import "net"

type ClipBoardEntity struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Hash    string `json:"hash"`
	Time    string `json:"time"`
}

type Client struct {
	Conn  net.Conn
	Write chan []byte
	Read  chan []byte
}

type TCPMessage struct {
	DeviceId   string `json:"device_id"`
	DeviceType string `json:"device_type"`
	Timestamp  int64  `json:"timestamp"`
	Data       string `json:"data"`
}
