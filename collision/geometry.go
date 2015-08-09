package collision

import (
	. "github.com/stojg/pants/vector"
)

type CollisionGeometry interface{}

func NewCircle(x, y, radius float64) *Circle {
	return &Circle{
		position: &Vec2{x, y},
		radius:   radius,
	}
}

type Circle struct {
	position *Vec2
	radius   float64
}
