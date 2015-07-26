package main

import (
	"log"
	"math/rand"
	"time"
)

const NetFPS = 50 * time.Millisecond
const FPS = 16 * time.Millisecond

func NewWorld(list *SpriteList) *World {
	current := time.Now()
	return &World{
		netTicked:  current,
		gameTicked: current,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
		spriteList: list,
		director:   &Director{},
	}
}

type World struct {
	netTicked  time.Time
	netTick    uint64
	gameTicked time.Time
	gameTick   uint64
	spriteList *SpriteList
	rand       *rand.Rand
	director   *Director
}

func (w *World) Run() {
	go w.worldTick()
	go w.networkTick()
}

func (w *World) networkTick() {
	ticker := time.NewTicker(NetFPS)
	for currentTime := range ticker.C {
		duration := currentTime.Sub(w.netTicked)
		if duration-NetFPS > 4*time.Millisecond {
			log.Printf("net is lagging %s", duration-NetFPS)
		}
		w.netTicked = currentTime
		w.netTick += 1
		// Send a snapshot to all connections
		h.Send(&Message{Topic: "update", Data: w.spriteList.Changed(true), Tick: w.gameTick, Timestamp: float64(time.Now().UnixNano()) / 1000000})
	}
}

func (w *World) worldTick() {
	ticker := time.NewTicker(FPS)
	for currentTime := range ticker.C {
		duration := currentTime.Sub(w.gameTicked)
		if duration-FPS > 4*time.Millisecond {
			log.Printf("world is lagging %s", duration-FPS)
		}
		w.gameTicked = currentTime
		w.gameTick += 1

		// world AI
		w.director.Update(w, duration)

		// Individual entities
		for _, sprite := range w.spriteList.sprites {
			sprite.Update(w, duration)
		}
	}
}
