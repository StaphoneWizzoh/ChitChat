package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

type Room struct{
	// Foward is a channel that holds incoming messages
	// that should be fowarded to the other clients.
	forward chan []byte

	// Join is a channel for clients wishing to join the room
	join chan *Client

	// Leave is a channel for clients wishing to leave the room
	leave chan *Client

	// Clients holds all current clients in this room
	clients map[*Client]bool
}

// newRoom makes a new room.
func newRoom() *Room{
	return &Room{
		forward: make(chan []byte),
		join: make(chan *Client),
		leave: make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

func (r *Room) run(){
	for{
		select{
		case client := <-r.join:
			// Joining the room
			r.clients[client] = true
		case client := <-r.leave:
			// Leaving the room
			delete(r.clients, client)
			close(client.send)
		case msg := <- r.forward:
			// Foward message to all clients
			for client := range r.clients{
				client.send <- msg
			}
		}
	}
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request){
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP for Room Object:", err)
		return
	}

	client := &Client{
		socket: socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}

	r.join <- client
	defer func ()  {
		r.leave <- client	
	}()

	go client.write()
	client.read()
}