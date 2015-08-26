package main

import (
	"github.com/stojg/pants/grid"
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
	inputs []*InputRequest
	Type   string
	Tile   *grid.Node
}

func (s *Entity) AddInput(i *InputRequest) {
	s.inputs = append(s.inputs, i)
}
