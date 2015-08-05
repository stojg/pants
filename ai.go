package main

import (
	"math"
)

type AI interface {
	Update(s *Sprite, w *World, t float64)
}

type AIDrunkard struct {
	state string
	steering Steering
}

func (d *AIDrunkard) Update(s *Sprite, w *World, t float64) {

	w.entities.physics[s.Id].forces = &Vec2{0,10}

	for _, input := range s.inputs {
		switch input.Action {
		case "mouseOver":
			d.Flee(w, s)
		case "mouseOut":
			d.Idle(w, s)
		case "clickUp":
			s.Kill()
		}
	}
	if len(s.inputs) == 0 {
		d.Stagger(w, s, t)
	}
	s.inputs = s.inputs[0:0]
}

func (d *AIDrunkard) Flee(w *World, s *Sprite) {
	if d.state == "flee" {
		return
	}
	o := w.rand.Float64() * math.Pi * 2
	w.entities.physics[s.Id].Orientation = o
	w.entities.physics[s.Id].forces = RadiansVec2(o).Multiply(10)
	d.state = "flee"
}

func (d *AIDrunkard) Idle(w *World, s *Sprite) {
	if d.state == "idle" {
		return
	}
	d.state = "idle"
	w.entities.physics[s.Id].forces = (&Vec2{0, 0})
}

func (d *AIDrunkard) Stagger(w *World, s *Sprite, duration float64) {

	if d.state != "stagger" {
		d.state = "stagger"
		d.steering = &Arrive{
			source: w.entities.physics[s.Id],
			target: &PhysicsComponent{
				Position: &Vec2{X:w.rand.Float64()*800, Y:w.rand.Float64()*600},
			},
			targetRadius: 1,
			slowRadius: 300,
		}
	}
	output := d.steering.Get(duration)
	w.entities.physics[s.Id].AddForce(output.linear.Scale(w.entities.physics[s.Id].Mass()))
	w.entities.physics[s.Id].rotations = output.angular

}
