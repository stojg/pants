package main

import (
	"github.com/mitchellh/mapstructure"
	"github.com/stojg/pants/grid"
	"github.com/stojg/pants/network"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type Response struct {
	Tick      uint64
	Topic     string
	Timestamp float64
	Data      []*EntityUpdate
}

type MapResponse struct {
	Tick      uint64
	Topic     string
	Timestamp float64
	Data      []*grid.Node
}

type TimeRequest struct {
	Topic  string
	Server float64
	Client float64
}

type InputRequest struct {
	Topic string
	Type  string
	Data  []int
}

type NetworkManager struct {
	server *network.Server
	world  *World
}

func (n *NetworkManager) SendState(entities map[uint64]bool, w *World, current time.Time) {
	changedSprites := w.entityManager.Changed()
	ts := float64(current.UnixNano()) / 1000000
	if len(changedSprites) > 0 {
		msg := &Response{
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

func (n *NetworkManager) SendMap(connId uint64, wMap []*grid.Node) {
	ts := float64(time.Now().UnixNano()) / 1000000

	msg := &MapResponse{
		Topic:     "map",
		Data:      wMap,
		Timestamp: ts,
	}
	data, err := bson.Marshal(msg)
	if err != nil {
		log.Printf("error %s", err)
		return
	}
	n.server.Send(connId, data)
}

func (n *NetworkManager) BroadcastMap(wMap []*grid.Node) {
	ts := float64(time.Now().UnixNano()) / 1000000

	msg := &MapResponse{
		Topic:     "map",
		Data:      wMap,
		Timestamp: ts,
	}
	data, err := bson.Marshal(msg)
	if err != nil {
		log.Printf("error %s", err)
		return
	}
	n.server.Broadcast(data)
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

func (n *NetworkManager) handleInput(in *network.Request) {
	msg := map[string]interface{}{}
	err := bson.Unmarshal(in.Message, msg)
	if err != nil {
		log.Printf("error on bson.Umarshal: %s", err)
		return
	}
	in.DecodedMessage = msg

	switch msg["topic"] {
	case "time_request":
		n.handleTimeCheck(in)
		n.SendMap(in.ConnectionID, n.world.worldMap.Compress())
	case "input":
		n.handleInputRequest(in)
	default:
		log.Printf("unhandled message topic '%s'", msg["topic"])
	}
}

func (n *NetworkManager) handleTimeCheck(req *network.Request) {
	log.Printf("time_request received")
	var request TimeRequest
	if err := mapstructure.Decode(req.DecodedMessage, &request); err != nil {
		log.Printf("error: could not decode incoming message: %s", err)
	}
	response := &TimeRequest{
		Topic:  "time_request",
		Server: float64(time.Now().UnixNano()) / 1000000,
		Client: request.Client,
	}
	bson, _ := bson.Marshal(response)
	n.server.Send(req.ConnectionID, bson)
	log.Printf("time_request sent")
}

type ClickEvent struct {
}

func (evt *ClickEvent) Code() EventCode {
	return EVT_CLICK
}

func (n *NetworkManager) handleInputRequest(msg *network.Request) {
	var input InputRequest
	if err := mapstructure.Decode(msg.DecodedMessage, &input); err != nil {
		log.Printf("error: could not decode incoming message: %s", err)
	}
	events.ScheduleEvent(&ClickEvent{})
	log.Printf("recieved input: %#V", input)
}
