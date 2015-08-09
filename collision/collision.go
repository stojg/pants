package collision

func CircleVsCircle(g1, g2 *Circle) *Contact {
	sqrDistance := g1.position.Clone().Sub(g2.position).SquareLength()
	radiusSum := g1.radius + g2.radius
	if sqrDistance > (radiusSum)*(radiusSum) {
		return nil
	}

	// @todo(stig): calculate normal and penetration
	cc := &Contact{
		a: g1,
		b: g2,
	}
	return cc
}
