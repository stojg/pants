package physics

import (
	"github.com/stojg/pants/structs"
	"github.com/stojg/pants/vector"
	"testing"
)

func TestCMAdd(t *testing.T) {
	cm := NewCollisionManager()
	cm.Add(NewPhysics(structs.NewSetData(0, 0, 0, 0)), 0)
	if cm.Length() != 1 {
		t.Errorf("Expected that there would be one item in the list")
	}
}

func TestCMRemove(t *testing.T) {
	cm := NewCollisionManager()
	p1 := NewPhysics(structs.NewSetData(0, 0, 0, 0))
	cm.Add(p1, 0)
	cm.Remove(1)
	if cm.Length() != 1 {
		t.Errorf("Expected that there would be one item in the list")
	}
	if cm.physics[0] != p1 {
		t.Errorf("Expected that the remaining item would be p2")
	}
}

func TestDetectCollisions(t *testing.T) {
	cm := NewCollisionManager()
	cm.Add(NewPhysics(structs.NewSetData(0, 0, 0, 0)), 0)
	cm.Add(NewPhysics(structs.NewSetData(10, 0, 0, 0)), 1)
	cm.DetectCollisions()
}

func TestResolveCollisions(t *testing.T) {
	cm := NewCollisionManager()
	cm.ResolveCollisions(0.016)
}

func TestContactPairHit(t *testing.T) {
	cm := NewCollisionManager()
	a := NewPhysics(structs.NewSetData(0, 0, 20, 20))
	b := NewPhysics(structs.NewSetData(1, 0, 20, 20))
	contact, err := cm.GenerateContacts(a, b)
	if err != nil {
		t.Errorf("GenerateContacts error '%s'", err)
		return
	}

	if contact == nil {
		t.Errorf("No collision detected")
		return
	}

	actual := contact.penetration
	expected := 19.0
	if actual != expected {
		t.Errorf("Contact.a should be %v, not %v", expected, actual)
		return
	}
}

func TestDetectNoCollisionsOutsideOfGrid(t *testing.T) {
	cm := NewCollisionManager()
	// this entity is partially outside the grid and should not be checked
	a := NewPhysics(structs.NewSetData(9, 9, 20, 20))
	b := NewPhysics(structs.NewSetData(10, 10, 20, 20))
	cm.Add(a, 0)
	cm.Add(b, 1)

	cm.DetectCollisions()

	numCollisions := len(cm.collisions)
	expectedCollisions := 0
	if numCollisions != expectedCollisions {
		t.Errorf("Expected that there would be %d collision(s), not %d", expectedCollisions, numCollisions)
	}
}

func TestCMDetectCollisions(t *testing.T) {
	cm := NewCollisionManager()

	a := NewPhysics(structs.NewSetData(10, 10, 20, 20))
	// standing still
	a.Data.Velocity = &vector.Vec2{0, 0}
	b := NewPhysics(structs.NewSetData(20, 10, 20, 20))
	// moving towards a
	b.Data.Velocity = &vector.Vec2{-1, 0}
	cm.Add(a, 0)
	cm.Add(b, 1)

	cm.DetectCollisions()

	numCollisions := len(cm.collisions)
	expectedCollisions := 1
	if numCollisions != expectedCollisions {
		t.Errorf("Expected that there would be %d collision(s), not %d", expectedCollisions, numCollisions)
		return
	}

	contact := cm.collisions[0]
	expectedPenetration := 10.0
	if contact.penetration != expectedPenetration {
		t.Errorf("Expected that the collision would have a penetration of %f, not %f", expectedPenetration, contact.penetration)
	}

	contact.resolveInterpenetration()

	expectedPenetration = 0
	if contact.penetration != expectedPenetration {
		t.Errorf("Expected that the collision would have a penetration of %f, not %f", expectedPenetration, contact.penetration)
	}

	sepVel := contact.separatingVelocity()
	expectedSepVel := -1.0
	if sepVel != expectedSepVel {
		t.Errorf("The separating velocity should be %f, not %f", expectedSepVel, sepVel)
	}

	// full bounce
	contact.restitution = 1
	contact.resolveVelocity(1.0)

	sepVel = contact.separatingVelocity()
	expectedSepVel = 1.0
	if sepVel != expectedSepVel {
		t.Errorf("The separating velocity should be %f, not %f", expectedSepVel, sepVel)
	}

	aExpectedPosition := &vector.Vec2{5, 10}
	if !a.Data.Position.Equals(aExpectedPosition) {
		t.Errorf("a.Position should be %V, not %V", aExpectedPosition, a.Data.Position)
	}

	bExpectedPosition := &vector.Vec2{25, 10}
	if !b.Data.Position.Equals(bExpectedPosition) {
		t.Errorf("b.Position should be %V, not %V", bExpectedPosition, b.Data.Position)
	}

	// since a was standing still it should take all the energy
	aExpectedVelocity := &vector.Vec2{-1, 0}
	if !a.Data.Velocity.Equals(aExpectedVelocity) {
		t.Errorf("b.Velocity should be {%f, %f}, not {%f, %f}", aExpectedVelocity.X, aExpectedVelocity.Y, a.Data.Velocity.X, a.Data.Velocity.Y)
	}

	// since a was standing still it should take all the energy and b will stop
	bExpectedVelocity := &vector.Vec2{0, 0}
	if !b.Data.Velocity.Equals(bExpectedVelocity) {
		t.Errorf("b.Velocity should be {%f, %f}, not {%f, %f}", bExpectedVelocity.X, bExpectedVelocity.Y, b.Data.Velocity.X, b.Data.Velocity.Y)
	}
}

// ~50000 28567
func BenchmarkDetectCollisions(bench *testing.B) {
	cm := NewCollisionManager()

	a := NewPhysics(structs.NewSetData(10, 10, 0, 0))
	// standing still
	a.Data.Velocity = &vector.Vec2{0, 0}

	b := NewPhysics(structs.NewSetData(20, 10, 0, 0))
	// moving towards a
	b.Data.Velocity = &vector.Vec2{-1, 0}
	cm.Add(a, 0)
	cm.Add(b, 1)

	bench.ResetTimer()
	for n := 0; n < bench.N; n++ {
		cm.DetectCollisions()
	}
}
