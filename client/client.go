package main

import (
	"clipboard/utils"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

var conn net.Conn
var clipBoardType string

func Connect(host string, port int) error {
	var err error
	conn, err = net.Dial("tcp", "localhost:8081") // "localhost:8081
	if err != nil {
		return err
	}
	return nil
}

func GetClipBoardEnv() {
	envs := os.Environ()
	for _, env := range envs {
		key := env[:strings.Index(env, "=")]
		value := env[strings.Index(env, "=")+1:]
		if key == "XDG_SESSION_TYPE" {
			clipBoardType = value
			break
		}
	}
}

func ReadFromClipBoard() (string, error) {
	// TODO
	return "", nil
}

func WriteToClipBoard(text string) error {
	fmt.Println(text)
	switch clipBoardType {
	case "x11":
		err := exec.Command("echo", text, "|", "xclip", "-selection", "clipboard").Run()
		if err != nil {
			return err
		}
	case "wayland":
		// exec wl-copy
		//cmd := fmt.Sprintf("echo %s | base64 -d | wl-copy", text)
		copyCmd := fmt.Sprintf("echo %s | base64 -d | wl-copy", fmt.Sprintf("%q", strings.TrimSpace(text)))
		fmt.Println(copyCmd)
		err := exec.Command("/bin/sh", "-c", copyCmd).Run()
		if err != nil {
			return err
		}
	default:
		// TODO

	}
	return nil
}

func handleConnection() {
	if conn == nil {
		return
	}
	defer conn.Close()

	// wait for server
	go func() {
		for {
			message := make([]byte, 1024)
			l, err := conn.Read(message)
			if err != nil {
				return
			}
			//fmt.Println(message)
			msgStr := string(message[:l])

			msgBase64 := utils.GetBase64([]byte(strings.TrimSpace(msgStr)))
			err = WriteToClipBoard(msgBase64)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	select {}
}

func main() {
	host := flag.String("host", "localhost", "host")
	port := flag.Int("port", 8081, "port")
	flag.Parse()
	for {
		err := Connect(*host, *port)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	GetClipBoardEnv()
	handleConnection()
}
