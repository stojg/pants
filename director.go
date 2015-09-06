package main

import "log"

type Director struct {
	lastAction float64
}

func NewDirector() *Director {
	d := &Director{}
	events.Register(EVT_COLLISION, d)
	events.Register(EVT_CLICK, d)
	return d
}

func (director *Director) Update(w *World, frameTime, gameTime float64) {

	if gameTime-director.lastAction < 0.1 {
		return
	}

	if len(entityManager.entities) >= 400 {
		return
	}

	entType := ENT_ENT1
	rand := w.rand.Float64()
	if rand > 0.5 {
		entType = ENT_FOOD
	}

	director.lastAction = gameTime
	id := w.entityManager.NewEntity(
		w.rand.Float64()*1000,
		w.rand.Float64()*1000,
		entType,
	)
	evt := &EntityCreateEvent{id: id}
	events.ScheduleEvent(evt)
}

func (director *Director) Notify(evt Event) {
	switch evt := evt.(type) {
	default:
		log.Printf("unknown event type %T", evt)
	case *CollisionEvent:
		a := world.entityManager.entities[evt.a]
		b := world.entityManager.entities[evt.b]
		if a.Type == ENT_FOOD && b.Type == ENT_ENT1 {
			world.entityManager.Remove(evt.a)
			pB := world.entityManager.physics[evt.b]
			pB.Data.Height += 3
			pB.Data.Width += 3
			pB.Data.InvMass = 1 / (1/pB.Data.InvMass + 2.0)
		}
		if b.Type == ENT_FOOD && a.Type == ENT_ENT1 {
			world.entityManager.Remove(evt.b)
			pA := world.entityManager.physics[evt.a]
			pA.Data.Height += 3
			pA.Data.Width += 3
			pA.Data.InvMass = 1 / (1/pA.Data.InvMass + 2.0)
		}
	case *ClickEvent:
		log.Printf("Asdsad")
		for id := range world.entityManager.entities {
			world.entityManager.Remove(id)
		}
	}
}

type EntityCreateEvent struct {
	id uint64
}

func (evt *EntityCreateEvent) Code() EventCode {
	return EVT_NEW_ENTITY
}
