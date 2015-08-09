package main

import "time"

func NewStateMachine(initial State) *StateMachine {
	sm := &StateMachine{
		timer: time.Now(),
	}
	sm.SetState(initial)
	return sm
}

type StateMachine struct {
	state State
	timer time.Time
}

func (sm *StateMachine) SetState(state State) {
	sm.state = state
	sm.timer = time.Now()
}

func (sm *StateMachine) State() State {
	return sm.state
}

func (sm *StateMachine) Duration() time.Duration {
	return time.Now().Sub(sm.timer)
}
