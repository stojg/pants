package main

import (
	. "github.com/stojg/pants/vector"
	"math"
//	"log"
)

func NewPhysicsComponent(x, y, orientation float64) *PhysicsComponent {
	return &PhysicsComponent{
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

		invMass: 0,
		damping: 0.999,
	}
}

type PhysicsComponent struct {
	Position           *Vec2
	prevPosition       *Vec2
	Velocity           *Vec2
	maxVelocity        float64
	forces             *Vec2
	maxAcceleration    float64

	acceleration       *Vec2

	Orientation        float64
	prevOrientation    float64
	AngularVelocity    float64
	rotations          float64
	maxAngularVelocity float64

	forceAccum         *Vec2
	invMass            float64

	damping            float64
}

func (c *PhysicsComponent) Update(sprite *Sprite, w *World, seconds float64) {

	c.integrate(sprite, seconds)

	// mark as updated
	if !c.Position.Equals(c.prevPosition) {
		c.prevPosition.Copy(c.Position)
		list.updated[sprite.Id] = true
	}
	if math.Abs(c.prevOrientation-c.Orientation) < EPSILON {
		c.prevOrientation = c.Orientation
		list.updated[sprite.Id] = true
	}
}

func (c *PhysicsComponent) integrate(sprite *Sprite, duration float64) {
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

func (c *PhysicsComponent) Mass() float64 {
	return 1/c.invMass
}

func (c *PhysicsComponent) setMass(m float64) {
	c.invMass = 1/m
}

func (c *PhysicsComponent) setDamping(d float64) {
	c.damping = d
}

func (c *PhysicsComponent) setAcceleration(a *Vec2) {
	c.acceleration = a
}

func (c *PhysicsComponent) AddForce(v *Vec2) {
	c.forces.Add(v)
}

func (c *PhysicsComponent) clearForces() {
	c.forces.Clear()
}
