package ai

import (
	. "github.com/stojg/pants/physics"
	"testing"
)

func TestArrive(t *testing.T) {
	arrive := NewArrive(NewPhysics(0, 0, 0, 1, 0, 0), NewPhysics(1, 0, 0, 1, 0, 0), 10, 20)
	actual := arrive.Get().Linear.Length()
	expected := 0.0
	if actual != expected {
		t.Errorf("Arrive.Get expected %f, got %f", expected, actual)
	}
}

func TestArriveStillToClose(t *testing.T) {
	arrive := NewArrive(NewPhysics(0, 0, 0, 1, 0, 0), NewPhysics(10, 0, 0, 1, 0, 0), 10, 20)
	actual := arrive.Get().Linear.Length()
	expected := 0.0
	if actual != expected {
		t.Errorf("Arrive.Get expected %f, got %f", expected, actual)
	}
}

func TestArriveWillMove(t *testing.T) {
	entity := NewPhysics(0, 0, 0, 1, 0, 0)
	arrive := NewArrive(entity, NewPhysics(12, 0, 0, 1, 0, 0), 10, 20)
	actual := arrive.Get().Linear.Length()
	expected := entity.Data.MaxAcceleration
	if actual != expected {
		t.Errorf("Arrive.Get expected %f, got %f", expected, actual)
	}
}
