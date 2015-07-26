package main

import (
	"math/rand"
	"time"
)

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

func (w *World) Tick() {
	currentTime := time.Now()
	w.worldTick(currentTime)
	w.networkTick(currentTime)
}

func (w *World) networkTick(currentTime time.Time) {
	duration := currentTime.Sub(w.netTicked)
	if duration < 50*time.Millisecond {
		return
	}
	w.netTicked = currentTime
	w.netTick += 1
	// Send a snapshot to all connections
	h.Send(&Message{Topic: "update", Data: w.spriteList.Changed(true), Tick: w.gameTick, Timestamp: float64(time.Now().UnixNano())/1000000,})
}

func (w *World) worldTick(currentTime time.Time) {
	duration := currentTime.Sub(w.gameTicked)
	if duration < 16*time.Millisecond {
		return
	}
	w.gameTicked = currentTime
	w.gameTick += 1

	// world AI
	w.director.Update(w, duration);

	// Individual entities
	for _, sprite := range w.spriteList.sprites {
		sprite.Update(w, duration)
	}

}
