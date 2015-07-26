package main

import (
	"log"
	"math"
	"time"
)

type AI interface {
	Update(s *Sprite, w *World, t time.Duration)
}

type AIDrunkard struct {
	state string
}

func (d *AIDrunkard) Update(s *Sprite, w *World, t time.Duration) {
	for _, input := range s.inputs {
		switch input.Action {
		case "mouseOver":
			d.Flee(w, s)
			return
		case "mouseOut":
			d.Idle(w, s)
			return
		case "clickUp":
			log.Printf("kill me %d", s.Id)
			s.Kill()
			return
		}
	}
	d.Stagger(w, s)
}

func (d *AIDrunkard) Flee(w *World, s *Sprite) {
	if d.state != "flee" {
		rot := w.rand.Float64() * math.Pi * 2
		s.SetRotation(rot)
		vel := RadiansVec2(rot).Multiply(50)
		s.SetVelocity(vel.X, vel.Y)
		d.state = "flee"
	}
}

func (d *AIDrunkard) Idle(w *World, s *Sprite) {
	d.state = "idle"
	s.SetVelocity(0, 0)
}

func (d *AIDrunkard) Stagger(w *World, s *Sprite) {
	d.state = "stagger"
	if w.rand.Float32() < 0.99 {
		return
	}
	rand := w.rand.Float32()
	if rand > 0.75 {
		s.Rotation = (3.14 / 2)
	} else if rand > 0.50 {
		s.Rotation = (3 * 3.14 / 2)
	} else if rand > 0.25 {
		s.Rotation = 3.14
	} else {
		s.Rotation = 0
	}
	vel := RadiansVec2(s.Rotation).Multiply(10)
	s.SetVelocity(vel.X, vel.Y)
}
