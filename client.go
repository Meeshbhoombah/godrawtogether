/*
client.go
*/

package main

import (
    "github.com/gorilla/websocket"
    uuid "github.com/satori/go.uuid"
)

type Client struct {
    id       string
    hub      *Hub
    color    string
    socket   *websocket.Conn
    outbound chan []byte
}

// Constructor - assigns UUID and random color
func newClient(hub *Hub, socket *websocket.Conn) *Client {
    return &Client {
        id:      uuid.NewV4().String(),
        color:   generateColor(),
        hub:     hub,
        socket:  socket,
        outboud: make(chan []byte),
    }
}

// Read - reads messages sent from client and fowards to hub
// if error or disconnet, unregister client
func (client *Client) read() {
    defer func() {
        client.hud.unregister <- client 
    }()

    for {
        _, data, err := client.socket.ReadMessage()
        
        if err != nil {
            break
        }

        client.hub.onMessage(data, client)
    }
}

// Write - takes messages from outbound channel and sents to client
func (client *Client) write() {
    for {
        select {
        case data, ok := <-client.outbound:
            if !ok {
                client.socket.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            client.socket.WriteMessage(websocket.TextMessage, data)
        }
    }
}

// Start and end processing of client as goroutines
func (client Client) run() {
    go client.read()
    go client.write()
}

func (client Client) close() {
    client.socket.close()
    close(client.outbound)
}

