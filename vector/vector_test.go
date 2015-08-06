package vector

import (
	"math"
	"testing"
)

func TestRadiansToVec2(t *testing.T) {
	actual := RadiansVec2(0)
	expected := Vec2{X: 0, Y: -1}
	if actual.X != expected.X {
		t.Errorf("expected %f, got %f", expected.X, actual.X)
	}
	if actual.Y != expected.Y {
		t.Errorf("expected %f, got %f", expected.Y, actual.Y)
	}
}

func TestEquals(t *testing.T) {
	a := &Vec2{0, 0}
	b := &Vec2{0, 0}
	if !a.Equals(b) {
		t.Errorf("Equals not working")
	}
	a.X = 1
	b.X = 1
	if !a.Equals(b) {
		t.Errorf("Equals not working")
	}

	b.X = 1
	a.X = 0
	if a.Equals(b) {
		t.Errorf("Equals not working a.X %f = b.X %f diffX %f", a.X, b.X, math.Abs(a.X-b.X))
		t.Errorf("Equals not working a.Y %f = b.Y %f diffY %f", a.Y, b.Y, math.Abs(a.Y-b.Y))
	}

	b.X = 0
	a.X = 1
	if a.Equals(b) {
		t.Errorf("Equals not working %f = %f %f", a.X, b.X, math.Abs(a.X-b.X))
	}

	a.X = 1
	a.Y = 0
	b.X = 1
	b.Y = 1
	if a.Equals(b) {
		t.Errorf("Equals not working a.X %f = b.X %f diffX %f", a.X, b.X, math.Abs(a.X-b.X))
		t.Errorf("Equals not working a.Y %f = b.Y %f diffY %f", a.Y, b.Y, math.Abs(a.Y-b.Y))
	}

	a.Y = 1
	b.Y = 0
	if a.Equals(b) {
		t.Errorf("Equals not  working %f", math.Abs(a.Y-b.Y))
	}
}

func BenchmarkRadiansToVec2(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		RadiansVec2(0)
	}
}

func TestVec2Multiply(t *testing.T) {
	vec := &Vec2{X: 1, Y: 0}
	expected := &Vec2{X: 10, Y: 0}
	actual := vec.Multiply(10)
	if actual.X != expected.X {
		t.Errorf("expected %f, got %f", expected.X, actual.X)
	}
	if actual.Y != expected.Y {
		t.Errorf("expected %f, got %f", expected.Y, actual.Y)
	}
}

func BenchmarkVec2Multiply(b *testing.B) {
	vec := &Vec2{X: 1, Y: 0}
	for n := 0; n < b.N; n++ {
		vec.Multiply(10)
	}

}
