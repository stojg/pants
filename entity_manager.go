package main

import (
	. "github.com/stojg/pants/physics"
	. "github.com/stojg/pants/vector"
	"github.com/stojg/tree"
	"math"
	"github.com/stojg/pants/structs"
)

func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities:         make(map[uint64]*Entity),
		controllers:      make(map[uint64]Controller),
		physics:          make(map[uint64]*Physics),
		boundingBoxes:    make(map[uint64]*tree.Rectangle),
		updated:          make(map[uint64]bool),
		forceRegistry:    &PhysicsForceRegistry{},
		collisionManager: NewCollisionManager(),
	}
}

type EntityManager struct {
	lastEntityID     uint64
	entities         map[uint64]*Entity
	controllers      map[uint64]Controller
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

func (em *EntityManager) NewEntity(x, y float64, entType EntityType) uint64 {

	em.lastEntityID += 1

	// copy the properties from the data file
	props := &EntityProperty{}
	*props = *entityProperties[entType]

	entity := &Entity{
		Id:         em.lastEntityID,
		Type:       entType,
		inputs:     make([]*InputRequest, 0),
		Properties: props,
	}

	em.entities[entity.Id] = entity
	data := structs.NewData()
	data.Position.X = x
	data.Position.Y = y
	data.Orientation = 3.14 * 2
	data.InvMass = props.invMass
	data.Width = props.Width
	data.Height = props.Height
	em.physics[entity.Id] = NewPhysics(data)
	em.physics[entity.Id].Data.Damping = 0.99

	//	em.forceRegistry.Add(em.physics[entity.Id], gravity)
	em.collisionManager.Add(em.physics[entity.Id], entity.Id)

	if props.ai {
		em.controllers[entity.Id] = NewBasicAI(entity.Id)
	}
	em.updated[entity.Id] = true
	return entity.Id
}

func (em *EntityManager) Remove(id uint64) {
	em.entities[id].Dead = true
	delete(em.physics, id)
	delete(em.controllers, id)
	em.collisionManager.Remove(id)
	em.updated[id] = true
}

func (em *EntityManager) Update(w *World, duration float64) {

	em.forceRegistry.Update(duration)

	for id, ctrl := range em.controllers {
		ctrl.Update(id)
	}

	for id, p := range em.physics {
		data := p.Data
		p.Integrate(data, duration)

		// update collision volume
		if data.Width > data.Height {
			p.CollisionGeometry.(*Circle).Radius = data.Width / 2
		} else {
			p.CollisionGeometry.(*Circle).Radius = data.Height / 2
		}

		moved := false
		if !data.Position.Equals(p.PrevData.Position) {
			moved = true
		}
		if math.Abs(data.Orientation-p.PrevData.Orientation) < EPSILON {
			moved = true
		}
		p.PrevData.Copy(data)

		if moved {
			em.updated[id] = true
		}
	}

	em.collisionManager.DetectCollisions()
	for _, collision := range em.collisionManager.Collisions() {
		evt := &CollisionEvent{}
		evt.a, evt.b = collision.Pair()
		events.ScheduleEvent(evt)
	}

	em.collisionManager.ResolveCollisions(duration)
}

func (em *EntityManager) Changed() []*EntityUpdate {
	// @todo(stig): it is possible that this might be called before the em.sprite has been fully setup, ie a race condition
	ids := make([]uint64, 0)
	for k := range em.entities {
		ids = append(ids, k)
	}
	return em.entityUpdates(ids)
}

func (em *EntityManager) entityUpdates(entityIds []uint64) []*EntityUpdate {
	entities := make([]*EntityUpdate, 0, len(entityIds))
	for _, id := range entityIds {

		update := &EntityUpdate{
			Id:         id,
			Type:       em.entities[id].Type,
			Properties: em.entities[id].Properties,
			Dead:       em.entities[id].Dead,
		}
		if physics, ok := em.physics[id]; ok {
			update.Height = physics.Data.Height
			update.Width = physics.Data.Width
			update.X = physics.Data.Position.X
			update.Y = physics.Data.Position.Y
			update.Orientation = physics.Data.Orientation
		}
		entities = append(entities, update)
	}
	return entities
}

type CollisionEvent struct {
	a uint64
	b uint64
}

func (c *CollisionEvent) Code() EventCode {
	return EVT_COLLISION
}
