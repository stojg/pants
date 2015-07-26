package main

import "testing"

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
