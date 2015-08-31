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

// Add a physics object to this collision manager
func (cm *CollisionManager) Add(p *Physics) {
	cm.physics = append(cm.physics, p)
}

// Remove a physics object to this collisions manager
func (cm *CollisionManager) Remove(p *Physics) {
	for i, obj := range cm.physics {
		if obj == p {
			cm.physics = append(cm.physics[:i], cm.physics[i+1:]...)
		}
	}
}

// Length returns the total number of physics object in the collision manager
func (cm *CollisionManager) Length() int {
	return len(cm.physics)
}

func (cm *CollisionManager) GenerateContacts(a, b *Physics) (*Contact, error) {
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

	cm.collisions = make([]*Contact, 0)
	// todo(stig): do broad phase here

	width := 1000.0
	height := 1000.0
	cellSize := 20.0

	cols := int(width / 20)
	rows := int(height / 20)

	totalCells := cols * rows

	entityWidth := 20.0
	entityHeight := 20.0

	// create the grid
	var grid [][]*Physics
	grid = make([][]*Physics, totalCells)

	//	hashChecks := 0

	for _, p := range cm.physics {
		// if the entity is outside the grid ignore it
		if p.Position.X-entityWidth/2.0 < 0 || p.Position.X+entityWidth/2.0 > width || p.Position.Y-entityHeight/2.0 < 0 || p.Position.Y+entityHeight/2.0 > height {
			continue
		}
		// find extremes of cells that entity overlaps
		// subtract min to shift grid to avoid negative numbers
		entityMinX := math.Floor(((p.Position.X - entityWidth/2.0) - 0) / cellSize)
		entityMaxX := math.Floor(((p.Position.X + entityWidth/2.0) - 0) / cellSize)
		entityMinY := math.Floor(((p.Position.Y - entityHeight/2.0) - 0) / cellSize)
		entityMaxY := math.Floor(((p.Position.Y + entityHeight/2.0) - 0) / cellSize)

		// insert entity into each cell it overlaps
		// we're looping to make sure that all cells between extremes are found
		for cX := entityMinX; cX <= entityMaxX; cX++ {
			for cY := entityMinY; cY <= entityMaxY; cY++ {
				index := int(cY)*cols + int(cX)
				grid[index] = append(grid[index], p)
			}
		}
	}

	checked := make(map[*Physics]map[*Physics]bool)

	// alright, we got all of the suckers in a nice array, it's time to see if
	// things are colliding
	for _, list := range grid {
		// no entities in this grid
		if len(list) < 2 {
			continue
		}

		for i := 0; i < len(list); i++ {
			a := list[i]
			for j := i + 1; j < len(list); j++ {
				b := list[j]

				if checked := checked[a][b]; checked {
					continue
				}
				if checked := checked[b][a]; checked {
					continue
				}

				collision, err := cm.GenerateContacts(a, b)
				if err != nil {
					log.Printf("error: %s", err)
				}
				if collision != nil {
					collision.a = a
					collision.b = b
					collision.restitution = 0.1
					cm.collisions = append(cm.collisions, collision)
				}
				if checked[a] == nil {
					checked[a] = make(map[*Physics]bool, 1)
				}
				checked[a][b] = true
				if checked[b] == nil {
					checked[b] = make(map[*Physics]bool, 1)
				}
				checked[b][a] = true
			}
		}
	}

	if len(cm.collisions) > 0 {
		return true
	}
	return false

}

// ResolveCollision will resolve all collisions found by DetectCollisions
func (cm *CollisionManager) ResolveCollisions(duration float64) {

	for _, collision := range cm.collisions {
		collision.Resolve(duration)
	}
	return

	numContacts := len(cm.collisions)
	//	log.Printf("collusions %d", numContacts)
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
			//			log.Printf("yeah nah")
			break
		}

		// resolve the collision with the lowest separating velocity
		cm.collisions[maxIndex].Resolve(duration)
		iterationsUsed += 1
	}
	cm.collisions = make([]*Contact, 0)
}
