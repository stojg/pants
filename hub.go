package main

import (
	"labix.org/v2/mgo/bson"
	"log"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	connections map[*connection]bool // Registered connections.
	broadcast   chan []byte          // Inbound messages from the connections.
	register    chan *connection     // Register requests from the connections.
	unregister  chan *connection     // Unregister requests from connections.
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			log.Printf("client connected")
			list.SendAll(c)
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
				log.Printf("client disconnected")
			}
		case m := <-h.broadcast:
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

func (h *hub) Send(m *Message) {
	bson, err := bson.Marshal(m)
	if err != nil {
		log.Printf("error %s", err)
		return
	}
	h.broadcast <- bson
}
