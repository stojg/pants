package main

import (
	. "github.com/stojg/pants/ai"
	"github.com/stojg/pants/physics"
	"github.com/stojg/pants/vector"
)

type AI interface {
	Update(id uint64, w *World, t float64)
}

type BasicAI struct {
	state    *StateMachine
	steering Steering
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
	world.Physic(id).AddForce(o.Linear.Scale(world.Physic(id).Mass()))
	world.Physic(id).SetRotations(o.Angular)
}

// Concrete state implementation
type IdleState struct {
	State
}

type HuntingState struct {
	State
	id       uint64
	position *vector.Vec2
	steering Steering
}

func (m *HuntingState) Actions() Actioner {
	return &Hunt{
		steering: m.steering,
	}
}

func (m *HuntingState) EntryActions() Actioner {
	m.position = &vector.Vec2{world.RandF64(800), world.RandF64(600)}
	target := physics.NewPhysics(m.position.X, m.position.Y, 0, 0)
	m.steering = NewArrive(world.Physic(m.id), target, 10, 500)
	return nil
}

// AI implementation
func NewBasicAI(id uint64) *BasicAI {

	idleState := &IdleState{}
	hunting := &HuntingState{id: id}

	huntingTransition := &Transition{
		condition:   &TrueCondition{},
		targetState: hunting,
	}

	idleTransition := &Transition{
		condition: &FalseCondition{},
		targetState: idleState,
	}

	idleState.transition = huntingTransition
	hunting.transition = idleTransition

	return &BasicAI{
		state: NewStateMachine(idleState),
	}
}

func (d *BasicAI) Update(id uint64, w *World, t float64) {
	//	d.handleInputs(s)
	for a := d.state.Update(); a != nil; a = a.Next() {
		if a != nil {
			a.Act(id)
		}
	}
}

func (d *BasicAI) handleInputs(s *Entity) {
}
