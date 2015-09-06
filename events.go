package main

type Listener interface {
	Notify(evt Event)
}

const (
	EVT_NEW_ENTITY EventCode = iota
	EVT_COLLISION
	EVT_CLICK
)

type EventCode byte

type Event interface {
	Code() EventCode
}

func NewEventManager() *EventManager {
	return &EventManager{
		listenerRegistrations: make([]*ListenerRegistration, 0),
		events:                make([]Event, 0),
	}
}

type EventManager struct {
	listenerRegistrations []*ListenerRegistration
	events                []Event
}

type ListenerRegistration struct {
	code     EventCode
	listener Listener
}

func (e *EventManager) Register(code EventCode, listener Listener) {
	e.listenerRegistrations = append(e.listenerRegistrations, &ListenerRegistration{
		code:     code,
		listener: listener,
	})
}

func (e *EventManager) ScheduleEvent(evt Event) {
	e.events = append(e.events, evt)
}

func (e *EventManager) DispatchEvents() {
	// @todo: not really thread safe.. not at all
	for _, event := range e.events {
		for _, listener := range e.listenerRegistrations {
			if listener.code == event.Code() {
				listener.listener.Notify(event)
			}
		}
	}
	e.events = make([]Event, 0)
}
