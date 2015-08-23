package physics

import (
	"github.com/stojg/pants/vector"
)

type Contact struct {
	a           Geometry
	b           Geometry
	normal      *vector.Vec2
	penetration float64
	restitution    float64
}

//func (c *Contact) SeparatingVelocity() float64 {
//	relativeVel := c.a.Velocity.Clone()
//	if c.b != nil {
//		relativeVel.Sub(c.b.Velocity)
//	}
//	return relativeVel.Dot(c.normal)
//}

func (c *Contact) Resolve(duration float64) {
	c.resolveVelocity(duration)
	c.resolveInterpenetration()
}

// resolveInterpenetration separates two objects that has penetrated
func (c *Contact) resolveInterpenetration() {

	if c.penetration <= 0 {
		return
	}

//	totalInvMass := c.a.physics.(*ParticlePhysics).InvMass
//	if c.b != nil {
//		totalInvMass += c.b.physics.(*ParticlePhysics).InvMass
//	}
//	// Both objects have infinite mass, so no velocity
//	if totalInvMass == 0 {
//		return
//	}
//
//	movePerIMass := c.normal.Clone().Scale(c.penetration / totalInvMass)
//
//	c.a.Position.Add(movePerIMass.Clone().Scale(c.a.physics.(*ParticlePhysics).InvMass))
//	if c.b != nil {
//		c.b.Position.Add(movePerIMass.Clone().Scale(-c.b.physics.(*ParticlePhysics).InvMass))
//	}
}

// resolveVelocity calculates the new velocity that is the result of the collision
func (collision *Contact) resolveVelocity(duration float64) {
//	// Find the velocity in the direction of the contact normal
//	separatingVelocity := collision.SeparatingVelocity()
//
//	// The objects are already separating, NOP
//	if separatingVelocity > 0 {
//		return
//	}
//
//	// Calculate the new separating velocity
//	newSepVelocity := -separatingVelocity * collision.restitution
//
//	// Check the velocity build up due to acceleration only
//	accCausedVelocity := collision.a.physics.(*ParticlePhysics).forces.Clone()
//	if collision.b != nil {
//		accCausedVelocity.Sub(collision.b.physics.(*ParticlePhysics).forces)
//	}
//
//	// If we have closing velocity due to acceleration buildup,
//	// remove it from the new separating velocity
//	accCausedSepVelocity := accCausedVelocity.Dot(collision.normal) * duration
//	if accCausedSepVelocity < 0 {
//		newSepVelocity += collision.restitution * accCausedSepVelocity
//		// make sure that we haven't removed more than was there to begin with
//		if newSepVelocity < 0 {
//			newSepVelocity = 0
//		}
//	}
//
//	deltaVelocity := newSepVelocity - separatingVelocity
//
//	totalInvMass := collision.a.physics.(*ParticlePhysics).InvMass
//	if collision.b != nil {
//		totalInvMass += collision.b.physics.(*ParticlePhysics).InvMass
//	}
//
//	// Both objects have infinite mass, so they can't actually move
//	if totalInvMass == 0 {
//		return
//	}
//
//	impulsePerIMass := collision.normal.Clone().Scale(deltaVelocity / totalInvMass)
//
//	velocityChangeA := impulsePerIMass.Clone().Scale(collision.a.physics.(*ParticlePhysics).InvMass)
//	collision.a.Velocity.Add(velocityChangeA)
//	if collision.b != nil {
//		velocityChangeB := impulsePerIMass.Clone().Scale(-collision.b.physics.(*ParticlePhysics).InvMass)
//		collision.b.Velocity.Add(velocityChangeB)
//	}
}

