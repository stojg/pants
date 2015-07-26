package main

type Message struct {
	Tick      uint64
	Topic     string
	Data      []*Sprite
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
