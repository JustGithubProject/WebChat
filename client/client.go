package main

import (
	"os"
	"log"
	"net"
	"encoding/base64"
)

// Constants
var targetWebSocketServer = "localhost:8000"

func generateSecWebSocketKey(username string) string {
	return base64.StdEncoding.EncodeToString([]byte(username))
}


func connectToWebSocketServer(username string) {
	key := generateSecWebSocketKey(username)

	// Manually create a TCP connection to send HTTP upgrade request
	conn, err := net.Dial("tcp", targetWebSocketServer)
	if err != nil {
		log.Println("Failed to connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Build HTTP Upgrade request
	req := "GET /ws HTTP/1.1\r\n" +
		"Host: " + targetWebSocketServer + "\r\n" +
		"Connection: Upgrade\r\n" +
		"Upgrade: websocket\r\n" +
		"Sec-WebSocket-Key: " + key + "\r\n" +
		"Sec-WebSocket-Version: 13\r\n\r\n"

	// Send request
	conn.Write([]byte(req))

	// Read and print server response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Failed to read response:", err)
		os.Exit(1)
	}
	log.Println("Server response:")
	log.Println(string(buf[:n]))
}



func main() {
	connectToWebSocketServer("Kai Cenat")
}