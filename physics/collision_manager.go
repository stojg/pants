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
	return nil, fmt.Errorf("No Collision test could be done between between %T and %T", a.collisionGeometry, b.collisionGeometry)
}

// DetectCollisions checks all managed Physics and checks for collisions
// if they to collide the collision contact will be stored for later usage by
// ResolveCollisions
func (cm *CollisionManager) DetectCollisions() bool {

	cm.collisions = make([]*Contact, 0)

	width := 1000.0
	height := 1000.0
	cellSize := 20.0

	cols := int(width / 20)

	entityWidth := 20.0
	entityHeight := 20.0

	// create the grid = {cols} * {height/cellsize}
	var grid [6000][]*Physics

	for _, p := range cm.physics {

		// pre-calculate these values for performance reasonss
		entityMinX := p.Position.X - entityWidth/2.0
		entityMaxX := p.Position.X + entityWidth/2.0
		entityMinY := p.Position.Y - entityHeight/2.0
		entityMaxY := p.Position.Y + entityHeight/2.0

		// entities outside of the grid should not be checked
		if entityMinX < 0 || entityMaxX > width || entityMinY < 0 || entityMaxY > height {
			continue
		}

		// find extremes of cells that entity overlaps
		// subtract min to shift grid to avoid negative numbers
		gridMinSizeX := 0.0
		gridMinSizeY := 0.0
		entityMinCellX := int((entityMinX - gridMinSizeX) / cellSize)
		entityMaxCellX := int((entityMaxX - gridMinSizeX) / cellSize)
		entityMinCellY := int((entityMinY - gridMinSizeY) / cellSize)
		entityMaxCellY := int((entityMaxY - gridMinSizeY) / cellSize)

		// insert entity into each cell it overlaps
		// loop to make sure that all cells between extremes are found
		for cellX := entityMinCellX; cellX <= entityMaxCellX; cellX++ {
			for cellY := entityMinCellY; cellY <= entityMaxCellY; cellY++ {
				index := cellY*cols + cellX
				grid[index] = append(grid[index], p)
			}
		}
	}

	checked := make(map[*Physics]map[*Physics]bool)

	// alright, we got all of the suckers in a nice array, it's time to see if
	// they are colliding
	for _, cell := range grid {
		// if there is zero or one entity in this cell there can't be a collision
		if len(cell) < 2 {
			continue
		}

		// iterate over the entities in the cell
		for i := 0; i < len(cell); i++ {
			a := cell[i]
			// compare a with all the other entities 'after' it in the cell
			for j := i + 1; j < len(cell); j++ {
				b := cell[j]
				// have this pair been checked before?
				if checked := checked[a][b]; checked {
					continue
				}
				if checked := checked[b][a]; checked {
					continue
				}
				// do the narrow phase collision detection
				collision, err := cm.GenerateContacts(a, b)
				if err != nil {
					log.Printf("error: %s", err)
				}
				// we got a proper contact, populate contact data with more data
				if collision != nil {
					collision.a = a
					collision.b = b
					collision.restitution = 0.1
					cm.collisions = append(cm.collisions, collision)
				}

				// if the checked maps not yet setup, initialise them here
				if checked[a] == nil {
					checked[a] = make(map[*Physics]bool, 1)
				}
				if checked[b] == nil {
					checked[b] = make(map[*Physics]bool, 1)
				}

				// mark this pair as checked
				checked[a][b] = true
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

	numContacts := len(cm.collisions)
	iterationsUsed := 0
	iterations := numContacts * 2

	// find the collision with the highest closing velocity and resolve that one
	// first.
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
	cm.collisions = make([]*Contact, 0)
}
