package models

var WriteChan chan []byte

func init() {
	WriteChan = make(chan []byte)
}

func GetWriteChan() chan []byte {
	return WriteChan
}
