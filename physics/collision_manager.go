package physics

import (
	"fmt"
	. "github.com/stojg/pants/structs"
	"log"
	"math"
)

func NewCollisionManager() *CollisionManager {
	return &CollisionManager{}
}

// CollisionManager does collision testing between Physics components
type CollisionManager struct{}

func (cm *CollisionManager) GenerateContacts(a, b *Physics) (*Contact, error) {
	switch a.CollisionGeometry.(type) {
	case *Circle:
		switch b.CollisionGeometry.(type) {
		case *Circle:
			return CircleVsCircle(a.CollisionGeometry.(*Circle), b.CollisionGeometry.(*Circle)), nil
		}
	}
	return nil, fmt.Errorf("No Collision test could be done between between %T and %T", a.CollisionGeometry, b.CollisionGeometry)
}

// DetectCollisions checks all managed Physics and checks for collisions
// if they to collide the collision contact will be stored for later usage by
// ResolveCollisions
func (cm *CollisionManager) DetectCollisions(entities map[uint64]*Physics, grid *Grid) []*Contact {

	collisionList := make([]*Contact, 0)

	for id, entity := range entities {
		d := entity.Data
		if d.Height > d.Width {
			entity.CollisionGeometry.(*Circle).Radius = d.Height / 2
		} else {
			entity.CollisionGeometry.(*Circle).Radius = d.Width / 2
		}
		// pre-calculate these values for performance reasons
		entityMinX := int(d.Position.X - d.Width/2)
		entityMaxX := int(d.Position.X + d.Width/2)
		entityMinY := int(d.Position.Y - d.Height/2)
		entityMaxY := int(d.Position.Y + d.Height/2)

		// entities outside of the grid should not be checked
		if entityMinX < 0 || entityMaxX >= grid.Width || entityMinY < 0 || entityMaxY >= grid.Height {
			continue
		}

		// find extremes of cells that entity overlaps
		// subtract min to shift grid to avoid negative numbers
		gridMinSizeX := 0
		gridMinSizeY := 0
		entityMinCellX := (entityMinX - gridMinSizeX) / grid.CellSize
		entityMaxCellX := (entityMaxX - gridMinSizeX) / grid.CellSize
		entityMinCellY := (entityMinY - gridMinSizeY) / grid.CellSize
		entityMaxCellY := (entityMaxY - gridMinSizeY) / grid.CellSize

		// insert entity into each cell it overlaps
		// loop to make sure that all cells between extremes are found
		for cellX := entityMinCellX; cellX <= entityMaxCellX; cellX++ {
			for cellY := entityMinCellY; cellY <= entityMaxCellY; cellY++ {
				index := cellY*grid.Cols + cellX
				grid.Cells[index] = append(grid.Cells[index], id)
			}
		}
	}

	checked := make(map[uint64]map[uint64]bool)

	// alright, we got all of the suckers in a nice array, it's time to see if
	// they are colliding
	for _, cell := range grid.Cells {
		// if there is zero or one entity in this cell there can't be a collision
		if len(cell) < 2 {
			continue
		}

		// iterate over the entities in the cell
		for i := 0; i < len(cell); i++ {
			aId := cell[i]
			// compare a with all the other entities 'after' it in the cell
			for j := i + 1; j < len(cell); j++ {
				bId := cell[j]
				// have this pair been checked before?
				if checked := checked[aId][bId]; checked {
					continue
				}
				if checked := checked[bId][aId]; checked {
					continue
				}
				a := entities[aId]
				b := entities[bId]
				// do the narrow phase collision detection
				collision, err := cm.GenerateContacts(a, b)
				if err != nil {
					log.Printf("error: %s", err)
				}
				// we got a proper contact, populate contact data with more data
				if collision != nil {
					collision.a = a
					collision.aId = aId
					collision.b = b
					collision.bId = bId
					collision.restitution = 0.1
					collisionList = append(collisionList, collision)
				}

				// if the checked maps not yet setup, initialise them here
				if checked[aId] == nil {
					checked[aId] = make(map[uint64]bool, 1)
				}
				if checked[bId] == nil {
					checked[bId] = make(map[uint64]bool, 1)
				}

				// mark this pair as checked
				checked[aId][bId] = true
				checked[bId][aId] = true
			}
		}
	}

	return collisionList
}

// ResolveCollision will resolve all collisions found by DetectCollisions
func (cm *CollisionManager) ResolveCollisions(collisions []*Contact, duration float64) {

	numContacts := len(collisions)
	iterationsUsed := 0
	iterations := numContacts * 2

	// find the collision with the highest closing velocity and resolve that one
	// first.
	for iterationsUsed < iterations {
		// find the contact with the largest closing velocity
		max := math.MaxFloat64
		maxIndex := numContacts
		for i, collision := range collisions {
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
		collisions[maxIndex].Resolve(duration)
		iterationsUsed += 1
	}
}
