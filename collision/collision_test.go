package collision

import (
	"testing"
)

var circleVsCircleTests = []struct {
	a   *Circle
	b   *Circle
	hit bool
}{
	{NewCircle(0, 0, 10), NewCircle(0, 0, 10), true},
	{NewCircle(0, 0, 10), NewCircle(0, 1, 10), true},
	{NewCircle(0, 0, 1), NewCircle(0, 10, 1), false},
	{NewCircle(0, 0, 1), NewCircle(0, 3, 1), false},
	{NewCircle(0, 0, 5), NewCircle(0, 10, 5), true},
	{NewCircle(0, -5, 5), NewCircle(0, 5, 5), true},
	{NewCircle(0, -5, 5), NewCircle(0, 6, 5), false},
}

func TestCircleVsCircle(t *testing.T) {
	for _, td := range circleVsCircleTests {
		contact := CircleVsCircle(td.a, td.b)
		if contact.hit != td.hit {
			t.Errorf("Expected %#v vs %#v to be %t, got %t", td.a, td.b, td.hit, contact.hit)
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
