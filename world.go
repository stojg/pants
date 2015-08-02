package main

import (
	"log"
	"math/rand"
	"time"
)

const NetFPS = 50 * time.Millisecond

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
			log.Printf("net lag %s", duration-NetFPS)
		}
		w.netTicked = currentTime
		w.netTick += 1
		// Send a snapshot to all connections
		h.Send(&Message{Topic: "update", Data: w.spriteList.Changed(true), Tick: w.gameTick, Timestamp: float64(time.Now().UnixNano()) / 1000000})
	}
}

func (w *World) worldTick() {

	var gameTime float64
	dt := 0.01
	currentTime := time.Now()
	var accumulator float64

	for {
		newTime := time.Now()
		frameTime := newTime.Sub(currentTime).Seconds()
		currentTime = newTime

		if frameTime > 0.016 {
			log.Printf("world lag: %d ms", int(frameTime*1000))
		}

		accumulator += frameTime
		for accumulator >= dt {
			// Individual entities
			w.spriteList.Update(w, dt, gameTime)
			accumulator -= (dt)
			gameTime += dt
		}

		// world AI
		w.director.Update(w, dt, gameTime)
	}

}
