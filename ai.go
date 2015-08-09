package main

import (
	"time"
	. "github.com/stojg/pants/physics"
)

type State int

const (
	STATE_IDLE State = iota // 0
	STATE_FLEE              // 1
	STATE_HUNT              // 2
)

type AI interface {
	Update(s *Sprite, w *World, t float64)
}

type AIDrunkard struct {
	state    *StateMachine
	steering Steering
}

func (d *AIDrunkard) Update(s *Sprite, w *World, t float64) {
	d.handleInputs(s)
	d.handleState(s, w, t)
}

func (d *AIDrunkard) handleInputs(s *Sprite) {
	for _, input := range s.inputs {
		switch input.Action {
		case "mouseOver":
			d.state.SetState(STATE_FLEE)
			//		case "mouseOut":
			//			d.state.SetState(STATE_IDLE)
		case "clickUp":
			d.state.SetState(STATE_IDLE)
		}
	}
	s.inputs = nil
}

func (d *AIDrunkard) handleState(s *Sprite, w *World, t float64) {
	var o *SteeringOutput
	switch d.state.State() {
	case STATE_IDLE:
		o = d.idle(w, s, t)
	case STATE_FLEE:
		o = d.flee(w, s, t)
	case STATE_HUNT:
		o = d.hunt(w, s, t)
	}
	if o != nil {
		w.Physic(s.Id).AddForce(o.linear.Scale(w.Physic(s.Id).Mass()))
		w.Physic(s.Id).SetRotations(o.angular)
	}
}

func (d *AIDrunkard) flee(w *World, s *Sprite, duration float64) *SteeringOutput {
	if d.state.Duration() > 1*time.Second {
		d.steering = nil
		d.state.SetState(STATE_IDLE)
		return nil
	}

	if _, ok := d.steering.(*Flee); !ok {
		d.steering = NewFlee(w.Physic(s.Id), NewPhysics(w.RandF64(800), w.RandF64(600), 0, 0))
	}
	return d.steering.Get(duration)
}

func (d *AIDrunkard) idle(w *World, s *Sprite, duration float64) *SteeringOutput {
	if d.state.Duration() > 3*time.Second {
		d.steering = nil
		d.state.SetState(STATE_HUNT)
		return nil
	}
	if _, ok := d.steering.(*Arrive); !ok {
		d.steering = NewArrive(w.Physic(s.Id), w.Physic(s.Id).Clone(), 1, 500)
	}
	return d.steering.Get(duration)
}

func (d *AIDrunkard) hunt(w *World, s *Sprite, duration float64) *SteeringOutput {
	if _, ok := d.steering.(*Arrive); !ok {
		d.steering = NewArrive(w.Physic(s.Id), NewPhysics(w.RandF64(800), w.RandF64(600), 0, 0), 1, 500)
	}

	direction := d.steering.Target().Position.Clone().Sub(w.Physic(s.Id).Position)
	length := direction.Length()
	if length < 2 {
		d.steering = nil
		d.state.SetState(STATE_IDLE)
		return nil
	}
	return d.steering.Get(duration)
}
