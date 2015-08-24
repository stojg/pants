package physics

import (
	"fmt"
	"log"
	"math"
)

// CollisionManager keeps tracks and does collision testing between Physics
// components
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
			cm.physics = append(cm.physics[:i], cm.physics[i + 1:]...)
		}
	}
}

func (cm *CollisionManager) Length() int {
	return len(cm.physics)
}

func (cm *CollisionManager) Collision(a, b *Physics) (*Contact, error) {
	switch a.collisionGeometry.(type) {
	case *Circle:
		switch b.collisionGeometry.(type) {
		case *Circle:
			return CircleVsCircle(a.collisionGeometry.(*Circle), b.collisionGeometry.(*Circle)), nil
		}
	}
	return nil, fmt.Errorf("No Collision check could be done between between %s and %s", a.collisionGeometry, b.collisionGeometry)
}

// DetectCollisions checks all managed Physics and checks for collisions
// if they to collide the collision contact will be stored for later usage by
// ResolveCollisions
func (cm *CollisionManager) DetectCollisions() bool {
	// todo(stig): do broad phase  here

	cm.collisions = make([]*Contact, 0)
	checked := make(map[int]map[int]bool)

	// todo(stig): implement narrow phase here
	for idx, a := range cm.physics {
		checked[idx] = make(map[int]bool)
		for idy, b := range cm.physics {
			if a == b {
				continue
			}
			if checked := checked[idx][idy]; checked {
				continue
			}
			collision, err := cm.Collision(a, b)
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

// ResolveCollision will resolve all collisions found by DetectCollisions
func (cm *CollisionManager) ResolveCollisions(duration float64) {

	numContacts := len(cm.collisions)
	iterationsUsed := 0
	iterations := numContacts * 2

	for iterationsUsed < iterations {
		// find the contact with the largest closing velocity
		max := math.MaxFloat64
		maxIndex := numContacts
		for i, collision := range cm.collisions {
			sepVel := collision.separatingVelocity()
			// we found a collision with the lowest separation velocity that is
			// intersection
			if sepVel < max && (sepVel < 0 || collision.penetration > 0) {
				max = sepVel
				maxIndex = i
			}
		}

		// do we have anything worth resolving
		if maxIndex == numContacts {
			break
		}

		// resolve the collision with the lowest separating velocity
		cm.collisions[maxIndex].Resolve(duration)
		iterationsUsed += 1
	}
}
