package main

import (
	"labix.org/v2/mgo/bson"
	"time"
)

func init() {
	list = &SpriteList{
		sprites: make(map[uint64]*Sprite),
	}
}

type Sprite struct {
	Id       uint64 `bson:",minsize"`
	AIs      []AI
	X, Y     float64 `bson:",minsize,omitempty"`
	velocity struct {
		x, y float64
	}
	Rotation float64 `bson:",minsize"`
	Image    string  `bson:",minsize,omitempty"`
	inputs   []*InputRequest
	changed  bool
	Dead     bool
}

func (s *Sprite) SetVelocity(x, y float64) {
	s.velocity.x = x
	s.velocity.y = y
	s.changed = true
}

func (s *Sprite) SetRotation(orientation float64) {
	s.Rotation = orientation
	s.changed = true
}

func (s *Sprite) Kill() {
	s.Dead = true
	s.changed = true
}

func (s *Sprite) Velocity() (x, y float64) {
	return s.velocity.x, s.velocity.y
}

func (s *Sprite) SetPosition(x, y float64) {
	s.X, s.Y = x, y
	s.changed = true
}

func (s *Sprite) Update(w *World, t time.Duration) {
	if s.Dead {
		return
	}
	s.handleAI(w, t)
	s.inputs = s.inputs[0:0]
	s.integrate(t)
}

func (s *Sprite) integrate(t time.Duration) {
	moveX := s.velocity.x * t.Seconds()
	moveY := s.velocity.y * t.Seconds()
	s.SetPosition(s.X+moveX, s.Y+moveY)
}

func (s *Sprite) AddInput(i *InputRequest) {
	s.inputs = append(s.inputs, i)
}

func (s *Sprite) handleAI(w *World, t time.Duration) {
	for _, ai := range s.AIs {
		ai.Update(s, w, t)
	}
}

// SpriteList is a simple struct that contains and interacts with Sprites /
// Entities.
type SpriteList struct {
	lastEntityID uint64
	sprites map[uint64]*Sprite
}

func (s *SpriteList) NewSprite(x, y float64, image string) {
	sprite := &Sprite{}
	sprite.Dead = false
	sprite.SetPosition(x, y)
	sprite.Image = image
	sprite.inputs = make([]*InputRequest, 0)
	sprite.Rotation = 3.14 / 2
	sprite.AIs = make([]AI, 0)
	sprite.AIs = append(sprite.AIs, &AIDrunkard{})
	s.Add(sprite)
}

func (s *SpriteList) Add(sprite *Sprite) {
	s.lastEntityID++
	sprite.Id = s.lastEntityID
	s.sprites[sprite.Id] = sprite
}

func (s *SpriteList) SendAll(c *connection) {
	t := &Message{
		Topic:     "all",
		Data:      s.All(),
		Timestamp: float64(time.Now().UnixNano()) / 1000000,
	}
	msg, _ := bson.Marshal(t)
	c.send <- msg
}

func (s *SpriteList) All() []*Sprite {
	v := make([]*Sprite, 0, len(s.sprites))
	for _, value := range s.sprites {
		v = append(v, value)
	}
	return v
}

func (s *SpriteList) Changed(reset bool) []*Sprite {
	toSend := make([]*Sprite, 0)
	for _, spr := range s.sprites {
		if spr.changed {
			toSend = append(toSend, spr)
			if reset {
				spr.changed = false
			}
		}
	}
	return toSend
}
