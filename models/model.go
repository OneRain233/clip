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
