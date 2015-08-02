package main

import (
	"math"
)

type AI interface {
	Update(s *Sprite, w *World, t float64)
}

type AIDrunkard struct {
	state string
}

func (d *AIDrunkard) Update(s *Sprite, w *World, t float64) {
	for _, input := range s.inputs {
		switch input.Action {
		case "mouseOver":
			d.Flee(w, s)
			return
		case "mouseOut":
			d.Idle(w, s)
			return
		case "clickUp":
			s.Kill()
			return
		}
	}
	d.Stagger(w, s)
}

func (d *AIDrunkard) Flee(w *World, s *Sprite) {
	if d.state == "flee" {
		return
	}
	o := w.rand.Float64() * math.Pi * 2
	s.SetOrientation(o)
	vel := RadiansVec2(o).Multiply(100)
	s.SetVelocity(vel)
	d.state = "flee"
}

func (d *AIDrunkard) Idle(w *World, s *Sprite) {
	if d.state == "idle" {
		return
	}
	d.state = "idle"
	s.SetVelocity(&Vec2{0, 0})
}

func (d *AIDrunkard) Stagger(w *World, s *Sprite) {
	d.state = "stagger"
	if w.rand.Float32() < 0.99 {
		return
	}
	rand := w.rand.Float32()
	if rand > 0.75 {
		s.Orientation = (math.Pi / 2)
	} else if rand > 0.50 {
		s.Orientation = (3 * math.Pi / 2)
	} else if rand > 0.25 {
		s.Orientation = math.Pi
	} else {
		s.Orientation = 0
	}
	vel := RadiansVec2(s.Orientation).Multiply(10)
	s.velocity.X = vel.X
	s.velocity.Y = vel.Y
}
