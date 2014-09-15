package main

import (
	"github.com/conformal/websocket"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// Websocket connection settings
var upgrader = websocket.Upgrader{
	ReadBufferSize:  64,
	WriteBufferSize: 1024,
	CheckOrigin: func(req *http.Request) bool {
		return true
	},
}

// Write data to websocket connection
func (c *connection) write(messageType int, data []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(messageType, data)
}

// Wait for data on send-channel and write data to client
func (c *connection) writePump() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	for {
		if _, _, err := c.ws.ReadMessage(); err != nil {
			break
		}
	}
}

// Handle incoming websocket connections
func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	conn := &connection{
		send: make(chan []byte, 256),
		ws:   ws,
	}

	// Save connection and start message pumps
	h.register <- conn
	go conn.writePump()
	conn.readPump()
}
