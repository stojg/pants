package main

import (
	. "github.com/stojg/pants/vector"
	"math"
)

type ForceGenerator interface {
	UpdateForce(*Physics, float64)
}

type StaticForce struct {
	force *Vec2
}

func (pg *StaticForce) UpdateForce(p *Physics, duration float64) {
	if p.invMass == 0 {
		return
	}
	p.AddForce(pg.force.Multiply(1 / p.invMass))
}

type ParticleForceRegistration struct {
	particle *Physics
	fg       ForceGenerator
}

type SpringForce struct {
	other           *Physics
	springConstance float64
	restLength      float64
}

func (sf *SpringForce) UpdateForce(p *Physics, duration float64) {
	// calculate the vector of the spring
	force := &Vec2{}
	p.getPosition(force)
	force.Sub(sf.other.Position)

	// calculate the magnitude of the force
	magnitude := force.Length()

	magnitude = math.Abs(magnitude - sf.restLength)
	magnitude *= sf.springConstance

	// calculate the final force and apply it
	force.Normalize()
	force.Scale(-magnitude)
	p.AddForce(force)
}

type PhysicsForceRegistry struct {
	registrations []*ParticleForceRegistration
}

func (pfr *PhysicsForceRegistry) Add(p *Physics, fg ForceGenerator) {
	pfr.registrations = append(pfr.registrations, &ParticleForceRegistration{
		particle: p,
		fg:       fg,
	})
}

func (pfr *PhysicsForceRegistry) Remove(p *Physics, fg ForceGenerator) {
	panic("not implemented")
}

func (pfr *PhysicsForceRegistry) Clear() {
	pfr.registrations = pfr.registrations[0:0]
}

func (pfr *PhysicsForceRegistry) Update(duration float64) {
	for _, i := range pfr.registrations {
		i.fg.UpdateForce(i.particle, duration)
	}
}
