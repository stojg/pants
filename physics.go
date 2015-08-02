package main
import (
)

func NewPhysicsComponent(x, y, orientation float64) *PhysicsComponent {
	return &PhysicsComponent{
		Position:     &Vec2{X: x, Y: y},
		prevPosition: &Vec2{X: x, Y: y},
		Orientation:  orientation,
		Velocity:     &Vec2{0, 0},
	}
}

type PhysicsComponent struct {
	Position     *Vec2
	prevPosition *Vec2
	Orientation  float64
	Velocity     *Vec2
	Damping      float64
}

func (c *PhysicsComponent) Update(sprite *Sprite, w *World, seconds float64) {
	c.Position.Add(c.Velocity.Multiply(seconds))
	if !c.Position.Equals(c.prevPosition) {
		c.prevPosition.Copy(c.Position)
		list.updated[sprite.Id] = true
	}
}
