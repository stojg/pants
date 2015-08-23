package ai

import (
	. "github.com/stojg/pants/physics"
	"testing"
)

func TestArrive(t *testing.T) {
	arrive := NewArrive(NewPhysics(0, 0, 0, 1), NewPhysics(1, 0, 0, 1), 10, 20)
	actual := arrive.Get(0.016).Linear.Length()
	expected := 0.0
	if actual != expected {
		t.Errorf("Arrive.Get expected %f, got %f", expected, actual)
	}
}

func TestArriveStillToClose(t *testing.T) {
	arrive := NewArrive(NewPhysics(0, 0, 0, 1), NewPhysics(10, 0, 0, 1), 10, 20)
	actual := arrive.Get(0.016).Linear.Length()
	expected := 0.0
	if actual != expected {
		t.Errorf("Arrive.Get expected %f, got %f", expected, actual)
	}
}

func TestArriveWillMove(t *testing.T) {
	entity := NewPhysics(0, 0, 0, 1)
	arrive := NewArrive(entity, NewPhysics(12, 0, 0, 1), 10, 20)
	actual := arrive.Get(0.016).Linear.Length()
	expected := entity.MaxAcceleration()
	if actual != expected {
		t.Errorf("Arrive.Get expected %f, got %f", expected, actual)
	}
}