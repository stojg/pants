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
		case "mouseOut":
			d.Idle(w, s)
		case "clickUp":
			s.Kill()
		}
	}
	if len(s.inputs) == 0 {
		d.Stagger(w, s)
	}
	s.inputs = s.inputs[0:0]
}

func (d *AIDrunkard) Flee(w *World, s *Sprite) {
	if d.state == "flee" {
		return
	}
	o := w.rand.Float64() * math.Pi * 2
	list.physics[s.Id].Orientation = o
	list.physics[s.Id].Velocity = RadiansVec2(o).Multiply(100)
	d.state = "flee"
}

func (d *AIDrunkard) Idle(w *World, s *Sprite) {
	if d.state == "idle" {
		return
	}
	d.state = "idle"
	list.physics[s.Id].Velocity = (&Vec2{0, 0})
}

func (d *AIDrunkard) Stagger(w *World, s *Sprite) {
	d.state = "stagger"
	if w.rand.Float32() < 0.99 {
		return
	}
	rand := w.rand.Float32()
	if rand > 0.75 {
		list.physics[s.Id].Orientation = (math.Pi / 2)
	} else if rand > 0.50 {
		list.physics[s.Id].Orientation = (3 * math.Pi / 2)
	} else if rand > 0.25 {
		list.physics[s.Id].Orientation = math.Pi
	} else {
		list.physics[s.Id].Orientation = 0
	}
	vel := RadiansVec2(list.physics[s.Id].Orientation).Multiply(10)
	list.physics[s.Id].Velocity = vel
}
