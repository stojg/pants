package main

import ()

type Director struct {
	lastAction float64
}

func (director *Director) Update(w *World, frameTime, gameTime float64) {

	if gameTime-director.lastAction < 0.1 {
		return
	}

	if len(entityManager.sprites) >= 50 {
		return
	}
	director.lastAction = gameTime
	w.entityManager.NewEntity(
		w.rand.Float64()*1000,
		w.rand.Float64()*1000,
		"assets/basics/arrow.png",
	)
}
