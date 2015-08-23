package physics

import (
	. "github.com/stojg/pants/vector"
	"math"
)

type DebugWriter interface {
	Line(*Vec2, *Vec2)
}

type Physics struct {
	Position           *Vec2    // current position
	prevPosition       *Vec2    // previous position
	Velocity           *Vec2    // current velocity
	maxVelocity        float64  // max possible velocity
	Acceleration       *Vec2    // current acceleration
	maxAcceleration    float64  // max acceleration
	Orientation        float64  // current orientation in radians
	prevOrientation    float64  // previous orientation in radians
	AngularVelocity    float64  // rotational velocity in radians/sec
	maxAngularVelocity float64  // max angular velocity
	rotations          float64  // rotational acceleration
	forces             *Vec2    // accumulated forces acting on this component
	invMass            float64  // the inverse mass of this object
	damping            float64  // general damping
	collisionGeometry  Geometry // the geometry used for collision detection
}

func NewPhysics(x, y, orientation, mass float64) *Physics {

	p := &Physics{
		Position:        &Vec2{X: x, Y: y},
		prevPosition:    &Vec2{X: x, Y: y},
		Velocity:        &Vec2{},
		maxVelocity:     300,
		Acceleration:    &Vec2{},
		maxAcceleration: 100,
		Orientation:     orientation,
		prevOrientation: orientation,
		forces:          &Vec2{},
		damping:         0.999,
	}

	if mass > 0 {
		p.invMass = 1 / mass
	}

	p.collisionGeometry = &Circle{
		position: p.Position,
		radius:   10,
	}

	return p
}

func (c *Physics) Clone() *Physics {
	t := c.Position.Clone()
	pc := NewPhysics(t.X, t.Y, 0, 1/c.invMass)
	pc.Velocity = c.Velocity.Clone()
	return pc
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

	// check if this component has changed position or rotations since last
	// update
	updated := false
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
	resultingAcc := c.Acceleration.Clone()
	resultingAcc.Add(c.forces.Multiply(c.invMass))

	// update linear velocity from the acceleration
	c.Velocity.Add(resultingAcc.Multiply(duration))

	// impose drag
	c.Velocity = c.Velocity.Multiply(math.Pow(c.damping, duration))

	// look in the direction where you want to go
	c.Orientation = c.Velocity.ToOrientation()

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
	c.Acceleration = a
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
