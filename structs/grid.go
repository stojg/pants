package structs

func NewGrid(width, height, cellSize int) *Grid {

	return &Grid{
		Cells:    make([][]uint64, int(width/cellSize)*int(height/cellSize)),
		Width:    width,
		Height:   height,
		CellSize: cellSize,
		Cols:     int(width / cellSize),
		Rows:     int(height / cellSize),
	}
}

type Grid struct {
	Cells    [][]uint64
	Width    int
	Height   int
	CellSize int
	Rows     int
	Cols     int
}

func (g *Grid) Clear() {
	g.Cells = make([][]uint64, g.Rows*g.Cols)
}
