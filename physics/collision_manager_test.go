package physics

import (
	"github.com/stojg/pants/vector"
	"testing"
)

func TestCMAdd(t *testing.T) {
	cm := &CollisionManager{}
	cm.Add(NewPhysics(0, 0, 0, 0))
	if cm.Length() != 1 {
		t.Errorf("Expected that there would be one item in the list")
	}
}

func TestCMRemove(t *testing.T) {
	cm := &CollisionManager{}
	p1 := NewPhysics(0, 0, 0, 0)
	p2 := NewPhysics(0, 0, 0, 0)
	cm.Add(p1)
	cm.Remove(p2)
	if cm.Length() != 1 {
		t.Errorf("Expected that there would be one item in the list")
	}
	if cm.physics[0] != p1 {
		t.Errorf("Expected that the remaining item would be p2")
	}
}

func TestDetectCollisions(t *testing.T) {
	cm := &CollisionManager{}
	cm.Add(NewPhysics(0, 0, 0, 0))
	cm.Add(NewPhysics(10, 0, 0, 0))
	cm.DetectCollisions()
}

func TestResolveCollisions(t *testing.T) {
	cm := &CollisionManager{}
	cm.ResolveCollisions(0.016)
}

func TestContactPairHit(t *testing.T) {
	cm := &CollisionManager{}
	a := NewPhysics(0, 0, 0, 0)
	b := NewPhysics(1, 0, 0, 0)
	contact, err := cm.GenerateContacts(a, b)
	if err != nil {
		t.Errorf("GenerateContacts error '%s'", err)
		return
	}
	actual := contact.penetration
	expected := 19.0
	if actual != expected {
		t.Errorf("Contact.a should be %v, not %v", expected, actual)
	}
}

func TestDetectNoCollisionsOutsideOfGrid(t *testing.T) {
	cm := &CollisionManager{}
	// this entity is partially outside the grid and should not be checked
	a := NewPhysics(9, 9, 0, 0)
	b := NewPhysics(10, 10, 0, 0)
	cm.Add(a)
	cm.Add(b)

	cm.DetectCollisions()

	numCollisions := len(cm.collisions)
	expectedCollisions := 0
	if numCollisions != expectedCollisions {
		t.Errorf("Expected that there would be %d collision(s), not %d", expectedCollisions, numCollisions)
	}
}

func TestCMDetectCollisions(t *testing.T) {
	cm := &CollisionManager{}
	a := NewPhysics(10, 10, 0, 1)
	// standing still
	a.Velocity = &vector.Vec2{0, 0}
	b := NewPhysics(20, 10, 0, 1)
	// moving towards a
	b.Velocity = &vector.Vec2{-1, 0}
	cm.Add(a)
	cm.Add(b)

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
	if !a.Position.Equals(aExpectedPosition) {
		t.Errorf("a.Position should be %V, not %V", aExpectedPosition, a.Position)
	}

	bExpectedPosition := &vector.Vec2{25, 10}
	if !b.Position.Equals(bExpectedPosition) {
		t.Errorf("b.Position should be %V, not %V", bExpectedPosition, b.Position)
	}

	// since a was standing still it should take all the energy
	aExpectedVelocity := &vector.Vec2{-1, 0}
	if !a.Velocity.Equals(aExpectedVelocity) {
		t.Errorf("b.Velocity should be {%f, %f}, not {%f, %f}", aExpectedVelocity.X, aExpectedVelocity.Y, a.Velocity.X, a.Velocity.Y)
	}

	// since a was standing still it should take all the energy and b will stop
	bExpectedVelocity := &vector.Vec2{0, 0}
	if !b.Velocity.Equals(bExpectedVelocity) {
		t.Errorf("b.Velocity should be {%f, %f}, not {%f, %f}", bExpectedVelocity.X, bExpectedVelocity.Y, b.Velocity.X, b.Velocity.Y)
	}
}

// ~50000 28567
func BenchmarkDetectCollisions(bench *testing.B) {
	cm := &CollisionManager{}
	a := NewPhysics(10, 10, 0, 1)
	// standing still
	a.Velocity = &vector.Vec2{0, 0}
	b := NewPhysics(20, 10, 0, 1)
	// moving towards a
	b.Velocity = &vector.Vec2{-1, 0}
	cm.Add(a)
	cm.Add(b)

	bench.ResetTimer()
	for n := 0; n < bench.N; n++ {
		cm.DetectCollisions()
	}
}
