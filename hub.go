package main

import (
	"log"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

// Wait for new connections to connect/disconnect and keeps internal conenction
// list clean.
func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			log.Println("New websocket connection: ", c.ws.RemoteAddr().String())
			h.connections[c] = true
		case c := <-h.unregister:
			log.Println("Lost websocket connection: ", c.ws.RemoteAddr().String())
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			// Send message to every connection
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}
