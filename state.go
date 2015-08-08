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
	//	log.Printf("State changing from %d to %d after %s", sm.state, state, sm.Duration())
	sm.state = state
	sm.timer = time.Now()
}

func (sm *StateMachine) State() State {
	return sm.state
}

func (sm *StateMachine) Duration() time.Duration {
	return time.Now().Sub(sm.timer)
}
