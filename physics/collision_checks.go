package physics

import "math"

// CircleVsCircle does a collision check between two cicles
func CircleVsCircle(g1, g2 *Circle) *Contact {
	sqrDistance := g1.position.Clone().Sub(g2.position).SquareLength()
	radiusSum := g1.radius + g2.radius
	if sqrDistance >= radiusSum*radiusSum {
		return nil
	}
	cc := &Contact{
		normal:      g1.position.Clone().Sub(g2.position).Normalize(),
		penetration: radiusSum - math.Sqrt(sqrDistance),
	}
	return cc
}
