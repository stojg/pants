package collision

import (
	"github.com/stojg/pants/vector"
)


type Contact struct {
	a Geometry
	b Geometry
	normal *vector.Vec2
	penetration float64
}
