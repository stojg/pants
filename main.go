package main

import (
	"log"
	_ "net/http/pprof"
	"runtime"
)

var entityManager *EntityManager

func main() {
	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)
	log.Printf("running on %d CPUs", nCPU)

	world := NewWorld(entityManager)
	webserver := webserver{port: "8081"}

	go world.worldTick()
	go world.networkTick()
	go h.run(entityManager)
	webserver.Start()
}
