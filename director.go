package main

import ()

type Director struct {
	lastAction float64
}

func (director *Director) Update(w *World, frameTime, gameTime float64) {

	if gameTime-director.lastAction < 0.1 {
		return
	}

	if len(entityManager.sprites) >= 20 {
		return
	}
	director.lastAction = gameTime
	w.entityManager.NewEntity(
		w.rand.Float64()*800,
		w.rand.Float64()*600,
		"assets/basics/arrow.png",
	)
}
