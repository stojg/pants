package main

import (
	. "github.com/stojg/pants/vector"
)

func NewEntityList() *EntityList {
	return &EntityList{
		sprites:       make(map[uint64]*Sprite),
		ais:           make(map[uint64]AI),
		physics:       make(map[uint64]*Physics),
		updated:       make(map[uint64]bool),
		forceRegistry: &PhysicsForceRegistry{},
	}
}

// SpriteList is a simple struct that contains and interacts with Sprites /
// Entities.
type EntityList struct {
	lastEntityID  uint64
	sprites       map[uint64]*Sprite
	ais           map[uint64]AI
	physics       map[uint64]*Physics
	updated       map[uint64]bool
	forceRegistry *PhysicsForceRegistry
}

var gravity *ParticleGravity

func init() {
	gravity = &ParticleGravity{
		gravity: &Vec2{0,10},
	}
}

func (s *EntityList) NewEntity(x, y float64, image string) uint64 {
	sprite := &Sprite{}
	sprite.Dead = false
	sprite.Type = "sprite"
	sprite.Image = image
	sprite.inputs = make([]*InputRequest, 0)
	s.lastEntityID++
	sprite.Id = s.lastEntityID
	s.sprites[sprite.Id] = sprite
	s.ais[sprite.Id] = &AIDrunkard{state: NewStateMachine(STATE_IDLE)}
	s.physics[sprite.Id] = NewPhysicsComponent(x, y, 3.14*2)
	s.physics[sprite.Id].setMass(1)
	s.physics[sprite.Id].setDamping(0.99)
	s.forceRegistry.Add(s.physics[sprite.Id], gravity)
	s.updated[sprite.Id] = true
	return sprite.Id
}

func (s *EntityList) Add(e *Sprite) {
	s.lastEntityID++
	e.Id = s.lastEntityID
	s.sprites[s.lastEntityID] = e
}

func (s *EntityList) all() []*EntityUpdate {
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

func (s *EntityList) Update(w *World, duration, gameTime float64) {

	s.forceRegistry.Update(duration)

	for _, sprite := range s.sprites {
		s.ais[sprite.Id].Update(sprite, w, duration)
	}
	for _, sprite := range s.sprites {
		s.physics[sprite.Id].Update(sprite, w, duration)
	}
}
