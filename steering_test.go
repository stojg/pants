package main

import (
	. "github.com/stojg/pants/vector"
	"testing"
)

func TestSeek(t *testing.T) {
	arrive := &Arrive{}
	arrive.source = &PhysicsComponent{
		Position: &Vec2{
			X: 0,
			Y: 0,
		},
	}
	arrive.target = &PhysicsComponent{
		Position: &Vec2{
			X: 0,
			Y: 0,
		},
	}
	arrive.targetRadius = 10
	arrive.timeToTarget = 0.1
	arrive.slowRadius = 20

	s := arrive.Get(0.016)
	if s.linear.Length() != 0 {
		t.Errorf("Arrive.Get should return zero length")
	}
}

func TestSeekStillToClose(t *testing.T) {
	arrive := &Arrive{}
	arrive.source = &PhysicsComponent{
		Position: &Vec2{
			X: 0,
			Y: 0,
		},
		Velocity: &Vec2{},
	}
	arrive.target = &PhysicsComponent{
		Position: &Vec2{
			X: 11,
			Y: 0,
		},
	}
	arrive.targetRadius = 10
	arrive.timeToTarget = 1
	arrive.slowRadius = 20

	s := arrive.Get(0.016)
	if s.linear.Length() != 0 {
		t.Errorf("Arrive.Get should return zero length")
	}
}

func TestSeekWillMove(t *testing.T) {
	arrive := &Arrive{}
	arrive.source = &PhysicsComponent{
		Position: &Vec2{
			X: 0,
			Y: 0,
		},
		Velocity: &Vec2{},
		maxVelocity: 10,
	}
	arrive.target = &PhysicsComponent{
		Position: &Vec2{
			X: 12,
			Y: 0,
		},
	}
	arrive.targetRadius = 10
	arrive.timeToTarget = 1
	arrive.slowRadius = 20

	s := arrive.Get(0.016)
	if s.linear.Length() == 0 {
		t.Errorf("Arrive.Get should not zero length %f", s.linear.Length())
	}
}
