package physics

import (
	"github.com/stojg/pants/structs"
	. "github.com/stojg/pants/vector"
	"math"
)

type Physics struct {
	Data              *structs.Data
	PrevData          *structs.Data
	collisionGeometry Geometry // the geometry used for collision detection
}

func NewPhysics(x, y, orientation, invMass, width, height float64) *Physics {

	data := structs.NewData()
	data.Position = &Vec2{X: x, Y: y}
	data.Orientation = orientation
	data.InvMass = invMass
	data.Width = width
	data.Height = height
	data.MaxVelocity = 300
	data.MaxAcceleration = 100

	p := &Physics{
		Data:              data,
		PrevData:          structs.NewData(),
		collisionGeometry: &Circle{position: data.Position},
	}

	p.PrevData.Copy(p.Data)
	p.SetSize(width, height)

	return p
}

func (p *Physics) SetSize(width, height float64) {
	p.Data.Width = width
	p.Data.Height = height
	if width > height {
		p.collisionGeometry.(*Circle).radius = width / 2
	} else {
		p.collisionGeometry.(*Circle).radius = height / 2
	}
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
