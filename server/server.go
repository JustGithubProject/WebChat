package main

import (
	// "fmt"
	"log"
	"os"
	"io"
	"net/http"
	"crypto/sha1"
	"encoding/base64"
)

var GUID = "b218e612-c447-4bb0-8513-3c8319601232"

func generateSecWebSocketAccept(secWebSocketKey string) string {
	hasher := sha1.New()
	hasher.Write([]byte(secWebSocketKey + GUID))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func handleClient(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Connection") != "Upgrade" || r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "400 - Invalid Request", http.StatusBadRequest)
		os.Exit(1) 
	}

	// Generating value for Sec-WebSocket-Accept header
	secWebSocketKey := r.Header.Get("Sec-WebSocket-Key")
	if secWebSocketKey == "" {
		http.Error(w, "400 - Missing Sec-WebSocket-Key", http.StatusBadRequest)
		os.Exit(1)
	}
	generatedSecWebSocketAcceptString := generateSecWebSocketAccept(secWebSocketKey)

	// Setting needed headers for websocket server
	w.Header().Set("Connection", "Upgrade")
	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Sec-WebSocket-Accept", generatedSecWebSocketAcceptString)

	w.WriteHeader(http.StatusSwitchingProtocols)
	

	// Handshake complete. Now upgrade the connection
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Println("Failed to hijack connection: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	log.Println("WebSocket connection established")
	io.WriteString(conn, "Hello, WebSocket\n")
}

func main() {
	http.HandleFunc("/ws", handleClient)
	log.Fatal(http.ListenAndServe(":8000", nil))
}