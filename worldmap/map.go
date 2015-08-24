package worldmap

const size = 128

type Layer struct {
	width  int
	height int
	Tiles  []byte
}

func (l *Layer) Place(x, y int, tileType byte) {
	l.Tiles[x+(y*l.width)] = tileType
}

func (l Layer) Get(x, y int) byte {
	return l.Tiles[x+(y*l.width)]
}

type WorldMap struct {
	width  int
	height int
	base   *Layer
}

func NewWorldMap(width, height int) *WorldMap {
	base := &Layer{
		width:  width,
		height: height,
		Tiles:  make([]byte, width*height),
	}
	return &WorldMap{
		base: base,
	}
}

func (m *WorldMap) Compress() [][]byte {
	layers := make([][]byte, 1, 1)
	layers[0] = make([]byte, len(m.base.Tiles), len(m.base.Tiles))
	copy(layers[0], m.base.Tiles)
	return layers
}

func (m *WorldMap) Base() *Layer {
	return m.base
}
