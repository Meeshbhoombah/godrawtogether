/*
client.go
*/

package main

import (
    'github.com/gorilla/websocket'
    uuid 'github.com/satori/go.uuid'
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
        id:      uuid.NewV4().String().
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

