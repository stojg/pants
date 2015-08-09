package main

import (
	. "github.com/stojg/pants/physics"
	. "github.com/stojg/pants/vector"
	"github.com/stojg/pants/collision"
)

func NewEntityManager() *EntityManager {
	return &EntityManager{
		sprites:       make(map[uint64]*Sprite),
		ais:           make(map[uint64]AI),
		physics:       make(map[uint64]*Physics),
		updated:       make(map[uint64]bool),
		forceRegistry: &PhysicsForceRegistry{},
		collisionManager: &collision.CollisionManager{},
	}
}

// SpriteList is a simple struct that contains and interacts with Sprites /
// Entities.
type EntityManager struct {
	lastEntityID  uint64
	sprites       map[uint64]*Sprite
	ais           map[uint64]AI
	physics       map[uint64]*Physics
	updated       map[uint64]bool
	forceRegistry *PhysicsForceRegistry
	collisionManager *collision.CollisionManager
}

var gravity *StaticForce

func init() {
	gravity = &StaticForce{}
	gravity.SetForce(&Vec2{0, 10})
}

func (s *EntityManager) NewEntity(x, y float64, image string) uint64 {
	sprite := &Sprite{}
	sprite.Dead = false
	sprite.Type = "sprite"
	sprite.Image = image
	sprite.inputs = make([]*InputRequest, 0)
	s.lastEntityID++
	sprite.Id = s.lastEntityID
	s.sprites[sprite.Id] = sprite
	s.ais[sprite.Id] = &AIDrunkard{state: NewStateMachine(STATE_IDLE)}
	s.physics[sprite.Id] = NewPhysics(x, y, 3.14*2, 1)
	s.physics[sprite.Id].SetDamping(0.99)
	s.forceRegistry.Add(s.physics[sprite.Id], gravity)
	s.updated[sprite.Id] = true
	s.collisionManager.Add(s.physics[sprite.Id])
	return sprite.Id
}

func (s *EntityManager) all() []*EntityUpdate {
	entities := make([]*EntityUpdate, 0, len(s.sprites))
	for id, spr := range s.sprites {

		data := make(map[string]string)
		data["Image"] = spr.Image
		entities = append(entities, &EntityUpdate{
			Id:          id,
			Type:        "sprite",
			X:           s.physics[id].Position.X,
			Y:           s.physics[id].Position.Y,
			Orientation: s.physics[id].Orientation,
			Data:        data,
			Dead:        s.sprites[id].Dead,
		})
	}
	return entities
}

func (s *EntityManager) Update(w *World, duration float64) {

	s.forceRegistry.Update(duration)

	for _, sprite := range s.sprites {
		s.ais[sprite.Id].Update(sprite, w, duration)
	}

	for _, sprite := range s.sprites {
		moved := s.physics[sprite.Id].Update(sprite.Id, w, duration)
		if moved {
			s.updated[sprite.Id] = true
		}
	}

	s.collisionManager.DetectCollisions(duration)
	s.collisionManager.ResolveCollisions(duration)
}
