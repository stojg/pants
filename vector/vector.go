package vector

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

func MapToRange(rotation float64) float64 {
	for rotation < -math.Pi {
		rotation += math.Pi * 2
	}
	for rotation > math.Pi {
		rotation -= math.Pi * 2
	}
	return rotation
}

type Vec2 struct {
	X, Y float64
}


func (a *Vec2) Direction(b *Vec2) *Vec2 {
	return a.Clone().Sub(b)
}

func (a *Vec2) Clear() {
	a.X = 0
	a.Y = 0
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

func (a *Vec2) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y)
}

func (a *Vec2) Add(o *Vec2) *Vec2 {
	a.X += o.X
	a.Y += o.Y
	return a
}

func (a *Vec2) Sub(b *Vec2) *Vec2 {
	a.X -= b.X
	a.Y -= b.Y
	return a
}

func (a *Vec2) Clone() *Vec2 {
	return &Vec2{
		X: a.X,
		Y: a.Y,
	}
}

func (a *Vec2) Copy(b *Vec2) {
	a.X = b.X
	a.Y = b.Y
}

func (a *Vec2) Normalize() *Vec2 {
	length := a.Length()
	if length > 0 {
		a.Scale(1 / length)
	}
	return a
}

func (a *Vec2) Scale(alpha float64) *Vec2 {
	a.X *= alpha
	a.Y *= alpha
	return a
}

func (v *Vec2) ToOrientation() float64 {
	return math.Atan2(v.X, -v.Y)
}
