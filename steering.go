package main

import (
	. "github.com/stojg/pants/vector"
)

func NewSteeringOutput() *SteeringOutput {
	return &SteeringOutput{
		linear:  &Vec2{},
		angular: 0,
	}
}

type SteeringOutput struct {
	linear  *Vec2
	angular float64
}

type Steering interface {
	Get(dt float64) *SteeringOutput
}

type Seek struct {
	source *PhysicsComponent
	target *PhysicsComponent
}

func (s Seek) Get(dt float64) *SteeringOutput {
	steering := NewSteeringOutput()
	// Get the direction to the target
	steering.linear = s.target.Position.Clone().Sub(s.source.Position)
	// Go full speed ahead
	steering.linear.Normalize().Scale(s.source.maxAcceleration)
	return steering
}

type Arrive struct {
	source *PhysicsComponent
	target *PhysicsComponent
	// Holds the radius that says we are at the target
	targetRadius float64
	// Start slowing down at this radius
	slowRadius float64
	// How fast we are trying to get to the target, 0.1
	timeToTarget float64
}

func (s Arrive) Get(dt float64) *SteeringOutput {
	steering := NewSteeringOutput()

	direction := s.target.Position.Clone().Sub(s.source.Position)
	distance := direction.Length()

	// We have arrived, no output
	if distance < s.targetRadius {
		return steering
	}

	targetSpeed := s.source.maxVelocity
	// We are inside the slow radius, so don't go full speed
	if distance < s.slowRadius {
		targetSpeed *= distance / s.slowRadius
	}

	// The target velocity combines speed and direction
	targetVelocity := direction.Clone()
	targetVelocity.Normalize().Scale(targetSpeed)

	// Acceleration tries to get to the target velocity
	steering.linear = targetVelocity.Sub(s.source.Velocity)
	// try to get there in 0.1 seconds
	steering.linear.Scale(1 / 0.1)

	if steering.linear.Length() > s.source.maxAcceleration {
		steering.linear = steering.linear.Normalize().Scale(s.source.maxAcceleration)
	}
	return steering
}
