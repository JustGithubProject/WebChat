package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	
	"github.com/gorilla/websocket"
)


/*
	Upgrader will be used to upgrade
	HTTP connections to WebSocket connections

	CheckOrigin allows websocket connections
	from any origin 
*/
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


type Message struct {
	Username string `json:"username"`
	Text string `json:"text"`
}


/*
	clients - It keeps track of connected WebSocket clients
	broadcast - It's used to broadcast messages to all connected clients
*/
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)


func handleHomePage(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("home.html");
	if err != nil {
		log.Println("Failed to read home.html file")
		return
	}
	fmt.Fprint(w, string(content))
}


func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade HTTP server connection to websocket protocol")
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Failed to read json")
			delete(clients, conn)
			return 
		}
		// Writing msg to channel
		broadcast <- msg
	}
}


func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(&msg)
			if err != nil {
				log.Println("Failed to write json")
				client.Close()
				delete(clients, client)
			}
		}
	}
}


func main() {
	http.HandleFunc("/", handleHomePage)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()
	
	http.ListenAndServe(":8000", nil)
}