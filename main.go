package main

import (
	"github.com/stojg/pants/network"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
)

var entityManager *EntityManager

func main() {
	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)
	log.Printf("running on %d CPUs", nCPU)

	server := network.NewServer("8081")
	server.Start()

	world := NewWorld(entityManager, server)

	go world.worldTick()
	go world.networkTick()

	// wait for signal
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			log.Printf("Received an interrupt, stopping services...\n")
			server.Stop()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
