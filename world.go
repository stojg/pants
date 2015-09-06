package main

import (
	. "github.com/stojg/pants/grid"
	"github.com/stojg/pants/network"
	. "github.com/stojg/pants/physics"
	"github.com/stojg/pants/timer"
	. "github.com/stojg/pants/vector"
	"log"
	"math"
	"math/rand"
	"time"
)

const NetFPS = 50 * time.Millisecond
const WORLD_FRAME_TIME = 0.016

var events *EventManager

func init() {
	events = NewEventManager()
}

func NewWorld(list *EntityManager, s *network.Server) *World {
	currentTime := time.Now()

	w := &World{
		netTicked:     currentTime,
		gameTicked:    currentTime,
		rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
		entityManager: list,
		director:      NewDirector(),
		server:        s,
	}

	w.networkManager = &NetworkManager{server: s, world: w}

	wMap := NewGrid(60, 60)

	for x := 0; x < wMap.Width; x++ {
		for y := 0; y < wMap.Height; y++ {
			v := math.Floor(w.rand.Float64() * 3)
			wMap.Place(x, y, byte(v))
		}
	}
	w.worldMap = wMap
	return w
}

type World struct {
	netTicked      time.Time
	netTick        uint64
	gameTicked     time.Time
	gameTick       uint64
	entityManager  *EntityManager
	networkManager *NetworkManager
	rand           *rand.Rand
	director       *Director
	debug          []*Line
	server         *network.Server
	worldMap       *Grid
}

type EntityUpdate struct {
	Id                  uint64     `bson:",minsize"`
	X, Y, Height, Width float64    `bson:",minsize,omitempty"`
	Orientation         float64    `bson:",minsize"`
	Type                EntityType `bson:",minsize"`
	Dead                bool       `bson:",minsize,omitempty"`
	Properties          *EntityProperty
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
		if duration-NetFPS > 5*time.Millisecond {
			log.Printf("net lag %s", duration-NetFPS)
		}
		w.netTicked = currentTime
		w.netTick += 1

		// handle inputs
		w.networkManager.Inputs()
		w.networkManager.SendState(entityManager.updated, w, currentTime)

		w.entityManager.updated = make(map[uint64]bool, 0)
		w.debug = make([]*Line, 0)
	}
}

func (w *World) worldTick() {
	var gameTime float64
	var accumulator float64
	var frameTime float64
	dt := 0.01
	timer := timer.NewTimer()
	currentTime := timer.Seconds()
	for {
		newTime := timer.Seconds()
		frameTime = newTime - currentTime
		currentTime = newTime

		accumulator += frameTime
		for accumulator >= dt {
			// Individual entities
			w.entityManager.Update(w, dt)
			accumulator -= (dt)
			gameTime += dt
		}

		// world AI
		w.director.Update(w, dt, gameTime)

		// tell the event manager to dispatch all queued events
		events.DispatchEvents()
	}

}
