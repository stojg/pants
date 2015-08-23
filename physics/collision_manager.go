package physics

import (
	"fmt"
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

func (cm *CollisionManager) DetectCollisions(entities map[uint64]*Physics) {
	// todo(stig): do broad phase detection here

	cm.collisions = make([]*Contact,0)
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

			collision, err:= cm.Contact(cm.Geometry(a), cm.Geometry(b))
			if err != nil {
				log.Printf("error: %s", err)
			}
			if collision != nil {
				cm.collisions = append(cm.collisions, collision)
			}
			checked[idx][idy] = true
		}
	}
}

func (cm *CollisionManager) ResolveCollisions(duration float64) {
	// todo(stig): implement contact resolution here
	for _, _ = range cm.collisions {
	}
}

func (cm *CollisionManager) Geometry(p *Physics) (Geometry) {
	circle := &Circle{
		position: p.Position.Clone(),
		radius:   10,
	}
	return circle
}
