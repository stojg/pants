package collision

func CircleVsCircle(g1, g2 *Circle) *CollisionContact {
	cc := &CollisionContact{}
	sqrDistance := g1.position.Clone().Sub(g2.position).SquareLength()
	radiusSum := g1.radius + g2.radius
	if sqrDistance > (radiusSum)*(radiusSum) {
		return cc
	}
	cc.hit = true
	cc.a = g1
	cc.b = g2
	return cc
}
