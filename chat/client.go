package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct{
	// Socket is the web socket for this client
	socket *websocket.Conn

	// Send is a buffered channel on which messages are sent
	send chan []byte

	// Room is the room thid client is chatting in
	room *Room

	// Associated client 
	user *User
}

func (c *Client) read(){
	defer c.socket.Close()
	for{
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			log.Println("Error occured in read method of Client struct", err)
			return
		}
		c.room.forward <- msg
	}
}

func (c *Client) write(){
	defer c.socket.Close()
	for msg := range c.send{
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil{
			log.Println("Error occured in write method of Client struct", err)
			return
		}
	}
}