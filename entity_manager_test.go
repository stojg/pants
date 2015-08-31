package main

import (
	"github.com/stojg/pants/network"
	"testing"
	"math/rand"
)

func TestAddSprite(t *testing.T) {
	list := NewEntityManager()
	spriteID := list.NewEntity(0, 0, "sprite.png")
	s := network.NewServer("8080")
	w := NewWorld(list, s)
	list.Update(w, 0.016)
	if spriteID != 1 {
		t.Errorf("Could not get sprite id")
	}
}

// ~ 294222ns
func BenchmarkSpriteUpdate(b *testing.B) {
	b.StopTimer()
	list := NewEntityManager()
	rand.Seed(1)
	for i := 0; i < 100; i++ {
		list.NewEntity(rand.Float64()*1000, rand.Float64()*1000, "sprite.png")
	}
	s := network.NewServer("8080")
	w := NewWorld(list, s)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		list.Update(w, 0.0016)
	}
}

func TestUpdate(t *testing.T) {
	e := NewEntityManager()
	s := network.NewServer("8080")
	w := NewWorld(e, s)
	e.Update(w, 0.016)
}
