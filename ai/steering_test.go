package ai

import (
	. "github.com/stojg/pants/physics"
	"github.com/stojg/pants/structs"
	"testing"
)

func TestArrive(t *testing.T) {
	arrive := NewArrive(NewPhysics(structs.NewSetData(0, 0, 0, 0)), NewPhysics(structs.NewSetData(1, 0, 0, 0)), 10, 20)
	actual := arrive.Get().Linear.Length()
	expected := 0.0
	if actual != expected {
		t.Errorf("Arrive.Get expected %f, got %f", expected, actual)
	}
}

func TestArriveStillToClose(t *testing.T) {
	arrive := NewArrive(NewPhysics(structs.NewSetData(0, 0, 0, 0)), NewPhysics(structs.NewSetData(10, 0, 0, 0)), 10, 20)
	actual := arrive.Get().Linear.Length()
	expected := 0.0
	if actual != expected {
		t.Errorf("Arrive.Get expected %f, got %f", expected, actual)
	}
}

func TestArriveWillMove(t *testing.T) {
	entity := NewPhysics(structs.NewSetData(0, 0, 0, 0))
	arrive := NewArrive(entity, NewPhysics(structs.NewSetData(12, 0, 0, 0)), 10, 20)
	actual := arrive.Get().Linear.Length()
	expected := entity.Data.MaxAcceleration
	if actual != expected {
		t.Errorf("Arrive.Get expected %f, got %f", expected, actual)
	}
}
