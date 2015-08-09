package collision

import (
	"fmt"
	"github.com/stojg/pants/vector"
	. "github.com/stojg/pants/physics"
)
type CollisionGeometry interface {}

type CollisionCircle struct {
	position *vector.Vec2
	radius float64
}

type CollisionManager struct {
	physics []*Physics
}

func (cm *CollisionManager) Add(p *Physics){
	cm.physics = append(cm.physics, p)
}

func (cm *CollisionManager) Remove(p *Physics){
	for i, obj := range cm.physics {
		if obj == p {
			cm.physics = append(cm.physics[:i], cm.physics[i+1:]...)
		}
	}
}

func (cm *CollisionManager) Length() int {
	return len(cm.physics)
}

func (cm *CollisionManager) Update(duration float64) {

}

func (cm *CollisionManager) getCollisionGeometry(p *Physics) (CollisionGeometry, error) {

	circle := &CollisionCircle{
		position: p.Position.Clone(),
		radius: 20,
	}
	return circle, nil

	return nil, fmt.Errorf("No collision geometry found for %v", p)
}