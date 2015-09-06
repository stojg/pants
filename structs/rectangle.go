package structs

import (
	. "github.com/stojg/pants/vector"
)

type Rectangle struct {
	Position   *Vec2
	HalfWidth  float64
	HalfHeight float64
}
