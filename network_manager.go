package main

import (
	"github.com/mitchellh/mapstructure"
	"github.com/stojg/pants/network"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type Message struct {
	Tick      uint64
	Topic     string
	Data      []*EntityUpdate
	Timestamp float64
}

type TimeRequest struct {
	Topic  string
	Server float64
	Client float64
}

type InputRequest struct {
	Topic  string
	Action string
	Id     uint64
}

type NetworkManager struct {
	server *network.Server
}

func (n *NetworkManager) SendState(entities map[uint64]bool, w *World, current time.Time) {
	changedSprites := w.entityManager.Changed()
	ts := float64(current.UnixNano()) / 1000000
	if len(changedSprites) > 0 {
		msg := &Message{
			Topic:     "update",
			Data:      changedSprites,
			Tick:      w.gameTick,
			Timestamp: ts,
		}
		bson, err := bson.Marshal(msg)
		if err != nil {
			log.Printf("error %s", err)
			return
		}
		n.server.Broadcast(bson)
	}
}

func (n *NetworkManager) Inputs() {
	messages := n.server.Incoming()
	if len(messages) == 0 {
		return
	}
	for _, in := range messages {
		n.handleInput(in)
	}
}

func (n *NetworkManager) handleInput(in []byte) {
	msg := map[string]interface{}{}
	err := bson.Unmarshal(in, msg)
	if err != nil {
		log.Printf("error on bson.Umarshal: %s", err)
		return
	}

	switch msg["topic"] {
	case "time_request":
		n.handleTimeCheck(msg)
	case "input":
		n.handleInputRequest(msg)
	default:
		log.Printf("unhandled message topic '%s'", msg["topic"])
	}
}

func (n *NetworkManager) handleTimeCheck(msg map[string]interface{}) {
	log.Printf("time_request received")
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
	n.server.Broadcast(bson)
	log.Printf("time_request sent")
}

func (n *NetworkManager) handleInputRequest(msg map[string]interface{}) {
	var input InputRequest
	if err := mapstructure.Decode(msg, &input); err != nil {
		log.Printf("error: could not decode incoming message: %s", err)
	}

	//	input.Id

}
