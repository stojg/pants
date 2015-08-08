package main

import (
	. "github.com/stojg/pants/vector"
)

type ForceGenerator interface {
	UpdateForce(*Physics, float64)
}

type ParticleGravity struct {
	gravity *Vec2
}

func (pg *ParticleGravity) UpdateForce(p *Physics, duration float64) {
	if p.invMass == 0 {
		return
	}
	p.AddForce(pg.gravity.Multiply(1/p.invMass))
}

type ParticleForceRegistration struct {
	particle *Physics
	fg ForceGenerator
}

type PhysicsForceRegistry struct {
	registrations []*ParticleForceRegistration
}

func (pfr *PhysicsForceRegistry) Add(p *Physics, fg ForceGenerator) {
	pfr.registrations = append(pfr.registrations, &ParticleForceRegistration{
		particle: p,
		fg: fg,
	})
}

func (pfr *PhysicsForceRegistry) Remove(p *Physics, fg ForceGenerator) {
	panic("not implemented")
}

func (pfr *PhysicsForceRegistry) Clear() {
	pfr.registrations = pfr.registrations[0:0]
}

func (pfr *PhysicsForceRegistry) Update(duration float64) {
	for _, i:= range pfr.registrations {
		i.fg.UpdateForce(i.particle, duration)
	}
}
