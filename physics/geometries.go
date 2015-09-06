package physics

import (
	. "github.com/stojg/pants/vector"
)

type Geometry interface{}

func NewCircle(x, y, radius float64) *Circle {
	return &Circle{
		Position: &Vec2{x, y},
		Radius:   radius,
	}
}

type Circle struct {
	Position *Vec2
	Radius   float64
}
