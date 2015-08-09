package physics

import (
	. "github.com/stojg/pants/vector"
)

type StaticForce struct {
	force *Vec2
}

func (sf *StaticForce) SetForce(a *Vec2) {
	sf.force = a
}
