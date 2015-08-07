package main

import (
	"log"
	"math/rand"
	"time"
)

const NetFPS = 50 * time.Millisecond

func NewWorld(list *EntityList) *World {
	current := time.Now()
	return &World{
		netTicked:  current,
		gameTicked: current,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
		//		rand:       rand.New(rand.NewSource(1)),
		entities: list,
		director: &Director{},
	}
}

type World struct {
	netTicked  time.Time
	netTick    uint64
	gameTicked time.Time
	gameTick   uint64
	entities   *EntityList
	rand       *rand.Rand
	director   *Director
}

type EntityUpdate struct {
	Id          uint64            `bson:",minsize"`
	X, Y        float64           `bson:",minsize,omitempty"`
	Orientation float64           `bson:",minsize"`
	Type        string            `bson:",minsize,omitempty"`
	Dead        bool              `bson:",minsize,omitempty"`
	Data        map[string]string `bson:",minsize,omitempty"`
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
		changedSprites := make([]*EntityUpdate, 0)
		for id := range list.updated {
			changedSprites = append(changedSprites, &EntityUpdate{
				Id:          id,
				X:           w.entities.physics[id].Position.X,
				Y:           w.entities.physics[id].Position.Y,
				Orientation: w.entities.physics[id].Orientation,
				Type:        w.entities.sprites[id].Type,
				Dead:        w.entities.sprites[id].Dead,
				Data: map[string]string{
					"Image": list.sprites[id].Image,
				},
			})
			delete(list.updated, id)
		}

		if len(changedSprites) > 0 {
			h.Send(&Message{
				Topic:     "update",
				Data:      changedSprites,
				Tick:      w.gameTick,
				Timestamp: float64(currentTime.UnixNano()) / 1000000,
			})
		}
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
			log.Printf("world lag: %d ms", int(frameTime*1000))
		}

		accumulator += frameTime
		for accumulator >= dt {
			// Individual entities
			w.entities.Update(w, dt, gameTime)
			accumulator -= (dt)
			gameTime += dt
		}

		// world AI
		w.director.Update(w, dt, gameTime)
	}

}
