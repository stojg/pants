package structs

import (
	. "github.com/stojg/pants/vector"
)

func NewData() *Data {
	return &Data{
		Position:           &Vec2{},
		Velocity:           &Vec2{},
		MaxVelocity:        0.0,
		Acceleration:       &Vec2{},
		MaxAcceleration:    0.0,
		Orientation:        0.0,
		AngularVelocity:    0.0,
		MaxAngularVelocity: 0.0,
		InvMass:            0.0,
		Rotations:          0.0,
		Height:             0.0,
		Width:              0.0,
		Forces:             &Vec2{},
		Damping:            0.99,
	}
}

type Data struct {
	Position           *Vec2
	Velocity           *Vec2
	MaxVelocity        float64
	Acceleration       *Vec2   // current acceleration
	MaxAcceleration    float64 // max acceleration
	Orientation        float64 // orientation in radians
	AngularVelocity    float64 // rotational velocity in radians/sec
	MaxAngularVelocity float64 // max angular velocity
	InvMass            float64
	Rotations          float64 // rotational acceleration
	Height             float64 // bounding box height
	Width              float64 // bounding box width
	Forces             *Vec2   // accumulated forces acting on this component

	Damping float64 // general damping
}

func (d *Data) Mass() float64 {
	return 1 / d.InvMass
}

// Copy copies the data from one Data struct to another Data struct
func (d *Data) Copy(o *Data) {
	d.Position.Copy(o.Position)
	d.Velocity.Copy(o.Velocity)
	d.MaxVelocity = o.MaxVelocity
	d.Acceleration.Copy(o.Acceleration)
	d.MaxAcceleration = o.MaxAcceleration
	d.Orientation = o.Orientation
	d.AngularVelocity = o.AngularVelocity
	d.MaxAngularVelocity = o.MaxAngularVelocity
	d.InvMass = o.InvMass
	d.Rotations = o.Rotations
	d.Height = o.Height
	d.Width = o.Width
	d.Forces.Copy(o.Forces)
	d.Damping = o.Damping
}
