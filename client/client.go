package client

import (
	"os"
	"log"
	"net/http"
	"encoding/base64"
)

// Constants
var targetWebSocketServer = "ws://localhost:8000"

func generateSecWebSocketKey(username string) string {
	return base64.StdEncoding.EncodeToString([]byte(username))
}


func connectToWebSocketServer(username string) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", targetWebSocketServer, nil)
	if err != nil {
		log.Println("Failed to wrap NewRequestWithContext")
		os.Exit(1)
	}
	// Setting headers
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "websocket")
	request.Header.Set("Sec-WebSocket-Key", generateSecWebSocketKey(username))
	request.Header.Set("Sec-WebSocket-Version", "13")

	response, err := client.Do(request)
	if err != nil {
		log.Println("Failed to connect to websocket server")
		os.Exit(1)
	}
	log.Println(response.StatusCode)
}


func main() {
	connectToWebSocketServer("Kai Cenat")
}