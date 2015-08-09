package tree

import (
	. "github.com/stojg/pants/vector"
)

type Rectangle struct {
	position   *Vec2
	halfHeight float64
	halfWidth  float64
	minX       float64
	maxX       float64
	minY       float64
	maxY       float64
}

func NewRectangle(x, y, halfWidth, halfHeight float64) *Rectangle {
	return &Rectangle{
		position:   &Vec2{x, y},
		halfWidth:  halfWidth,
		halfHeight: halfHeight,
		minX:       x - halfWidth,
		maxX:       x + halfWidth,
		minY:       y - halfHeight,
		maxY:       y + halfHeight,
	}
}

// Intersects returns true if other intersects this
func (r *Rectangle) Intersects(other *Rectangle) bool {
	return r.minX < other.maxX && r.minY < other.maxY && r.maxX > other.minX && r.maxY > other.minY
}

// Contains returns true if other can fit within this
func (r *Rectangle) Contains(other *Rectangle) bool {
	return r.minX <= other.minX && r.minY <= other.minY && r.maxX >= other.maxX && r.maxY >= other.maxY
}
