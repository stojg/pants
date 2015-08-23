package physics

import (
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
	g1 := NewPhysics(0, 0, 0, 0)
	g2 := NewPhysics(0, 0, 0, 0)
	_, err := cm.Contact(g1, g2)
	if err != nil {
		t.Errorf("Error reported %s", err)
		return
	}
	//	actual := contact.a
	//	if actual != g1 {
	//		t.Errorf("Contact should report %v, not %v", g1, actual)
	//	}
}
