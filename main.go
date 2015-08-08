package main

import (
	"log"
	_ "net/http/pprof"
	"runtime"
)

var list *EntityList

func main() {
	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)
	log.Printf("running on %d CPUs", nCPU)

	world := NewWorld(list)
	webserver := webserver{port: "8081"}

	world.Run()
	go h.run()
	webserver.Start()
}
