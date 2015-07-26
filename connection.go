package main

import (
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second    // Time allowed to write a message to the peer.
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize = 512                 // Maximum message size allowed from peer.
)

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	ws   *websocket.Conn // The websocket connection.
	send chan []byte     // Buffered channel of outbound messages.
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("error on socket read: %s", err)
			break
		}
		c.handleMessage(message)
	}
}

func (c *connection) handleMessage(in []byte) {
	msg := map[string]interface{}{}
	err := bson.Unmarshal(in, msg)
	if err != nil {
		log.Printf("error on bson umarshal read: %s", err)
		return
	}
	switch msg["topic"] {
	case "time_request":
		c.handleTimeCheck(msg)
	case "input":
		c.handleInputRequest(msg)
	default:
		log.Printf("unhandled message topic '%s'", msg["topic"])
	}
}

func (c *connection) handleInputRequest(msg map[string]interface{}) {
	var request InputRequest
	if err := mapstructure.Decode(msg, &request); err != nil {
		log.Printf("error: could not decode incoming message: %s", err)
	}
	sprite, ok := list.sprites[request.Id]
	if !ok {
		log.Printf("error: no sprite with id %d for input command", request.Id)
		return
	}
	sprite.AddInput(&request)
}

func (c *connection) handleTimeCheck(msg map[string]interface{}) {
	var request TimeRequest
	if err := mapstructure.Decode(msg, &request); err != nil {
		log.Printf("error: could not decode incoming message: %s", err)
	}
	response := &TimeRequest{
		Topic:  "time_request",
		Server: float64(time.Now().UnixNano()) / 1000000,
		Client: request.Client,
	}
	bson, _ := bson.Marshal(response)
	c.send <- bson
	log.Printf("time_request sent: %f", response.Server)

}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				if err := c.write(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("error %s", err)
				}
				return
			}
			if err := c.write(websocket.BinaryMessage, message); err != nil {
				log.Printf("error %s", err)
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("error %s", err)
			}
		}
	}
}
