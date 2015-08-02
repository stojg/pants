package main

import ()

func init() {
	list = NewSpriteList()
}

type Sprite struct {
	Id          uint64 `bson:",minsize"`
	State       string
	AIs         []AI
	X, Y        float64 `bson:",minsize,omitempty"`
	Orientation float64 `bson:",minsize"`
	Image       string  `bson:",minsize,omitempty"`
	Dead        bool
	velocity    *Vec2
	inputs      []*InputRequest
	changed     bool
}

func (s *Sprite) SetVelocity(vec *Vec2) {
	s.velocity = vec
	s.changed = true
}

func (s *Sprite) SetOrientation(orientation float64) {
	s.Orientation = orientation
	s.changed = true
}

func (s *Sprite) Kill() {
	s.Dead = true
	s.changed = true
}

func (s *Sprite) Velocity() *Vec2 {
	return s.velocity
}

func (s *Sprite) SetPosition(x, y float64) {
	s.X, s.Y = x, y
	s.changed = true
}

func (s *Sprite) AddInput(i *InputRequest) {
	s.inputs = append(s.inputs, i)
}
