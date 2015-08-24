package worldmap_test

import (
	. "github.com/stojg/pants/worldmap"
	"testing"
)

func TestNew(t *testing.T) {

	w := NewWorldMap(10, 10)
	w.Base().Place(1, 1, 1)
	w.Base().Place(1, 2, 2)
	w.Base().Place(2, 1, 3)
	w.Base().Place(2, 2, 4)

	if w.Base().Get(1, 1) != 1 {
		t.Errorf("test failed")
	}

	if w.Base().Get(1, 2) != 2 {
		t.Errorf("test failed")
	}

	if w.Base().Get(2, 1) != 3 {
		t.Errorf("test failed")
	}

	if w.Base().Get(2, 2) != 4 {
		t.Errorf("test failed")
	}
}

func TestCompress(t *testing.T) {
	w := NewWorldMap(10, 10)
	w.Base().Place(1, 1, 1)
	w.Base().Place(1, 2, 2)
	w.Base().Place(2, 1, 3)
	w.Base().Place(2, 2, 4)

	layers := w.Compress()

	for i := 0; i < len(layers); i++ {
		if len(layers[i]) != 10*10 {
			t.Errorf("Expected a lenght of 100, not %d", len(layers[i]))
		}
	}
}
