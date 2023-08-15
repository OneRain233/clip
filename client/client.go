package main

import (
	"clipboard/models"
	"clipboard/utils"
	"encoding/json"
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
var currentClipBoardHash string

const Help = `
Usage:
	clipboard [options]
Options:
	-h, --help
		Show this help message and exit
	--host <host>
		Host to connect to
	--port <port>
		Port to connect to
`

func Connect(host string, port int) error {
	var err error
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
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
	var content string
	switch clipBoardType {
	case "x11":
		out, err := exec.Command("xclip", "-selection", "clipboard", "-o").Output()
		if err != nil {
			return "", err
		}
		content = string(out)
	case "wayland":
		out, err := exec.Command("xclip", "-selection", "clipboard", "-o").Output()
		if err != nil {
			return "", err
		}
		content = string(out)
	default:
		// TODO
		return "", nil
	}
	//fmt.Println("Read content: ", content)
	return content, nil
}

func WriteToClipBoard(text string, useBase64 bool) error {
	log.Default().Println("[+] WriteToClipBoard: ", text)
	switch clipBoardType {
	case "x11":
		err := exec.Command("echo", text, "|", "xclip", "-selection", "clipboard").Run()
		if err != nil {
			return err
		}
	case "wayland":
		if useBase64 {
			copyCmd := fmt.Sprintf("echo %s | base64 -d | wl-copy", fmt.Sprintf("%q", strings.TrimSpace(text)))
			fmt.Println(copyCmd)
			err := exec.Command("/bin/sh", "-c", copyCmd).Run()
			if err != nil {
				return err
			}
		} else {
			copyCmd := fmt.Sprintf("echo %s | wl-copy", fmt.Sprintf("%q", strings.TrimSpace(text)))
			fmt.Println(copyCmd)
			err := exec.Command("/bin/sh", "-c", copyCmd).Run()
			if err != nil {
				return err
			}
		}

	default:
		// TODO

	}
	return nil
}

func MonitorRecvMessage(useBase64 bool) {
	if conn == nil {
		return
	}
	for {
		message := make([]byte, 1024)
		cnt, err := conn.Read(message)
		if err != nil {
			return
		}
		message = message[:cnt]
		var msg models.TCPMessage
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Default().Println("Unmarshal message error: ", err)
			continue
		}
		err = WriteToClipBoard(msg.Data, useBase64)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func MonitorClipBoard() {
	for {
		content, err := ReadFromClipBoard()
		if err != nil {
			log.Default().Println("Get clipboard content error: ", err)
			time.Sleep(1 * time.Second)
			continue
		}
		content = strings.TrimSpace(content)
		contentHash := utils.GetHash(content)
		if contentHash == currentClipBoardHash {
			continue
		}

		currentClipBoardHash = contentHash
		fmt.Println("ContentHash: ", contentHash)
		fmt.Println("Content: ", content)

		messageEntity := models.TCPMessage{
			DeviceId:   "test_PC",
			DeviceType: "PC",
			Timestamp:  time.Now().Unix(),
			Data:       content,
		}
		//contentBase64 := utils.GetBase64([]byte(strings.TrimSpace(content)))
		//contentLen := len(content)
		message, err := json.Marshal(messageEntity)
		if err != nil {
			log.Default().Println("Marshal message error: ", err)
		}
		conn.Write(message)
		time.Sleep(1 * time.Second)
	}
}

func handleConnection() {
	if conn == nil {
		return
	}
	defer conn.Close()

	go MonitorRecvMessage(false)
	go MonitorClipBoard()

	select {}
}

func main() {
	host := flag.String("host", "localhost", "host")
	port := flag.Int("port", 8081, "port")
	flag.Parse()
	// print help
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println(Help)
		return
	}
	for {
		err := Connect(*host, *port)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	log.Default().Println("Connected to server")
	GetClipBoardEnv()
	handleConnection()
}
