package main

import (
	"time"
	//	"log"
	"log"
)

type Director struct {
}

func (director *Director) Update(w *World, duration time.Duration) {
	if w.rand.Float32() > 0.99 {
		w.spriteList.NewSprite(
			w.rand.Float64()*800,
			w.rand.Float64()*600,
			"assets/basics/arrow.png",
		)
		log.Printf("%d", len(w.spriteList.sprites))
	}
}
