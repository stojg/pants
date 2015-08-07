package main

import (
	. "github.com/stojg/pants/vector"
)

func init() {
	list = NewEntityList()
}

type Line struct {
	Position *Vec2
	End      *Vec2
}

type Sprite struct {
	Id     uint64 `bson:",minsize"`
	Image  string `bson:",minsize,omitempty"`
	Dead   bool
	inputs []*InputRequest
	Type   string
}

func (s *Sprite) Kill() {
	s.Dead = true
	list.updated[s.Id] = true
}

func (s *Sprite) AddInput(i *InputRequest) {
	s.inputs = append(s.inputs, i)
}
