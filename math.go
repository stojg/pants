package main

import "math"

type Vec2 struct {
	X, Y float64
}

func (v *Vec2) Multiply(scalar float64) *Vec2 {
	return &Vec2{
		v.X * scalar,
		v.Y * scalar,
	}
}

func RadiansVec2(radians float64) *Vec2 {
	return &Vec2{
		X: math.Sin(radians),
		Y: -math.Cos(radians),
	}
}
