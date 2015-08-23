package main

import (
	. "github.com/stojg/pants/physics"
	. "github.com/stojg/pants/vector"
	"github.com/stojg/tree"
)

func NewEntityManager() *EntityManager {
	return &EntityManager{
		sprites:          make(map[uint64]*Entity),
		ais:              make(map[uint64]AI),
		physics:          make(map[uint64]*Physics),
		boundingBoxes:    make(map[uint64]*tree.Rectangle),
		updated:          make(map[uint64]bool),
		forceRegistry:    &PhysicsForceRegistry{},
		collisionManager: &CollisionManager{},
		quadTree:         tree.NewQuadTree(0, tree.NewRectangle(500, 500, 500, 500)),
	}
}

type EntityManager struct {
	lastEntityID     uint64
	sprites          map[uint64]*Entity
	ais              map[uint64]AI
	physics          map[uint64]*Physics
	boundingBoxes    map[uint64]*tree.Rectangle
	updated          map[uint64]bool
	forceRegistry    *PhysicsForceRegistry
	collisionManager *CollisionManager
	quadTree         *tree.QuadTree
}

var gravity *StaticForce

func init() {
	gravity = &StaticForce{}
	gravity.SetForce(&Vec2{0, 10})
}

func (em *EntityManager) NewEntity(x, y float64, image string) uint64 {
	sprite := &Entity{}
	sprite.Dead = false
	sprite.Type = "sprite"
	sprite.Image = image
	sprite.inputs = make([]*InputRequest, 0)
	em.lastEntityID++
	sprite.Id = em.lastEntityID
	em.sprites[sprite.Id] = sprite
	em.ais[sprite.Id] = &AIDrunkard{state: NewStateMachine(STATE_IDLE)}
	em.physics[sprite.Id] = NewPhysics(x, y, 3.14*2, 1)
	em.physics[sprite.Id].SetDamping(0.99)
	em.forceRegistry.Add(em.physics[sprite.Id], gravity)
	em.updated[sprite.Id] = true
	em.collisionManager.Add(em.physics[sprite.Id])
	return sprite.Id
}

func (em *EntityManager) Update(w *World, duration float64) {

	em.forceRegistry.Update(duration)

	for _, sprite := range em.sprites {
		em.ais[sprite.Id].Update(sprite, w, duration)
	}

	for _, sprite := range em.sprites {
		moved := em.physics[sprite.Id].Update(sprite.Id, w, duration)
		if moved {
			em.updated[sprite.Id] = true
		}
	}

	tries := 0
	if em.collisionManager.DetectCollisions() || tries > 3 {
		tries += 1
		em.collisionManager.ResolveCollisions(duration)
	}
}

func (em *EntityManager) All() []*EntityUpdate {
	ids := make([]uint64, len(em.sprites))
	i := 0
	for k := range em.sprites {
		ids[i] = k
		i++
	}

	return em.entityUpdates(ids)
}

func (em *EntityManager) Changed() []*EntityUpdate {
	ids := make([]uint64, len(em.sprites))
	i := 0
	for k := range em.sprites {
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
			Type:        em.sprites[id].Type,
			Dead:        em.sprites[id].Dead,
			Data: map[string]string{
				"Image": em.sprites[id].Image,
			},
		})
	}
	return entities
}
