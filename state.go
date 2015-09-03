package main

import (
	"time"
)

type Condition interface {
	Test() bool
}

type Actioner interface {
	Front() Actioner
	Back() Actioner
	Next() Actioner
	SetNext(Actioner)
	Act()
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

func (a *Action) Act() {

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

// StateMachineState
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
			// Transition to the first state
			sm.currentState = sm.initialState
			// return initial state's actions
			return sm.initialState.EntryActions()
		}
	}

	// We have a current state to work with

	// Start off with no transition
	var transition *Transition

	// Check through each transition in the current state.
	testTransition := sm.currentState.Transition()
	for t := testTransition.Front(); t != nil; t = t.Next() {
		// check if this transition should trigger
		if testTransition.isTriggered() {
			transition = testTransition
			break
		}
	}

	// no transitions, so return current states actions
	if transition == nil {
		return sm.currentState.Actions()
	}

	// we are going through a transition so create a list of actions from the
	// - current states exit action(s)
	// - the transitions action(s)
	// - the next states action(s)

	// get all the exit actions from the current state
	actions := sm.currentState.ExitActions()
	last := actions.Back()
	// get the actions from the transition
	transitionActions := transition.Actions()
	last.SetNext(transitionActions)
	// add the actions from the next state
	nextState := transition.getTargetState()
	transitionActions.Back().SetNext(nextState.EntryActions())
	// Update the change of state
	sm.currentState = nextState
	return actions
}
