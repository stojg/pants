package vector

import (
	"testing"
)

var equalsTests = []struct {
	a   *Vec2
	b   *Vec2
	out bool
}{
	{&Vec2{0, 0}, &Vec2{0, 0}, true},
	{&Vec2{1, 0}, &Vec2{1, 0}, true},
	{&Vec2{0, 1}, &Vec2{0, 1}, true},
	{&Vec2{0, 0}, &Vec2{1, 0}, false},
	{&Vec2{0, 0}, &Vec2{0, 1}, false},
	{&Vec2{0, 1}, &Vec2{1, 1}, false},
	{&Vec2{1, 0}, &Vec2{1, 1}, false},
	{&Vec2{1, 1}, &Vec2{0, 1}, false},
	{&Vec2{1, 1}, &Vec2{1, 0}, false},
}

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
	for _, tt := range equalsTests {
		s := tt.a.Equals(tt.b)
		if s != tt.out {
			t.Errorf("%v.Equals(%v) => %t, want %t", tt.a, tt.b, s, tt.out)
		}
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
