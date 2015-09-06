package main

import (
	. "github.com/stojg/pants/ai"
	"github.com/stojg/pants/physics"
	"github.com/stojg/pants/timer"
)

type Controller interface {
	Update(id uint64, w *World, t float64)
}

/// Conditions
type TrueCondition struct{}

func (c *TrueCondition) Test() bool {
	return true
}

type FalseCondition struct {
}

func (c *FalseCondition) Test() bool {
	return false
}

type Hunt struct {
	Action
	steering Steering
}

func (a *Hunt) Act(id uint64) {
	o := a.steering.Get()
	world.Physic(id).Data.Forces.Add(o.Linear.Scale(world.Physic(id).Data.Mass()))
	world.Physic(id).Data.Rotations = o.Angular
}

// Concrete state implementation
type IdleState struct {
	State
}

type HuntingState struct {
	State
	id       uint64
	steering Steering
}

func (m *HuntingState) Actions() []Actioner {
	actions := make([]Actioner, 1)
	actions = append(actions, &Hunt{
		steering: m.steering,
	})
	return actions
}

func (m *HuntingState) EntryActions() []Actioner {
	target := physics.NewPhysics(world.RandF64(800), world.RandF64(600), 0, 0, 0, 0)
	m.steering = NewArrive(world.Physic(m.id), target, 10, 500)
	return nil
}

type BasicAI struct {
	stateMachine *StateMachine
}

type TimerCondition struct {
	timer *timer.Timer
	value float64
}

func (c *TimerCondition) Test() bool {
	return c.timer.Seconds() > c.value
}

// AI implementation
func NewBasicAI(id uint64) *BasicAI {

	idle := &IdleState{}
	hunting := &HuntingState{id: id}

	ai := &BasicAI{stateMachine: NewStateMachine(idle)}

	idle.AddTransitions(NewTransition(hunting, &TimerCondition{
		value: 3.0,
		timer: ai.stateMachine.timer,
	}))
	hunting.AddTransitions(NewTransition(idle, &FalseCondition{}))

	return ai
}

func (d *BasicAI) Update(id uint64, w *World, t float64) {
	for _, a := range d.stateMachine.Update() {
		if a != nil {
			a.Act(id)
		}
	}
}
