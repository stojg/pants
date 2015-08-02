package main

import (
	"math"
)

const EPSILON = 0.001

func RadiansVec2(radians float64) *Vec2 {
	return &Vec2{
		X: math.Sin(radians),
		Y: -math.Cos(radians),
	}
}

type Vec2 struct {
	X, Y float64
}

func (a *Vec2) Multiply(scalar float64) *Vec2 {
	return &Vec2{
		a.X * scalar,
		a.Y * scalar,
	}
}

func (a *Vec2) Equals(b *Vec2) bool {
	if a.X == b.X && a.Y == b.Y {
		return true
	}
	if math.Abs(a.X-b.X) > EPSILON {
		return false
	}
	if math.Abs(a.Y-b.Y) > EPSILON {
		return false
	}
	return true
}

func (a *Vec2) Add(o *Vec2) {
	a.X += o.X
	a.Y += o.Y
}

func (a *Vec2) Clone() *Vec2 {
	return &Vec2{
		X: a.X,
		Y: a.Y,
	}
}

func (a *Vec2) Copy(b *Vec2) {
	a.X = b.X
	a.Y = a.Y
}
