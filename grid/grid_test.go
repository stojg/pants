package grid_test

import (
	. "github.com/stojg/pants/grid"
	. "github.com/stojg/pants/vector"
	"testing"
)

func TestNew(t *testing.T) {

	w := NewGrid(10, 10)
	w.Place(1, 1, 1)
	w.Place(1, 2, 2)
	w.Place(2, 1, 3)
	w.Place(2, 2, 4)

	if w.Get(1, 1).TileType != 1 {
		t.Errorf("test failed")
	}

	if w.Get(1, 2).TileType != 2 {
		t.Errorf("test failed")
	}

	if w.Get(2, 1).TileType != 3 {
		t.Errorf("test failed")
	}

	if w.Get(2, 2).TileType != 4 {
		t.Errorf("test failed")
	}
}

func TestCompress(t *testing.T) {
	w := NewGrid(10, 10)
	w.Place(1, 1, 1)
	w.Place(1, 2, 2)
	w.Place(2, 1, 3)
	w.Place(2, 2, 4)
	layers := w.Compress()
	if len(layers) != 10*10 {
		t.Errorf("Expected a length of 100, not %d", len(layers))
	}
}

func TestTranslateFromWorld(t *testing.T) {
	w := NewGrid(10, 10)
	w.Place(1, 1, 1)

	tile := w.TileFromWorld(&Vec2{20, 20})
	if tile == nil {
		t.Errorf("Did not get a tile")
		return
	}

	w.Place(1, 2, 2)
	tile = w.TileFromWorld(&Vec2{20, 40})
	if tile.TileType != 2 {
		t.Errorf("Did not get expected tile")
	}

	w.Place(2, 1, 3)
	tile = w.TileFromWorld(&Vec2{40, 20})
	if tile.TileType != 3 {
		t.Errorf("Did not get expected tile")
	}

	w.Place(2, 2, 4)
	tile = w.TileFromWorld(&Vec2{40, 40})
	if tile.TileType != 4 {
		t.Errorf("Did not get expected tile")
	}
}

func TestNeighbours(t *testing.T) {
	w := NewGrid(10, 10)

	centerTile := w.Get(5, 5)

	tiles := w.Neighbours(centerTile)

	if len(tiles) != 8 {
		t.Errorf("Expected that the center tile have 8 neighbours, got %d", len(tiles))
	}

	for _, tile := range tiles {
		if tile == centerTile {
			t.Errorf("tile cant itself be in a neighbours ")
		}
	}

}
