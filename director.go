package main

import (
	//	"log"
	"log"
)

type Director struct {
	lastAction float64
}

func (director *Director) Update(w *World, frameTime, gameTime float64 ) {

	if gameTime - director.lastAction < 1 {
		return
	}
	director.lastAction = gameTime
	w.spriteList.NewSprite(
		w.rand.Float64()*800,
		w.rand.Float64()*600,
		"assets/basics/arrow.png",
	)
	log.Printf("%d", len(w.spriteList.sprites))
}
