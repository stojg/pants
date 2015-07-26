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
}

func (d *AIDrunkard) Update(s *Sprite, w *World, t time.Duration) {
	for _, input := range s.inputs {
		switch input.Action {
		case "mouseOver":
//			d.Flee(w, s)
//			return
		case "mouseOut":
			d.Stop(w, s)
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
	x, y := s.Velocity()
	if math.Abs(x+y) == 0 {
		if w.rand.Float32() < 0.5 {
			x = 10
		} else {
			y = 10
		}
	}
	s.SetVelocity(x*10, y*10)
}

func (d *AIDrunkard) Stop(w *World, s *Sprite) {
	s.SetVelocity(0, 0)
}

func (d *AIDrunkard) Stagger(w *World, s *Sprite) {
	if w.rand.Float32() < 0.99 {
		return
	}
	rand := w.rand.Float32()
	if rand > 0.75 {
		s.Rotation = (3.14 / 2)
		s.SetVelocity(10, 0)
	} else if rand > 0.50 {
		s.Rotation = (3 * 3.14 / 2)
		s.SetVelocity(-10, 0)
	} else if rand > 0.25 {
		s.Rotation = 3.14
		s.SetVelocity(0, 10)
	} else {
		s.Rotation = 0
		s.SetVelocity(0, -10)
	}
}
