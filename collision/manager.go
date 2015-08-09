package collision

import (
	"fmt"
	. "github.com/stojg/pants/physics"
)

type CollisionManager struct {
	physics []*Physics
	collisions []*Contact
}

func (cm *CollisionManager) Add(p *Physics) {
	cm.physics = append(cm.physics, p)
}

func (cm *CollisionManager) Remove(p *Physics) {
	for i, obj := range cm.physics {
		if obj == p {
			cm.physics = append(cm.physics[:i], cm.physics[i+1:]...)
		}
	}
}

func (cm *CollisionManager) Length() int {
	return len(cm.physics)
}

func (cm *CollisionManager) Contact(g1, g2 Geometry) (*Contact, error) {
	switch g1.(type) {
	case *Circle:
		switch g2.(type) {
		case *Circle:
			return CircleVsCircle(g1.(*Circle), g2.(*Circle)), nil
		}
	}
	return nil, fmt.Errorf("No contact pair found")
}

func (cm *CollisionManager) DetectCollisions(duration float64) {
	// todo(stig): do broad phase detection here

	// todo(stig): implement fine phase detection here
	for _, _ = range cm.physics {
	}
}

func (cm *CollisionManager) ResolveCollisions(duration float64) {
	// todo(stig): implement contact resolution here
	for _, _ = range cm.collisions {
	}
}

func (cm *CollisionManager) Geometry(p *Physics) (Geometry, error) {
	circle := &Circle{
		position: p.Position.Clone(),
		radius:   10,
	}
	return circle, nil
}
