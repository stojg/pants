package vector

import (
	"math"
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
	{&Vec2{EPSILON, 0}, &Vec2{0, 0}, true},
	{&Vec2{0, EPSILON}, &Vec2{0, 0}, true},
	{&Vec2{0, 0}, &Vec2{1, 0}, false},
	{&Vec2{0, 0}, &Vec2{0, 1}, false},
	{&Vec2{0, 1}, &Vec2{1, 1}, false},
	{&Vec2{1, 0}, &Vec2{1, 1}, false},
	{&Vec2{1, 1}, &Vec2{0, 1}, false},
	{&Vec2{1, 1}, &Vec2{1, 0}, false},
}

func TestEquals(t *testing.T) {
	for _, tt := range equalsTests {
		actual := tt.a.Equals(tt.b)
		if actual != tt.out {
			t.Errorf("%v.Equals(%v) => %t, want %t", tt.a, tt.b, actual, tt.out)
		}
	}
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

var mapToRangeTests = []struct {
	in  float64
	out float64
}{
	{0, 0},
	{2 * math.Pi, 0},
	{-2 * math.Pi, 0},
	{math.Pi, math.Pi},
	{3 * math.Pi, math.Pi},
	{4 * math.Pi, 0},
}

func TestMapToRange(t *testing.T) {
	for _, tt := range mapToRangeTests {
		actual := MapToRange(tt.in)
		if actual != tt.out {
			t.Errorf("MapToRange(%f) => %f, want %f", tt.in, actual, tt.out)
		}
	}
}

var normalizeTests = []struct {
	in  *Vec2
	out *Vec2
}{
	{&Vec2{}, &Vec2{}},
	{&Vec2{1, 0}, &Vec2{1, 0}},
	{&Vec2{2, 0}, &Vec2{1, 0}},
	{&Vec2{-1, 0}, &Vec2{-1, 0}},
	{&Vec2{0, 2}, &Vec2{0, 1}},
}

func TestNormalize(t *testing.T) {
	for _, td := range normalizeTests {
		actual := td.in.Normalize()
		if actual.Equals(td.out) != true {
			t.Errorf("%v.Normalize() => %f, want %f", td.in, actual, td.out)
		}
	}
}

var addTests = []struct {
	in     *Vec2
	add    *Vec2
	result *Vec2
}{
	{&Vec2{0, 0}, &Vec2{0, 0}, &Vec2{0, 0}},
	{&Vec2{1, 0}, &Vec2{1, 0}, &Vec2{2, 0}},
	{&Vec2{2, 0}, &Vec2{1, 0}, &Vec2{3, 0}},
	{&Vec2{-1, 0}, &Vec2{-1, 0}, &Vec2{-2, 0}},
	{&Vec2{0, 2}, &Vec2{0, 1}, &Vec2{0, 3}},
	{&Vec2{-5, 5}, &Vec2{2, 3}, &Vec2{-3, 8}},
}

func TestAdd(t *testing.T) {
	for _, td := range addTests {
		td.in.Add(td.add)
		if td.in.Equals(td.result) != true {
			t.Errorf("vec2.Add(%v) => %f, want %f", td.add, td.in, td.result)
		}
	}
}

var subTests = []struct {
	in     *Vec2
	sub    *Vec2
	result *Vec2
}{
	{&Vec2{0, 0}, &Vec2{0, 0}, &Vec2{0, 0}},
	{&Vec2{1, 0}, &Vec2{1, 0}, &Vec2{0, 0}},
	{&Vec2{2, 0}, &Vec2{1, 0}, &Vec2{1, 0}},
	{&Vec2{-1, 0}, &Vec2{-1, 0}, &Vec2{0, 0}},
	{&Vec2{0, 2}, &Vec2{0, 1}, &Vec2{0, 1}},
	{&Vec2{-5, 5}, &Vec2{2, 3}, &Vec2{-7, 2}},
}

func TestSub(t *testing.T) {
	for _, td := range subTests {
		td.in.Sub(td.sub)
		if td.in.Equals(td.result) != true {
			t.Errorf("vec2.Sub(%v) => %f, want %f", td.sub, td.in, td.result)
		}
	}
}

var scaleTests = []struct {
	in     *Vec2
	alpha  float64
	result *Vec2
}{
	{&Vec2{0, 0}, 1, &Vec2{0, 0}},
	{&Vec2{1, 1}, 1, &Vec2{1, 1}},
	{&Vec2{1, 1}, 0.5, &Vec2{0.5, 0.5}},
	{&Vec2{1, 0}, 2, &Vec2{2, 0}},
	{&Vec2{1, 1}, 2, &Vec2{2, 2}},
	{&Vec2{2, 1}, 2, &Vec2{4, 2}},
	{&Vec2{1, 2}, 2, &Vec2{2, 4}},
}

func TestScale(t *testing.T) {
	for _, td := range scaleTests {
		td.in.Scale(td.alpha)
		if td.in.Equals(td.result) != true {
			t.Errorf("vec2.Scale(%v) => %f, want %f", td.alpha, td.in, td.result)
		}
	}
}

func TestCopy(t *testing.T) {
	original := &Vec2{1, 2}
	copy := &Vec2{}
	original.Copy(copy)
	if !copy.Equals(original) {
		t.Errorf("vec2.Copy() copy %v != original %v ", copy, original)
	}
	original.X = 5
	if copy.Equals(original) {
		t.Errorf("vec2.Copy() copy should not change when original changes. Copy %v == original %v ", copy, original)
	}
}

func TestClear(t *testing.T) {
	original := &Vec2{1, 2}
	original.Clear()
	expected := &Vec2{}
	if !original.Equals(expected) {
		t.Errorf("vec2.Clear() should have reset vec2 %v ", original)
	}
}

func TestClone(t *testing.T) {
	original := &Vec2{1, 2}
	copy := original.Clone()
	if !copy.Equals(original) {
		t.Errorf("vec2.Copy() copy %v != original %v ", copy, original)
	}
	original.X = 5
	if copy.Equals(original) {
		t.Errorf("vec2.Copy() copy should not change when original changes. Copy %v == original %v ", copy, original)
	}
}

var lengthTests = []struct {
	in  *Vec2
	out float64
}{
	{&Vec2{}, 0},
	{&Vec2{1, 0}, 1},
	{&Vec2{-1, 0}, 1},
	{&Vec2{0, 1}, 1},
	{&Vec2{0, -1}, 1},
	{&Vec2{3, 0}, 3},
	{&Vec2{0, 3}, 3},
	{&Vec2{1, 1}, math.Sqrt(1 + 1)},
	{&Vec2{-1, -1}, math.Sqrt(1 + 1)},
}

func TestLength(t *testing.T) {
	for _, td := range lengthTests {
		actual := td.in.Length()
		if actual != td.out {
			t.Errorf("%v.Length() => %f, want %f", td.in, actual, td.out)
		}
	}
}

var squareLengthTests = []struct {
	in  *Vec2
	out float64
}{
	{&Vec2{}, 0},
	{&Vec2{1, 0}, 1},
	{&Vec2{-1, 0}, 1},
	{&Vec2{0, 1}, 1},
	{&Vec2{0, -1}, 1},
	{&Vec2{3, 0}, 9},
	{&Vec2{0, 3}, 9},
	{&Vec2{1, 1}, (1 + 1)},
	{&Vec2{-1, -1}, (1 + 1)},
}

func TestSquareLength(t *testing.T) {
	for _, td := range squareLengthTests {
		actual := td.in.SquareLength()
		if actual != td.out {
			t.Errorf("%v.Length() => %f, want %f", td.in, actual, td.out)
		}
	}
}

var toOrientationTests = []struct {
	in  *Vec2
	out float64
}{
	{&Vec2{}, math.Pi},
	{&Vec2{1, 0}, math.Pi / 2},
	{&Vec2{-1, 0}, -math.Pi / 2},
	{&Vec2{0, 1}, math.Pi},
	{&Vec2{0, -1}, 0},
	{&Vec2{3, 0}, math.Pi / 2},
	{&Vec2{0, 3}, math.Pi},
}

func TestToOrientation(t *testing.T) {
	for _, td := range toOrientationTests {
		actual := td.in.ToOrientation()
		if actual != td.out {
			t.Errorf("%v.ToOrientation() => %f, want %f", td.in, actual, td.out)
		}
	}
}

func BenchmarkRadiansToVec2(b *testing.B) {
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
