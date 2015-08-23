package main

import (
	. "github.com/stojg/pants/vector"
)

func init() {
	entityManager = NewEntityManager()
}

type Line struct {
	Position *Vec2
	End      *Vec2
}

type Entity struct {
	Id     uint64 `bson:",minsize"`
	Image  string `bson:",minsize,omitempty"`
	Dead   bool
	inputs []*InputRequest
	Type   string
}

func (s *Entity) Kill() {
	s.Dead = true
	entityManager.updated[s.Id] = true
}

func (s *Entity) AddInput(i *InputRequest) {
	s.inputs = append(s.inputs, i)
}
