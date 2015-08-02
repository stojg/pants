package main

import (
	"testing"
)

func TestAddSprite(t *testing.T) {
	list := NewSpriteList()
	spriteID := list.NewSprite(0, 0, "sprite.png")
	w := NewWorld(list)
	list.Update(w, 0.016, 0.016)
	if spriteID != 1 {
		t.Errorf("Could not get sprite id")
	}
}

func BenchmarkSpriteUpdate(b *testing.B) {
	list := NewSpriteList()

	for i := 0; i < 100; i++ {
		list.NewSprite(0, 0, "sprite.png")
	}

	w := NewWorld(list)

	for n := 0; n < b.N; n++ {
		list.Update(w, 0.0016, 0.016)
	}
}

func BenchmarkSpriteChanged(b *testing.B) {
	list := NewSpriteList()
	list.NewSprite(0, 0, "sprite.png")
	for n := 0; n < b.N; n++ {
		list.Changed(false)
	}
}
