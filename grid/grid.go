package grid

import (
	. "github.com/stojg/pants/vector"
)

type Node struct {
	X        int
	Y        int
	TileType byte
}

type Grid struct {
	size   int
	Width  int
	Height int
	Nodes  []*Node
}

func (g *Grid) Place(x, y int, tileType byte) {
	g.Nodes[x+(y*g.Width)] = &Node{X: x, Y: y, TileType: tileType}
}

func (g *Grid) Get(x, y int) *Node {
	if x < 0 || x >= g.Width {
		return nil
	}
	if y < 0 || y >= g.Height {
		return nil
	}
	return g.Nodes[x+(y*g.Width)]
}

func NewGrid(width, height int) *Grid {

	wm := &Grid{
		size:   20,
		Nodes:  make([]*Node, width*height),
		Width:  width,
		Height: height,
	}

	i := 0
	for x := 0; x < width; x++ {
		for y := 0; y < width; y++ {
			wm.Place(x, y, 0)
			i += 1
		}
	}
	return wm
}

func (g *Grid) Neighbours(node *Node) []*Node {
	var dirs [8][2]int = [8][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, 1}, {-1, -1}, {1, -1}, {-1, 1}}
	tiles := make([]*Node, 0)
	for _, dir := range dirs {
		tile := g.Get(node.X+dir[0], node.X-dir[1])
		if tile != nil {
			tiles = append(tiles, tile)
		}
	}
	return tiles
}

func (g *Grid) Compress() []*Node {
	layers := make([]*Node, len(g.Nodes), len(g.Nodes))
	copy(layers, g.Nodes)
	return layers
}

func (g *Grid) TileFromWorld(position *Vec2) *Node {
	x := int(position.X) / g.size
	y := int(position.Y) / g.size
	return g.Get(x, y)
}
