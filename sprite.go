package main

import ()

func init() {
	list = NewEntityList()
}

type Sprite struct {
	Id          uint64 `bson:",minsize"`
	Image       string  `bson:",minsize,omitempty"`
	Dead        bool
	inputs      []*InputRequest
	changed     bool
}

func (s *Sprite) Kill() {
	s.Dead = true
	s.changed = true
}

func (s *Sprite) AddInput(i *InputRequest) {
	s.inputs = append(s.inputs, i)
}
