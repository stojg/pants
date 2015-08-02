package main

import (
	"labix.org/v2/mgo/bson"
	"time"
)

func NewEntityList() *EntityList {
	return &EntityList{
		sprites:           make(map[uint64]*Sprite),
		AIComponents:      make(map[uint64]AI),
		PhysicsComponents: make(map[uint64]*PhysicsComponent),
	}
}

// SpriteList is a simple struct that contains and interacts with Sprites /
// Entities.
type EntityList struct {
	lastEntityID      uint64
	sprites           map[uint64]*Sprite
	AIComponents      map[uint64]AI
	PhysicsComponents map[uint64]*PhysicsComponent
}

func (s *EntityList) NewEntity(x, y float64, image string) uint64 {
	sprite := &Sprite{}
	sprite.Dead = false
	sprite.Image = image
	sprite.inputs = make([]*InputRequest, 0)
	s.lastEntityID++
	sprite.Id = s.lastEntityID
	s.sprites[sprite.Id] = sprite
	s.AIComponents[sprite.Id] = &AIDrunkard{state: "idle"}
	s.PhysicsComponents[sprite.Id] = NewPhysicsComponent(x, y, 3.14/2)
	return s.lastEntityID
}

func (s *EntityList) SendAll(c *connection) {
	t := &Message{
		Topic:     "all",
		Data:      s.All(),
		Timestamp: float64(time.Now().UnixNano()) / 1000000,
	}
	msg, _ := bson.Marshal(t)
	c.send <- msg
}

func (s *EntityList) All() []*EntityUpdate {
	v := make([]*EntityUpdate, 0, len(s.sprites))
	for id, spr := range s.sprites {
		v = append(v, &EntityUpdate{
			Id:          id,
			X:           s.PhysicsComponents[id].Position.X,
			Y:           s.PhysicsComponents[id].Position.Y,
			Orientation: s.PhysicsComponents[id].Orientation,
			Image:       spr.Image,
		})
	}
	return v
}

func (s *EntityList) Update(w *World, duration, gameTime float64) {
	for _, sprite := range s.sprites {
		s.AIComponents[sprite.Id].Update(sprite, w, duration)
	}
	for _, sprite := range s.sprites {
		s.PhysicsComponents[sprite.Id].Update(sprite, w, duration)
	}
}

func (s *EntityList) Changed(reset bool) []*Sprite {
	toSend := make([]*Sprite, 0)
	for _, spr := range s.sprites {
		if !spr.changed {
			continue
		}
		toSend = append(toSend, spr)
		if reset {
			spr.changed = false
		}
	}
	return toSend
}
