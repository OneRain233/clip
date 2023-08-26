package main

import (
	"clipboard/forms"
	"clipboard/models"
	"clipboard/utils"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var conn net.Conn
var clipBoardType string
var currentClipBoardHash string
var deviceId string
var deviceType string
var apiHost string

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
	log.Default().Println("ClipBoardType: ", clipBoardType)
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
			DeviceId:   deviceId,
			DeviceType: deviceType,
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

func sendClipBoardOnce() {
	if conn == nil {
		return
	}
	content, err := ReadFromClipBoard()
	if err != nil {
		log.Default().Println("Get clipboard content error: ", err)
		return
	}
	content = strings.TrimSpace(content)
	log.Default().Println("Content: ", content)

	messageEntity := models.TCPMessage{
		DeviceId:   deviceId,
		DeviceType: deviceType,
		Timestamp:  time.Now().Unix(),
		Data:       content,
	}
	message, err := json.Marshal(messageEntity)
	if err != nil {
		log.Default().Println("Marshal message error: ", err)
	}
	_, err = conn.Write(message)
	if err != nil {
		log.Default().Println("Write message error: ", err)
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

func getLatestClipBoardViaHttp() {
	latestClipBoardApi := fmt.Sprintf("%s/clipboard/latest", apiHost)
	resp, err := http.Get(latestClipBoardApi)
	if err != nil {
		log.Default().Println("Send http request error: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Default().Println("Send http request error: ", resp.Status)
	}

	var respBody forms.GetLatestClipBoardResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)

	if err != nil {
		log.Default().Println("Decode response error: ", err)
	}

	if respBody.Code != 0 {
		log.Default().Println("Get latest clipboard content error: ", respBody.Code)
		return
	}

	content := respBody.Data.Content
	err = WriteToClipBoard(content, false)
	if err != nil {
		log.Default().Println("Write to clipboard error: ", err)
	}

}

func sendClipBoardViaHttp() {
	addClipBoardApi := fmt.Sprintf("%s/clipboard/add", apiHost)
	content, err := ReadFromClipBoard()
	if err != nil {
		log.Default().Println("Get clipboard content error: ", err)
		return
	}
	content = strings.TrimSpace(content)
	log.Default().Println("Content: ", content)

	// send http request
	resp, err := http.Post(
		addClipBoardApi,
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("content=%s&device_id=%s&device_type=%s", content, deviceId, deviceType)),
	)
	if err != nil {
		log.Default().Println("Send http request error: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Default().Println("Send http request error: ", resp.Status)
	}

	log.Default().Println("Send clipboard content to server")
}

func main() {
	host := flag.String("host", "localhost", "host")
	port := flag.Int("port", 8081, "port")
	mode := flag.String("mode", "listen", "mode")
	apiAddr := flag.String("api_host", "http://localhost:18081", "api host")
	flag.Parse()

	// print help
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println(Help)
		return
	}

	apiHost = *apiAddr

	GetClipBoardEnv()
	if *mode == "http_send" {
		sendClipBoardViaHttp()
		return
	} else if *mode == "http_get" {
		getLatestClipBoardViaHttp()
		return
	}
	for {
		err := Connect(*host, *port)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	defer func() {
		if conn != nil {
			log.Default().Println("Close connection")
			conn.Close()
		}
	}()
	log.Default().Println("Connected to server")
	if *mode == "send" {
		sendClipBoardOnce()
		log.Default().Println("Send clipboard content to server")
		return
	}

	handleConnection()
}
