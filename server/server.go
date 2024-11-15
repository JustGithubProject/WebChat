package server

import (
	// "fmt"
	"log"
	"net/http"
	"crypto/sha1"
	"encoding/hex"
	"encoding/base64"
)

var GUID = "b218e612-c447-4bb0-8513-3c8319601232"

func generateSecWebSocketAccept(secWebSocketKey string) string {
	tempString := GUID + secWebSocketKey
	hasher := sha1.New()
	hasher.Write([]byte(tempString))
	sha1Hash := hex.EncodeToString(hasher.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha1Hash))
}

func handleClient(w http.ResponseWriter, r *http.Request) {

	// Generating value for Sec-WebSocket-Accept header
	secWebSocketKey := r.Header.Get("Sec-WebSocket-Key")
	generatedSecWebSocketAcceptString := generateSecWebSocketAccept(secWebSocketKey)

	// Setting needed headers for websocket server
	w.Header().Set("Connection", "Upgrade")
	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Sec-WebSocket-Accept", generatedSecWebSocketAcceptString)

	w.WriteHeader(http.StatusSwitchingProtocols)
}

func main() {
	http.HandleFunc("/ws", handleClient)
	log.Fatal(http.ListenAndServe(":8000", nil))
}