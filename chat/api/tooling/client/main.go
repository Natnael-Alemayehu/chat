package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	url = "ws://localhost:3000/v1/connect"
)

func main() {
	if err := hack1(); err != nil {
		log.Fatal(err)
	}
}

func hack1() error {
	header := make(http.Header)

	// --------------------------------------------------------------

	// Client creatin a connection to the server
	socket, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		return fmt.Errorf("dail: %v", err)
	}
	defer socket.Close()

	_, msg, err := socket.ReadMessage()
	if err != nil {
		return fmt.Errorf("read message: %v", err)
	}

	// Checking if the message is "Hello"
	if string(msg) != "Hello" {
		return fmt.Errorf("unexpected message: %s", msg)
	}

	// --------------------------------------------------------------

	// Sending handshake message to the server
	user := struct {
		ID   uuid.UUID
		Name string
	}{
		ID:   uuid.New(),
		Name: "Nate",
	}

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshal: %v", err)
	}

	err = socket.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return fmt.Errorf("write message: %v", err)
	}

	// --------------------------------------------------------------

	// Reading the response from the server Expecting a welcome Nate message
	_, msg, err = socket.ReadMessage()
	if err != nil {
		return fmt.Errorf("read message: %v", err)
	}

	if string(msg) != "Welcome Nate" {
		return fmt.Errorf("unexpected message: %s", msg)
	}

	fmt.Println(string(msg))

	return nil
}
