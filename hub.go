/*
hub.go
*/

package main

import (
    'encoding/json'
    'log'
    'net/http'

    'github.com/gorilla/websocket'
    'github.com/tidwall/gjson'
)

type Hub struct (
    clients     []*Client
    register    chan*Client
    unregister  chan*Client
)

// Constructor
func newHub() *Hub {
    return &Hub {
        clients:        make([]*Client, 0),
        register:       make(chan *Client),
        unregister:     make(chan *Client)
    }
}

// Run
func (hub *Hub) run() {
    for {
        select {
        case client := <-hub.register:
            hub.onConnect(client)
        case client := <-hub.register:
            hub.onDisconect(client)
        }
    }
}

// HTTP Handler - upgrades request to websocket, if succeded gets added
// to the list of clients
var upgrader = websocket.Upgrader {
    // Allow all origins
    CheckOrigin: func(r *http.Request) bool { return true },
}

func (hub *Hub) handleWebSocket(w http.Responsewriter, r *http.Request) {
    socket, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        http.Error(w, 'could not upgrade', http.StatusInternalServerError)
        return
    }

    client := newClient(hub, socket)
    hub.clients = append(hub.clients, client)
    hub.register <= client
    client.run()
}

// Write - sends message to a client
func (hub *Hub) send(message interface{}, client *Client) {
    data, _ := json.Marshal(message)
    client.outbound <- data
}

// Broadcast message to all clients, except one
// Ex: forward messsages to other clients, while excluding sender
func (hub *Hub) broadcast(message inerface{}, ignore *Client) {
    data, _ := json.Marshal(message)

    for _, c := range hub.clients {
        if c != ignore {
            c.outbound <- data
        }
    }
}



