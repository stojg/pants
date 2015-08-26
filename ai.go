package main

import (
	. "github.com/stojg/pants/ai"
	. "github.com/stojg/pants/physics"
	"time"
)

type State int

const (
	STATE_IDLE State = iota // 0
	STATE_FLEE              // 1
	STATE_HUNT              // 2
)

type AI interface {
	Update(id uint64, w *World, t float64)
}

type AIDrunkard struct {
	state    *StateMachine
	steering Steering
}

func (d *AIDrunkard) Update(id uint64, w *World, t float64) {
//	d.handleInputs(s)
	d.handleState(id, w, t)
}

func (d *AIDrunkard) handleInputs(s *Entity) {
	for _, input := range s.inputs {
		switch input.Type {
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

func (d *AIDrunkard) handleState(id uint64, w *World, t float64) {
	var o *SteeringOutput
	switch d.state.State() {
	case STATE_IDLE:
		o = d.idle(w, id, t)
	case STATE_FLEE:
		o = d.flee(w, id, t)
	case STATE_HUNT:
		o = d.hunt(w, id, t)
	}
	if o != nil {
		w.Physic(id).AddForce(o.Linear.Scale(w.Physic(id).Mass()))
		w.Physic(id).SetRotations(o.Angular)
	}
}

func (d *AIDrunkard) flee(w *World, id uint64, duration float64) *SteeringOutput {
	if d.state.Duration() > 1*time.Second {
		d.steering = nil
		d.state.SetState(STATE_IDLE)
		return nil
	}

	if _, ok := d.steering.(*Flee); !ok {
		d.steering = NewFlee(w.Physic(id), NewPhysics(w.RandF64(1000), w.RandF64(1000), 0, 0))
	}
	return d.steering.Get(duration)
}

func (d *AIDrunkard) idle(w *World, id uint64, duration float64) *SteeringOutput {
	if d.state.Duration() > 3*time.Second {
		d.steering = nil
		d.state.SetState(STATE_HUNT)
		return nil
	}
	return nil
	if _, ok := d.steering.(*Arrive); !ok {
		d.steering = NewArrive(w.Physic(id), w.Physic(id).Clone(), 1, 500)
	}
	return d.steering.Get(duration)
}

func (d *AIDrunkard) hunt(w *World, id uint64, duration float64) *SteeringOutput {
	if _, ok := d.steering.(*Arrive); !ok {
		d.steering = NewArrive(w.Physic(id), NewPhysics(w.RandF64(800), w.RandF64(600), 0, 0), 10, 500)
	}

	direction := d.steering.Target().Position.Clone().Sub(w.Physic(id).Position)
	length := direction.Length()
	if length <= 10 {
		d.steering = nil
		d.state.SetState(STATE_IDLE)
		return nil
	}
	return d.steering.Get(duration)
}
