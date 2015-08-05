package main

import (
	//	"log"
	"log"
)

type Director struct {
	lastAction float64
}

func (director *Director) Update(w *World, frameTime, gameTime float64 ) {

	if gameTime - director.lastAction < 0.1 {
		return
	}

	if len(list.sprites) >= 100  {
		return;
	}
	director.lastAction = gameTime
	w.entities.NewEntity(
		w.rand.Float64()*800,
		w.rand.Float64()*600,
		"assets/basics/arrow.png",
	)
	log.Printf("%d", len(w.entities.sprites))
}
