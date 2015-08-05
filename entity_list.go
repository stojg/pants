package main

import (
	"labix.org/v2/mgo/bson"
	"time"
)

func NewEntityList() *EntityList {
	return &EntityList{
		sprites: make(map[uint64]*Sprite),
		ais:     make(map[uint64]AI),
		physics: make(map[uint64]*PhysicsComponent),
		updated: make(map[uint64]bool),
	}
}

// SpriteList is a simple struct that contains and interacts with Sprites /
// Entities.
type EntityList struct {
	lastEntityID uint64
	sprites      map[uint64]*Sprite
	ais          map[uint64]AI
	physics      map[uint64]*PhysicsComponent
	updated      map[uint64]bool
}

func (s *EntityList) NewEntity(x, y float64, image string) uint64 {
	sprite := &Sprite{}
	sprite.Dead = false
	sprite.Image = image
	sprite.inputs = make([]*InputRequest, 0)
	s.lastEntityID++
	sprite.Id = s.lastEntityID
	s.sprites[sprite.Id] = sprite
	s.ais[sprite.Id] = &AIDrunkard{state: "idle"}
	s.physics[sprite.Id] = NewPhysicsComponent(x, y, 3.14*2)
	s.physics[sprite.Id].setMass(10)
	s.physics[sprite.Id].setDamping(0.99)
	s.physics[sprite.Id].setAcceleration(&Vec2{0,10})
	s.updated[sprite.Id] = true
	return sprite.Id
}

func (s *EntityList) SendAll(c *connection) {
	t := &Message{
		Topic:     "all",
		Data:      s.all(),
		Timestamp: float64(time.Now().UnixNano()) / 1000000,
	}
	msg, _ := bson.Marshal(t)
	c.send <- msg
}

func (s *EntityList) all() []*EntityUpdate {
	entities := make([]*EntityUpdate, 0, len(s.sprites))
	for id, spr := range s.sprites {
		entities = append(entities, &EntityUpdate{
			Id:          id,
			X:           s.physics[id].Position.X,
			Y:           s.physics[id].Position.Y,
			Orientation: s.physics[id].Orientation,
			Image:       spr.Image,
		})
	}
	return entities
}

func (s *EntityList) Update(w *World, duration, gameTime float64) {
	for _, sprite := range s.sprites {
		s.ais[sprite.Id].Update(sprite, w, duration)
	}
	for _, sprite := range s.sprites {
		s.physics[sprite.Id].Update(sprite, w, duration)
	}
}
