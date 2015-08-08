package main

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	connections map[*connection]bool // Registered connections.
	broadcast   chan []byte          // Inbound messages from the connections.
	register    chan *connection     // Register requests from the connections.
	unregister  chan *connection     // Unregister requests from connections.
	list        *EntityList
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run(list *EntityList) {
	h.list = list
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			log.Printf("client connected")
			h.SendAll(c)
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

func (h *hub) SendAll(c *connection) {
	t := &Message{
		Topic:     "all",
		Data:      list.all(),
		Timestamp: float64(time.Now().UnixNano()) / 1000000,
	}
	msg, _ := bson.Marshal(t)
	c.send <- msg
}

func (h *hub) SendUpdates(list map[uint64]bool, w *World, currentTime time.Time) {
	changedSprites := make([]*EntityUpdate, 0)
	for id := range list {
		changedSprites = append(changedSprites, &EntityUpdate{
			Id:          id,
			X:           w.Physic(id).Position.X,
			Y:           w.Physic(id).Position.Y,
			Orientation: w.Physic(id).Orientation,
			Type:        w.entities.sprites[id].Type,
			Dead:        w.entities.sprites[id].Dead,
			Data: map[string]string{
				"Image": w.entities.sprites[id].Image,
			},
		})
	}

	for _, line := range w.debug {
		// @todo(stig): make sure deleted entities are .. dead
		changedSprites = append(changedSprites, &EntityUpdate{
			X:    line.Position.X,
			Y:    line.Position.Y,
			Type: "graphics",
			Data: map[string]string{
				"toX": fmt.Sprintf("%9.f", line.End.X),
				"toY": fmt.Sprintf("%9.f", line.End.Y),
			},
		})
	}
	w.debug = nil

	if len(changedSprites) > 0 {
		h.Send(&Message{
			Topic:     "update",
			Data:      changedSprites,
			Tick:      w.gameTick,
			Timestamp: float64(currentTime.UnixNano()) / 1000000,
		})
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
