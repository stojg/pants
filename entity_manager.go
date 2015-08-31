package main

import (
	. "github.com/stojg/pants/physics"
	. "github.com/stojg/pants/vector"
	"github.com/stojg/tree"
)

func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities:         make(map[uint64]*Entity),
		ais:              make(map[uint64]AI),
		physics:          make(map[uint64]*Physics),
		boundingBoxes:    make(map[uint64]*tree.Rectangle),
		updated:          make(map[uint64]bool),
		forceRegistry:    &PhysicsForceRegistry{},
		collisionManager: &CollisionManager{},
	}
}

type EntityManager struct {
	lastEntityID     uint64
	entities         map[uint64]*Entity
	ais              map[uint64]AI
	physics          map[uint64]*Physics
	boundingBoxes    map[uint64]*tree.Rectangle
	updated          map[uint64]bool
	forceRegistry    *PhysicsForceRegistry
	collisionManager *CollisionManager
}

var gravity *StaticForce

func init() {
	gravity = &StaticForce{}
	gravity.SetForce(&Vec2{0, 10})
}

func (em *EntityManager) NewEntity(x, y float64, image string) uint64 {

	em.lastEntityID += 1;
	entity := &Entity{
		Id: em.lastEntityID,
		Type: "sprite",
		Image: image,
		inputs: make([]*InputRequest, 0),
	}

	em.entities[entity.Id] = entity
	em.ais[entity.Id] = &AIDrunkard{state: NewStateMachine(STATE_IDLE)}
	em.physics[entity.Id] = NewPhysics(x, y, 3.14*2, 1)
	em.physics[entity.Id].SetDamping(0.99)

	em.updated[entity.Id] = true

	em.forceRegistry.Add(em.physics[entity.Id], gravity)
	em.collisionManager.Add(em.physics[entity.Id])
	return entity.Id
}

func (em *EntityManager) Update(w *World, duration float64) {

	em.forceRegistry.Update(duration)

	for id, ai := range em.ais {
		ai.Update(id, w, duration)
	}

	for id, p := range em.physics {
		moved := p.Update(id, w, duration)
		if moved {
			em.updated[id] = true
		}
	}

	em.collisionManager.DetectCollisions()
	em.collisionManager.ResolveCollisions(duration)
}

func (em *EntityManager) Changed() []*EntityUpdate {
	// @todo(stig): it is possible that this might be called before the em.sprite has been fully setup, ie a race condition
	ids := make([]uint64, len(em.entities))
	i := 0
	for k := range em.entities {
		ids[i] = k
		i++
	}
	return em.entityUpdates(ids)
}

func (em *EntityManager) entityUpdates(entityIds []uint64) []*EntityUpdate {
	entities := make([]*EntityUpdate, 0, len(entityIds))
	for _, id := range entityIds {
		entities = append(entities, &EntityUpdate{
			Id:          id,
			X:           em.physics[id].Position.X,
			Y:           em.physics[id].Position.Y,
			Orientation: em.physics[id].Orientation,
			Type:        em.entities[id].Type,
			Data: map[string]string{
				"Image": em.entities[id].Image,
			},
		})
	}
	return entities
}
