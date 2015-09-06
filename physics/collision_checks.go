package physics

import "math"

// CircleVsCircle does a collision check between two cicles
func CircleVsCircle(g1, g2 *Circle) *Contact {
	sqrDistance := g1.Position.Clone().Sub(g2.Position).SquareLength()
	radiusSum := g1.Radius + g2.Radius
	if sqrDistance >= radiusSum*radiusSum {
		return nil
	}
	cc := &Contact{
		normal:      g1.Position.Clone().Sub(g2.Position).Normalize(),
		penetration: radiusSum - math.Sqrt(sqrDistance),
	}
	return cc
}
