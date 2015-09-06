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

type EntityType byte

const (
	ENT_ENT1 EntityType = iota
	ENT_FOOD
)

// @todo send to client
type EntityProperty struct {
	Sprite  string  `bson:",minsize,omitempty"`
	Height  float64 `bson:",minsize,omitempty"`
	Width   float64 `bson:",minsize,omitempty"`
	invMass float64
	ai      bool
}

var entityProperties map[EntityType]*EntityProperty

func init() {
	entityProperties = make(map[EntityType]*EntityProperty, 0)
	entityProperties[ENT_ENT1] = &EntityProperty{
		Sprite:  "assets/basics/arrow.png",
		Height:  20,
		Width:   20,
		invMass: 1,
		ai:      true,
	}

	entityProperties[ENT_FOOD] = &EntityProperty{
		Sprite:  "assets/basics/small_circle.png",
		Height:  16,
		Width:   16,
		invMass: 0,
	}

}

type Entity struct {
	Id         uint64 `bson:",minsize"`
	Image      string `bson:",minsize,omitempty"`
	inputs     []*InputRequest
	Type       EntityType
	Properties *EntityProperty
	Tile       *grid.Node
}

func (s *Entity) AddInput(i *InputRequest) {
	s.inputs = append(s.inputs, i)
}
