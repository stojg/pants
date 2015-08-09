package collision

import (
	. "github.com/stojg/pants/physics"
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

func TestCMUpdate(t *testing.T) {
	cm := &CollisionManager{}
	cm.Update(0.016)
}

func TestGetCollisionGeometry(t *testing.T) {
	cm := &CollisionManager{}
	p1 := NewPhysics(0, 0, 0, 0)
	geometry, err := cm.getCollisionGeometry(p1)
	if err != nil {
		t.Errorf("getCollisionGeometry returned error: '%s'", err)
	}
	if ty, ok := geometry.(*Circle); !ok {
		t.Errorf("getCollisionGeometry should return a *CollisionCircle, not %s", ty)
	}
}

func TestContactPairHit(t *testing.T) {
	cm := &CollisionManager{}
	g1, _ := cm.getCollisionGeometry(NewPhysics(0, 0, 0, 0))
	g2, _ := cm.getCollisionGeometry(NewPhysics(0, 0, 0, 0))
	_, err := cm.ContactPair(g1, g2)
	if err != nil {
		t.Errorf("Error reported %s", err)
		return
	}
	//	actual := contact.a
	//	if actual != g1 {
	//		t.Errorf("Contact should report %v, not %v", g1, actual)
	//	}
}
