package models

type ClipBoardEntity struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Hash    string `json:"hash"`
	Time    string `json:"time"`
}
