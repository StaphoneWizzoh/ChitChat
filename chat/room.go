package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/StaphoneWizzoh/ChitChat/trace"
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

	// tracer will receive trace information of activity in the room
	tracer trace.Tracer
}

// newRoom makes a new room.
func newRoom() *Room{
	r:=  &Room{
		forward: make(chan []byte),
		join: make(chan *Client),
		leave: make(chan *Client),
		clients: make(map[*Client]bool),
		tracer: trace.Off(),
	}

	// Query messages from the database
	// Load messages into the room
	var messages []Message
	db.Find(&messages)

	// Iterate through messages and froward them to clients
	for _, message := range messages{
		r.forward <- []byte(fmt.Sprintf("[%s] %s", message.Time.Format("2006-01-02 15:04:05"), message.Content))
	}
	return r
}

func (r *Room) run(){
	for{
		select{
		case client := <-r.join:
			// Joining the room
			r.clients[client] = true
			r.tracer.Trace("New client joined.")

		case client := <-r.leave:
			// Leaving the room
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")

		case msg := <- r.forward:
			// Foward message to all clients
			for client := range r.clients{

				if client.user != nil {
					// Save the message to the database using the twmp variable
				senderID := client.user.ID
				SaveMessage(string(msg), senderID)

				// Foward message to the client
				client.send <- msg
				r.tracer.Trace(" -- sent to client")

				} else{
					client.send <- msg
					r.tracer.Trace(" -- sent to client")
				}
				
			}
		}
	}
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request){
	// Extracting the current user from the request context
	// user, ok := req.Context().Value("user").(*User)
	// if !ok{
	// 	http.Error(w, "User not found in context", http.StatusInternalServerError)
	// 	return
	// }


	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP for Room Object:", err)
		return
	}

	client := &Client{
		socket: socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
		// user: user,
	}

	r.join <- client
	defer func ()  {
		r.leave <- client	
	}()

	go client.write()
	client.read()
}