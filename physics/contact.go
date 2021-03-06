package physics

import (
	"github.com/stojg/pants/vector"
)

type Contact struct {
	a           *Physics
	aId         uint64
	b           *Physics
	bId         uint64
	normal      *vector.Vec2
	penetration float64
	restitution float64
}

func (c *Contact) Resolve(duration float64) {
	c.resolveVelocity(duration)
	c.resolveInterpenetration()
}

func (c *Contact) Pair() (uint64, uint64) {
	return c.aId, c.bId
}

// resolveInterpenetration separates two objects that has penetrated
func (c *Contact) resolveInterpenetration() {

	// These objects are not intersecting
	if c.penetration <= 0 {
		return
	}

	// Calculate the total inverse mass
	totalInvMass := c.a.Data.InvMass
	if c.b != nil {
		totalInvMass += c.b.Data.InvMass
	}

	// Both objects have infinite mass, so neither can be moved
	if totalInvMass == 0 {
		return
	}

	movePerInvMass := c.normal.NewScale(c.penetration / totalInvMass)
	// Move the objects out of contact depending on their mass
	c.a.Data.Position.Add(movePerInvMass.NewScale(c.a.Data.InvMass))
	if c.b.Data != nil {
		c.b.Data.Position.Add(movePerInvMass.NewScale(-c.b.Data.InvMass))
	}

	// set the penetration to zero
	c.penetration = 0
}

// resolveVelocity calculates the new velocity that is the result of the collision
func (collision *Contact) resolveVelocity(duration float64) {

	// Find the velocity in the direction of the contact normal
	separatingVelocity := collision.separatingVelocity()

	// check if it needs to be resolved
	if separatingVelocity > 0 {
		// the contact is either separating or stationary so no impulse is required
		return
	}

	// Calculate the new separating velocity
	newSepVelocity := -separatingVelocity * collision.restitution

	// Check the velocity build up due to acceleration only
	accCausedVelocity := collision.a.Data.Forces.Clone()
	if collision.b != nil {
		accCausedVelocity.Sub(collision.b.Data.Forces)
	}

	// If we've got closing velocity due to acceleration buildup,
	// remove it from the new separating velocity
	accCausedSepVelocity := accCausedVelocity.Dot(collision.normal) * duration
	if accCausedSepVelocity < 0 {
		newSepVelocity += collision.restitution * accCausedSepVelocity
		// make sure that we haven't removed more than was there to begin with
		if newSepVelocity < 0 {
			newSepVelocity = 0
		}
	}

	deltaVelocity := newSepVelocity - separatingVelocity

	// We apply the change of velocity to each object in proportion to
	// their inverse mass (i.e., those with lower inverse mass (higher mass)
	// get less change in velocity
	totalInvMass := collision.a.Data.InvMass
	if collision.b != nil {
		totalInvMass += collision.b.Data.InvMass
	}

	// If all objects have infinite mass the will not react to any impulses
	if totalInvMass == 0 {
		return
	}

	// calculate the total impulse to apply
	impulse := deltaVelocity / totalInvMass

	// find the impulse direction by scaling it in the direction of the normal
	impulseWithDirection := collision.normal.NewScale(impulse)

	// change velocity for a in proportion to its inverse mass
	velocityChangeA := impulseWithDirection.NewScale(collision.a.Data.InvMass)
	collision.a.Data.Velocity.Add(velocityChangeA)

	// change velocity for b in proportion to its inverse mass
	if collision.b != nil {
		velocityChangeB := impulseWithDirection.NewScale(-collision.b.Data.InvMass)
		collision.b.Data.Velocity.Add(velocityChangeB)
	}
}

func (c *Contact) separatingVelocity() float64 {
	relativeVel := c.a.Data.Velocity.Clone()
	if c.b != nil {
		relativeVel.Sub(c.b.Data.Velocity)
	}
	return relativeVel.Dot(c.normal)
}
