package main

import (
	"log"
	"runtime"
)

var list *SpriteList

func main() {
	w := NewWorld(list)
	ws := webserver{port: "8081"}

	go loop(w)
	go h.run()
	ws.Start()
}

func loop(w *World) {
	for {
		w.Tick()
	}
}
