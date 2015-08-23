package physics

import (
	"log"
)

type CollisionManager struct {
	physics    []*Physics
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

func (cm *CollisionManager) Contact(g1, g2 *Physics) (*Contact, error) {
	switch g1.collisionGeometry.(type) {
	case *Circle:
		switch g2.collisionGeometry.(type) {
		case *Circle:
			return CircleVsCircle(g1.collisionGeometry.(*Circle), g2.collisionGeometry.(*Circle)), nil
		}
	}
	panic("No contact pair found")
}

func (cm *CollisionManager) DetectCollisions() bool {
	// todo(stig): do broad phase detection here

	cm.collisions = make([]*Contact, 0)
	checked := make(map[int]map[int]bool)

	// todo(stig): implement fine phase detection here
	for idx, a := range cm.physics {
		checked[idx] = make(map[int]bool)
		for idy, b := range cm.physics {
			if a == b {
				continue
			}
			if checked := checked[idx][idy]; checked {
				continue
			}
			collision, err := cm.Contact(a, b)
			if err != nil {
				log.Printf("error: %s", err)
			}
			if collision != nil {
				collision.a = a
				collision.b = b
				collision.restitution = 0.1
				cm.collisions = append(cm.collisions, collision)
			}
			checked[idx][idy] = true
		}
	}

	if len(cm.collisions) > 0 {
		return true
	}
	return false

}

func (cm *CollisionManager) ResolveCollisions(duration float64) {
	// todo(stig): implement contact resolution here
	for _, collision := range cm.collisions {
		collision.Resolve(duration)
	}
}

func (cm *CollisionManager) Geometry(p *Physics) Geometry {
	circle := &Circle{
		position: p.Position.Clone(),
		radius:   10,
	}
	return circle
}
