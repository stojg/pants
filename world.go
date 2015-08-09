package main

import (
	. "github.com/stojg/pants/vector"
	. "github.com/stojg/pants/physics"
	"log"
	"math/rand"
	"time"
)

const NetFPS = 50 * time.Millisecond

func NewWorld(list *EntityManager) *World {
	current := time.Now()
	return &World{
		netTicked:  current,
		gameTicked: current,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
		//		rand:       rand.New(rand.NewSource(1)),
		entityManager: list,
		director: &Director{},
	}
}

type World struct {
	netTicked  time.Time
	netTick    uint64
	gameTicked time.Time
	gameTick   uint64
	entityManager *EntityManager
	rand       *rand.Rand
	director   *Director
	debug      []*Line
}

type EntityUpdate struct {
	Id          uint64            `bson:",minsize"`
	X, Y        float64           `bson:",minsize,omitempty"`
	Orientation float64           `bson:",minsize"`
	Type        string            `bson:",minsize,omitempty"`
	Dead        bool              `bson:",minsize,omitempty"`
	Data        map[string]string `bson:",minsize,omitempty"`
}

func (w *World) RandF64(x int) float64 {
	return w.rand.Float64() * 800
}

func (w *World) Physic(id uint64) *Physics {
	return w.entityManager.physics[id]
}

func (w *World) Line(start, end *Vec2) {
	w.debug = append(w.debug, &Line{
		Position: start,
		End:      end,
	})
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

		h.SendUpdates(entityManager.updated, w, currentTime)
		entityManager.updated = make(map[uint64]bool)
	}
}

func (w *World) worldTick() {
	var gameTime float64
	var accumulator float64
	dt := 0.01
	currentTime := time.Now()
	for {
		newTime := time.Now()
		frameTime := newTime.Sub(currentTime).Seconds()
		currentTime = newTime

		if frameTime > 0.016 {
			log.Printf("world lag: %d ms", int((frameTime-0.016)*1000))
		}

		accumulator += frameTime
		for accumulator >= dt {
			// Individual entities
			w.entityManager.Update(w, dt, gameTime)
			accumulator -= (dt)
			gameTime += dt
		}

		// world AI
		w.director.Update(w, dt, gameTime)
	}

}
