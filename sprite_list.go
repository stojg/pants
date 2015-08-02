package main
import (
	"time"
	"labix.org/v2/mgo/bson"
)

func NewSpriteList() *SpriteList {
	return  &SpriteList{
		sprites: make(map[uint64]*Sprite),
	}
}

// SpriteList is a simple struct that contains and interacts with Sprites /
// Entities.
type SpriteList struct {
	lastEntityID uint64
	sprites      map[uint64]*Sprite
}

func (s *SpriteList) NewSprite(x, y float64, image string) uint64 {
	sprite := &Sprite{}
	sprite.Dead = false
	sprite.SetPosition(x, y)
	sprite.SetVelocity(&Vec2{0, 0})
	sprite.Image = image
	sprite.inputs = make([]*InputRequest, 0)
	sprite.Orientation = 3.14 / 2
	sprite.AIs = make([]AI, 0)
	sprite.AIs = append(sprite.AIs, &AIDrunkard{})
	return s.Add(sprite)
}

func (s *SpriteList) Add(sprite *Sprite) uint64 {
	s.lastEntityID++
	sprite.Id = s.lastEntityID
	s.sprites[sprite.Id] = sprite
	return s.lastEntityID
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

func (s *SpriteList) Update(w *World, duration float64) {
	seconds :=  duration
	for _, sprite := range s.sprites {
		if sprite.Dead {
			continue
		}
		// handle AI
		for _, ai := range sprite.AIs {
			ai.Update(sprite, w, duration)
		}
		// clear inputs
		sprite.inputs = sprite.inputs[0:0]
		// shitty integrate
		if sprite.velocity != nil {
			sprite.X += sprite.velocity.X * seconds
			sprite.Y += sprite.velocity.Y * seconds
			sprite.changed = true
		}

	}
}

func (s *SpriteList) Changed(reset bool) []*Sprite {
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

