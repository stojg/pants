package physics

import (
	"github.com/stojg/pants/structs"
//	. "github.com/stojg/pants/vector"
	"math"
)

type Physics struct {
	Data              *structs.Data
	PrevData          *structs.Data
	CollisionGeometry Geometry // the geometry used for collision detection
}

func NewPhysics(data *structs.Data) *Physics {

	data.MaxVelocity = 300
	data.MaxAcceleration = 100

	p := &Physics{
		Data:              data,
		PrevData:          structs.NewData(),
		CollisionGeometry: &Circle{
			Position: data.Position,
		},
	}

	p.PrevData.Copy(p.Data)
	if data.Width > data.Height {
		p.CollisionGeometry.(*Circle).Radius = data.Width / 2
	} else {
		p.CollisionGeometry.(*Circle).Radius = data.Height / 2
	}
	return p
}

func (p *Physics) SetSize(width, height float64) {

}

func (c *Physics) Integrate(data *structs.Data, duration float64) {

	if data.InvMass <= 0 {
		return
	}

	if duration <= 0 {
		panic("Duration cant be zero")
	}

	// update linear position
	data.Position.Add(data.Velocity.Multiply(duration))

	// Work out the acceleration from the force
	resultingAcc := data.Acceleration.Clone()
	resultingAcc.Add(data.Forces.Multiply(data.InvMass))

	// update linear velocity from the acceleration
	data.Velocity.Add(resultingAcc.Multiply(duration))

	// impose drag
	data.Velocity = data.Velocity.Multiply(math.Pow(data.Damping, duration))

	// look in the direction where you want to go
	data.Orientation = data.Velocity.ToOrientation()
	// data.Orientation += data.Rotations * duration

	// clear forces
	c.Data.Forces.Clear()
}
