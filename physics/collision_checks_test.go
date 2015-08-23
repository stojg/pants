package physics

import (
	"testing"
)

var circleVsCircleTests = []struct {
	a   *Circle
	b   *Circle
	hit bool
	pen float64
}{
	{NewCircle(0, 0, 10), NewCircle(0, 0, 10), true, 20},
	{NewCircle(0, 0, 10), NewCircle(0, 1, 10), true, 19},
	{NewCircle(0, 0, 1), NewCircle(0, 10, 1), false, 0},
	{NewCircle(0, 0, 1), NewCircle(0, 3, 1), false, 0},
	{NewCircle(0, 0, 5), NewCircle(0, 9, 5), true, 1},
	{NewCircle(0, -5, 5), NewCircle(0, 4, 5), true, 1},
	{NewCircle(0, -5, 5), NewCircle(0, 6, 5), false, 0},
}

func TestCircleVsCircle(t *testing.T) {
	for _, td := range circleVsCircleTests {
		contact := CircleVsCircle(td.a, td.b)
		if (contact == nil) == td.hit {
			t.Errorf("Expected %#v vs %#v to be %t", td.a, td.b, td.hit)
		}
		if contact != nil && contact.penetration != td.pen {
			t.Errorf("Expected penetration to be %f, got %f", td.pen, contact.penetration)
		}
	}
}

func BenchmarkCircleVsCircleMiss(b *testing.B) {
	c1 := NewCircle(0, 0, 1)
	c2 := NewCircle(0, 3, 1)
	for n := 0; n < b.N; n++ {
		CircleVsCircle(c1, c2)
	}
}

func BenchmarkCircleVsCircleHit(b *testing.B) {
	c1 := NewCircle(0, 0, 2)
	c2 := NewCircle(0, 3, 2)
	for n := 0; n < b.N; n++ {
		CircleVsCircle(c1, c2)
	}
}
