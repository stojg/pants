package main

import (
	"time"
	"log"
)

type Condition interface {
	Test() bool
}

type Actioner interface {
	Front() Actioner
	Back() Actioner
	Next() Actioner
	SetNext(Actioner)
	Act(id uint64)
	Interruptable() bool
	CanDoBoth(other Actioner) bool
	IsComplete() bool
}

type Action struct {
	next Actioner
}

func (a *Action) SetNext(o Actioner) {
	a.next = o
}

func (a *Action) Front() Actioner {
	return a
}

func (a *Action) Next() Actioner {
	return a.next
}

func (a *Action) Back() Actioner {
	if a.Next() == nil {
		return a
	}
	var thisAction Actioner = a
	var nextAction Actioner = a.next
	for nextAction != nil {
		thisAction = a.next
		nextAction = a.Next().Next()
	}
	return thisAction
}

func (a *Action) Act(id uint64) {

}

func (a *Action) Interruptable() bool {
	return false
}

func (a *Action) CanDoBoth(other Actioner) bool {
	return false
}

func (a *Action) IsComplete() bool {
	return true
}

type Transition struct {
	condition   Condition
	targetState StateMachineState
	action      Actioner
	next        *Transition
}

func (t *Transition) Front() *Transition {
	return t
}

func (t *Transition) Next() *Transition {
	return t.next
}

func (t *Transition) Actions() Actioner {
	return t.action
}

func (t *Transition) isTriggered() bool {
	return t.condition.Test()
}

func (t *Transition) getTargetState() StateMachineState {
	return t.targetState
}

type StateMachineState interface {
	// Transition gets the first... transition...
	Transition() *Transition
	// Action returns the first in a sequence of actions that should be
	// performed while the character is in this state.
	// Note that this method should return one or more newly created action
	// instances, and the caller of this method should be responsible for the
	// deletion.
	Actions() Actioner
	// EntryActions returns the sequence of actions to perform when arriving in
	// this state.
	// Note that this method should return one or more newly created action
	// instances, and the caller of this method should be responsible for the
	// deletion.
	EntryActions() Actioner
	// ExitActions returns the sequence of actions to perform when leaving
	// this state.
	// Note that this method should return one or more newly created action
	// instances, and the caller of this methodshould be responsible for the
	// deletion.
	ExitActions() Actioner
}

// StateMachineState
type State struct {
	transition *Transition
}

// Transition gets the first... transition...
func (s *State) Transition() *Transition {
	return s.transition
}

// Action returns the first in a sequence of actions that should be
// performed while the character is in this state.
// Note that this method should return one or more newly created action
// instances, and the caller of this method should be responsible for the
// deletion.
func (s *State) Actions() Actioner {
	return nil
}

// EntryActions returns the sequence of actions to perform when arriving in
// this state.
// Note that this method should return one or more newly created action
// instances, and the caller of this method should be responsible for the
// deletion.
func (s *State) EntryActions() Actioner {
	return nil
}

// ExitActions returns the sequence of actions to perform when leaving
// this state.
// Note that this method should return one or more newly created action
// instances, and the caller of this methodshould be responsible for the
// deletion.
func (s *State) ExitActions() Actioner {
	return nil
}

func NewStateMachine(initial StateMachineState) *StateMachine {
	sm := &StateMachine{
		timer: time.Now(),
	}
	sm.initialState = initial
	return sm
}

// StateMachine encapsulates a single layer state machine (i.e. one without
// hierarchical transitions)
type StateMachine struct {
	initialState StateMachineState
	currentState StateMachineState
	timer        time.Time
}

// Update returns an Action
func (sm *StateMachine) Update() Actioner {
	// first case - we have no current state
	if sm.currentState == nil {
		// in this case we use the entry action for the initial state
		if sm.initialState == nil {
			// we have nothing to do
			return nil
		} else {
			log.Printf("initial state %T", sm.initialState)
			// Transition to the first state
			sm.currentState = sm.initialState
			// return initial state's actions
			return sm.initialState.EntryActions()
		}
	}

	// Start off with no transition
	var transition *Transition

	// Check through each transition in the current state.
	testTransition := sm.currentState.Transition()
//	log.Printf("transition %V", testTransition)
	for t := testTransition.Front(); t != nil; t = t.Next() {
		// check if this transition should trigger
		if testTransition.isTriggered() {
			log.Printf("transition %T triggered", testTransition)
			transition = testTransition
			break
		}
	}

	// no transitions, so return current states actions
	if transition == nil {
		return sm.currentState.Actions()
	}

	log.Printf("transitioning from %T to %T", sm.currentState, transition.targetState)

	// we are going through a transition so create a list of actions from the
	// - current states exit action(s)
	// - the transitions action(s)
	// - the next states action(s)
	actions := make([]Actioner,0)
	// get all the exit actions from the current state
	exitActions := sm.currentState.ExitActions()
	if exitActions != nil {
		actions = append(actions, exitActions)
	}

	// get the actions from the transition
	transitionActions := transition.Actions()
	if transitionActions != nil {
		actions = append(actions, transitionActions)
	}


	// add the actions from the next state
	nextState := transition.getTargetState()
	entryActions := nextState.EntryActions()
	if entryActions != nil {
		actions = append(actions, entryActions)
	}

	sm.currentState = nextState

	if len(actions) == 0 {
		return nil
	}

	first, actions := actions[0], actions[1:]

	for _, action := range actions {
		first.Back().SetNext(action)
	}

	// Update the change of state
	return first
}
