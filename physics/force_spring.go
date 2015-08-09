package physics

import (
	. "github.com/stojg/pants/vector"
	"math"
)

type SpringForce struct {
	other           *Physics
	springConstance float64
	restLength      float64
}

func (sf *SpringForce) UpdateForce(p *Physics, duration float64) {
	// calculate the vector of the spring
	force := &Vec2{}
	p.getPosition(force)
	force.Sub(sf.other.Position)

	// calculate the magnitude of the force
	magnitude := force.Length()

	magnitude = math.Abs(magnitude - sf.restLength)
	magnitude *= sf.springConstance

	// calculate the final force and apply it
	force.Normalize()
	force.Scale(-magnitude)
	p.AddForce(force)
}