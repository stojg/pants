package physics

import (
	. "github.com/stojg/pants/vector"
	"math"
)

type DebugWriter interface {
	Line(*Vec2, *Vec2)
}

type Physics struct {
	Position        *Vec2
	prevPosition    *Vec2
	Velocity        *Vec2
	maxVelocity     float64
	forces          *Vec2
	maxAcceleration float64

	acceleration *Vec2

	Orientation        float64
	prevOrientation    float64
	AngularVelocity    float64
	rotations          float64
	maxAngularVelocity float64

	forceAccum *Vec2
	invMass    float64

	damping float64
}

func NewPhysics(x, y, orientation, mass float64) *Physics {

	var invMass float64
	if mass > 0 {
		invMass = 1 / mass
	}

	return &Physics{
		Position:        &Vec2{X: x, Y: y},
		prevPosition:    &Vec2{X: x, Y: y},
		Velocity:        &Vec2{0, 0},
		maxVelocity:     300,
		maxAcceleration: 100,

		acceleration: &Vec2{0, 0},

		forces: &Vec2{0, 0},

		Orientation:     orientation,
		prevOrientation: orientation,
		AngularVelocity: 0,
		rotations:       0,

		invMass: invMass,
		damping: 0.999,
	}
}

func (c *Physics) Clone() *Physics {
	t := c.Position.Clone()
	pc := NewPhysics(t.X, t.Y, 0, 1/c.invMass)
	pc.Velocity = c.Velocity.Clone()
	return pc
}

func (c *Physics) getPosition(a *Vec2) {
	a.Copy(c.Position)
}

func (c *Physics) Update(id uint64, w DebugWriter, duration float64) bool {

	//	if c.Velocity.Length() > 1 {
	//		w.Line(c.Position.Clone(), c.Position.Clone().Add(c.Velocity))
	//	}
	//
	//	if c.forces.Length() > 1 {
	//		w.Line(c.Position.Clone(), c.Position.Clone().Add(c.forces.Multiply(duration*10).Multiply(-1)))
	//	}

	c.integrate(id, duration)

	updated := false
	// mark as updated
	if !c.Position.Equals(c.prevPosition) {
		c.prevPosition.Copy(c.Position)
		updated = true
	}
	if math.Abs(c.prevOrientation-c.Orientation) < EPSILON {
		c.prevOrientation = c.Orientation
		updated = true
	}
	return updated
}

func (c *Physics) integrate(id uint64, duration float64) {
	if c.invMass <= 0 {
		return
	}

	if duration <= 0 {
		panic("Duration cant be zero")
	}

	// update linear position
	c.Position.Add(c.Velocity.Multiply(duration))

	// Work out the acceleration from the force
	resultingAcc := c.acceleration.Clone()
	resultingAcc.Add(c.forces.Multiply(c.invMass))

	// update linear velocity from the acceleration
	c.Velocity.Add(resultingAcc.Multiply(duration))

	// impose drag
	c.Velocity = c.Velocity.Multiply(math.Pow(c.damping, duration))

	// look in the direction where you want to go
	c.Orientation = c.forces.ToOrientation()

	// clear forces
	c.clearForces()
}

func (c *Physics) Mass() float64 {
	return 1 / c.invMass
}

func (c *Physics) SetMass(m float64) {
	c.invMass = 1 / m
}

func (c *Physics) SetDamping(d float64) {
	c.damping = d
}

func (c *Physics) SetAcceleration(a *Vec2) {
	c.acceleration = a
}

func (p *Physics) MaxAcceleration() float64 {
	return p.maxAcceleration
}

func (p *Physics) MaxVelocity() float64 {
	return p.maxVelocity
}

func (c *Physics) SetRotations(a float64) {
	c.rotations = a
}

func (c *Physics) AddForce(v *Vec2) {
	c.forces.Add(v)
}

func (c *Physics) clearForces() {
	c.forces.Clear()
}
