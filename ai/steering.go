package ai

import (
	. "github.com/stojg/pants/physics"
	. "github.com/stojg/pants/vector"
)

// SteeringOutput is what every steering behaviour will return. It contains the
// linear and angular acceleration that the steering behaviour wish to set on
// the target entity
type SteeringOutput struct {
	Linear  *Vec2
	Angular float64
}

func NewSteeringOutput() *SteeringOutput {
	return &SteeringOutput{
		Linear:  &Vec2{},
		Angular: 0,
	}
}

// Steering is the interface that all steering behaviour needs to follow
type Steering interface {
	Get() *SteeringOutput
	Target() *Physics
	Entity() *Physics
}

// Seek behavior moves the owner towards the target position. Given a target,
// this behavior calculates the linear steering acceleration which will direct
// the agent towards the target as fast as possible.
type Seek struct {
	entity *Physics
	target *Physics
}

// Entity returns the entity that is the steering behaviour is acting on
func (s *Seek) Entity() *Physics {
	return s.entity
}

// Target returns the target for the steering behaviour
func (s *Seek) Target() *Physics {
	return s.target
}

// Get returns the steering output for Seek
func (s *Seek) Get() *SteeringOutput {
	steering := NewSteeringOutput()
	// Get the direction to the target
	steering.Linear = s.target.Data.Position.Clone().Sub(s.entity.Data.Position)
	// Go full speed ahead
	steering.Linear.Normalize().Scale(s.entity.Data.MaxAcceleration)
	return steering
}

// Flee behavior does the opposite of Seek. It produces a linear steering force
// that moves the agent away from a target position.
type Flee struct {
	*Seek
}

// NewFlee returns a new Flee steering behaviour
func NewFlee(entity, target *Physics) *Flee {
	return &Flee{
		Seek: &Seek{
			entity: entity,
			target: target,
		},
	}
}

// Get returns the steering output for Flee
func (s *Flee) Get() *SteeringOutput {
	steering := NewSteeringOutput()
	// Get the direction to the target
	steering.Linear = s.entity.Data.Position.Clone().Sub(s.target.Data.Position)
	// Go full speed ahead
	steering.Linear.Normalize().Scale(s.entity.Data.MaxAcceleration)
	return steering
}

// Arrive behavior moves the agent towards a target position. It is similar to
// seek but it attempts to arrive at the target position with a zero velocity.
//
// Arrive behavior uses two radii. The targetRadius lets the owner get near
// enough to the target without letting small errors keep it in motion. The
// slowRadius, usually much larger than the previous one, specifies when the
// incoming entity will begin to slow down. The algorithm calculates an ideal
// speed for the owner. At the slowing-down radius, this is equal to its maximum
// linear speed. At the target point, it is zero (we want to have zero speed
// when we arrive). In between, the desired speed is an interpolated
// intermediate value, controlled by the distance from the target.
//
// The direction toward the target is calculated and combined with the desired
// speed to give a target velocity. The algorithm looks at the current velocity
// of the character and works out the acceleration needed to turn it into the
// target velocity. We can't immediately change velocity, however, so the
// acceleration is calculated based on reaching the target velocity in a fixed
// time scale known as timeToTarget. This is usually a small value; it defaults
// to 0.1 seconds which is a good starting point.
type Arrive struct {
	*Seek
	targetRadius float64 // Holds the radius that says we are at the target
	slowRadius   float64 // Start slowing down at this radius
	timeToTarget float64 // How fast we are trying to get to the target, 0.1
}

// NewArrive returns a new Arrive steering behaviour
func NewArrive(entity, target *Physics, targetRadius, slowRadius float64) *Arrive {
	return &Arrive{
		Seek: &Seek{
			entity: entity,
			target: target,
		},
		targetRadius: targetRadius,
		slowRadius:   slowRadius,
		timeToTarget: 0.1,
	}
}

// Get returns the steering output for Arrive
func (s *Arrive) Get() *SteeringOutput {
	steering := NewSteeringOutput()

	direction := s.target.Data.Position.Clone().Sub(s.entity.Data.Position)
	distance := direction.Length()

	// We have arrived, no output
	if distance <= s.targetRadius {
		return steering
	}

	targetSpeed := s.entity.Data.MaxVelocity
	// We are inside the slow radius, so don't go full speed
	if distance < s.slowRadius {
		targetSpeed *= distance / s.slowRadius
	}

	// The target velocity combines speed and direction
	targetVelocity := direction.Clone()
	targetVelocity.Normalize().Scale(targetSpeed)

	// Acceleration tries to get to the target velocity
	steering.Linear = targetVelocity.Sub(s.entity.Data.Velocity)
	// try to get there in timeToTarget seconds
	steering.Linear.Scale(1 / s.timeToTarget)

	if steering.Linear.Length() > s.entity.Data.MaxAcceleration {
		steering.Linear = steering.Linear.Normalize().Scale(s.entity.Data.MaxAcceleration)
	}
	return steering
}
