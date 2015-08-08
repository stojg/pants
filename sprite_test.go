package main

import (
	"testing"
)

func TestAddSprite(t *testing.T) {
	list := NewEntityManager()
	spriteID := list.NewEntity(0, 0, "sprite.png")
	w := NewWorld(list)
	list.Update(w, 0.016, 0.016)
	if spriteID != 1 {
		t.Errorf("Could not get sprite id")
	}
}

func BenchmarkSpriteUpdate(b *testing.B) {
	list := NewEntityManager()

	for i := 0; i < 100; i++ {
		list.NewEntity(0, 0, "sprite.png")
	}

	w := NewWorld(list)

	for n := 0; n < b.N; n++ {
		list.Update(w, 0.0016, 0.016)
	}
}
