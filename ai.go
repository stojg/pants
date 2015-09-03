package main

import (
	. "github.com/stojg/pants/ai"
	//	. "github.com/stojg/pants/physics"
	//	"time"
)

type AI interface {
	Update(id uint64, w *World, t float64)
}

type BasicAI struct {
	state    *StateMachine
	steering Steering
}

/// Conditions
type Timer struct{}

func (c *Timer) Test() bool {
	return false
}

type CloseClick struct{}

func (c *CloseClick) Test() bool {
	return false
}

type IsSafe struct{}

func (c *IsSafe) Test() bool {
	return false
}

type IsCloseTo struct{}

func (c *IsCloseTo) Test() bool {
	return false
}

// Actions

type MyAction struct {
	Action
}

func (a *MyAction) Act() {
	// do something here
}

// Concrete state implementation

type SomeState struct {
	transition *Transition
}

func (m *SomeState) Actions() Actioner {
	return &MyAction{}
}

func (m *SomeState) ExitActions() Actioner {
	return nil
}

func (m *SomeState) EntryActions() Actioner {
	return nil
}

func (m *SomeState) Transition() *Transition {
	return m.transition
}

// AI implementation

func NewBasicAI() *BasicAI {

	idleState := &SomeState{}
	seekState := &SomeState{}
	fleeState := &SomeState{}

	seekTransition := &Transition{
		condition:   &Timer{},
		targetState: seekState,
		action:      nil,
		next:        nil, // next transition
	}

	fleeTransition := &Transition{
		condition:   &CloseClick{},
		targetState: fleeState,
		action:      nil,
		next:        nil,
	}

	stopFleeingTransition := &Transition{
		condition:   &IsSafe{},
		targetState: seekState,
		action:      nil,
	}

	idleTransition := &Transition{
		condition:   &IsCloseTo{},
		targetState: idleState,
		action:      nil,
	}

	seekTransition.next = fleeTransition
	idleState.transition = seekTransition
	seekState.transition = idleTransition
	fleeState.transition = stopFleeingTransition

	return &BasicAI{
		state: NewStateMachine(idleState),
	}
}

func (d *BasicAI) Update(id uint64, w *World, t float64) {
	//	d.handleInputs(s)
	action := d.state.Update()
	if action != nil {
		for a := action.Front(); a != nil; a = action.Next() {
			a.Act()
		}
	}
}

func (d *BasicAI) handleInputs(s *Entity) {
	for _, input := range s.inputs {
		switch input.Type {
		case "mouseOver":
		//			d.state.SetState(STATE_FLEE)
		//		case "mouseOut":
		//			d.state.SetState(STATE_IDLE)
		case "clickUp":
			//			d.state.SetState(STATE_IDLE)
		}
	}
	s.inputs = nil
}

func (d *BasicAI) flee(w *World, id uint64, duration float64) *SteeringOutput {
	return nil
	//	if d.state.Past(time.Second) {
	//		d.steering = nil
	//		d.state.SetState(STATE_IDLE)
	//		return nil
	//	}
	//
	//	if _, ok := d.steering.(*Flee); !ok {
	//		d.steering = NewFlee(w.Physic(id), NewPhysics(w.RandF64(1000), w.RandF64(1000), 0, 0))
	//	}
	//	return d.steering.Get(duration)
}

func (d *BasicAI) idle(w *World, id uint64, duration float64) *SteeringOutput {
	return nil
	//	if !d.state.Past(3*time.Second) {
	//		return nil
	//	}
	//	d.steering = nil
	//	d.state.SetState(STATE_HUNT)
	//	return nil
}

func (d *BasicAI) hunt(w *World, id uint64, duration float64) *SteeringOutput {
	return nil
	//	if _, ok := d.steering.(*Arrive); !ok {
	//		d.steering = NewArrive(w.Physic(id), NewPhysics(w.RandF64(800), w.RandF64(600), 0, 0), 10, 500)
	//	}
	//
	//	direction := d.steering.Target().Position.Clone().Sub(w.Physic(id).Position)
	//	length := direction.Length()
	//	if length <= 10 {
	//		d.steering = nil
	//		d.state.SetState(STATE_IDLE)
	//		return nil
	//	}
	//	return d.steering.Get(duration)
}
