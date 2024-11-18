package main

import (
	"bufio"
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


func encodeWebSocketMessage(message string) []byte {
	length := len(message)
	/*
		0x81 = 10000001
		It will set first bit to 1 (FIN bit)
	*/
	header := []byte{0x81}

	if length < 126 {
		/*
			EXAMPLE:
				if we have length = 0x33 (51)
				header = {0x81, 0x33}
		*/
		header = append(header, byte(length))

	} else if length <= 65535 {
		/*
			EXAMPLE:
				if we have length = 0xABCD
				0xABCD >> 8 = 0xAB
				0xABCD & 0xFF = 0xCD
				header = {0x81, 0xAB, 0xCD}
		*/
		header = append(header, 126, byte(length >> 8), byte(length & 0xFF))
	} else {
		log.Println("Message too long")
		return nil
	}
	return append(header, []byte(message)...)

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

	// Enter a loop to read response from websocket server
	go func() {
		buffer := make([]byte, 1024)

		for {
			n, err := conn.Read(buffer)
			if err != nil {
				log.Println("Failed to read response: ", err)
				os.Exit(1)
			}
			log.Println(string(buffer[:n]))
		}
	}()

	// Enter a loop to do request to websocket server
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "exit" {
			log.Println("Closing connection")
			os.Exit(1)
		}
		encodedMessage := encodeWebSocketMessage(message)
		conn.Write(encodedMessage)
	}
}

func main() {
	connectToWebSocketServer("Kai Cenat")
}