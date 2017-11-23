/*
hub.go
*/

package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
    "github.com/tidwall/gjson"
)

type Hub struct {
    clients     []*Client
    register    chan*Client
    unregister  chan*Client
}

// Constructor
func newHub() *Hub {
    return &Hub {
        clients:        make([]*Client, 0),
        register:       make(chan *Client),
        unregister:     make(chan *Client),
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
        http.Error(w, "could not upgrade", http.StatusInternalServerError)
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
func (hub *Hub) broadcast(message interface{}, ignore *Client) {
    data, _ := json.Marshal(message)

    for _, c := range hub.clients {
        if c != ignore {
            c.outbound <- data
        }
    }
}

// onConnect - called from run, sends user's color and information to
// other clients and notifies them a user has joined
func (hub *Hub) onConnect(client *Client) {
    log.Println("client connected: ", client.socket.RemoteAddr())
    // list of all users
    users := []message.User{}
    for _, c := range hub.clients {
        users = append(users, message.User { ID: c.id, Color: c.color })
    }

    // Notification
    hub.send(message.NewConnected(client.color, users), client)
    hub.broadcast(message.NewUserJoined(client.id, client.color), client)
}

// onDisconnect - removes disconnected client from list of clients
// and notifies all clients that someone left
func (hub *Hub) onDisconnect(client *Client) {
    log.Println("client disconnected: ", client.socket.RemoteAddr())
    client.close

    // find index
    i := -1

    for j, c := range hub.clients {
        if c.id == client.id {
            i = j
            break
        }
    }
 
    // Delete client
    copy(hub.clients[i:], hub.clients[i+1:])
    hub.clients[len(hub.clients)-1] = nil
    hub.clients = hub.clients[:len(hub.clients)-1]

    // Notification
    hub.broadcast(message.NewUserLeft(client.id), nil)
}

// onMessage - called whenever a message is recieved from the client, first
// reads the kind of message, then handles each case
func (hub *Hub) onMessage(data []byte, client *Client) {
    kind := gjson.GetBytes(data, "kind").Int()

    if kind == message.KindStroke {
        var msg message.Stroke

        if json.Unmarshal(data, &msg) != nil {
            return
        }

        msg.UserID = client.id
        hub.broadcast(msg, client)
    } else if kind == message.KindClear {
        var msg message.Clear

        if json.Unmarshal(data, &msg) != nil {
            return
        }
    }

    msg.UserID = client.id
    hub.broadcast(msg, client)
}


