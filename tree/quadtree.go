package tree

import (
	"fmt"
)

const (
	MAX_OBJECTS = 12
	MAX_LEVELS  = 5
)

const (
	TOP_LEFT = iota
	TOP_RIGHT
	BOTTOM_LEFT
	BOTTOM_RIGHT
)

type QuadTree struct {
	bounds     *Rectangle
	rectangles []*Rectangle
	nodes      map[int]*QuadTree
	level      int
}

func NewQuadTree(pLevel int, b *Rectangle) *QuadTree {
	return &QuadTree{
		level:      pLevel,
		bounds:     b,
		rectangles: make([]*Rectangle, 0),
		nodes:      make(map[int]*QuadTree, 4),
	}
}

/*
 * Insert the object into the QuadTree. If the node
 * exceeds the capacity, it will split and add all
 * objects to their corresponding nodes.
 */
func (q *QuadTree) Insert(rectangle *Rectangle) {
	// There are sub-nodes so push the object down to one of them if it
	// can fit in any of them
	if len(q.nodes) != 0 {
		index := q.index(rectangle)
		if index != -1 {
			q.nodes[index].Insert(rectangle)
			return
		}
	}
	q.rectangles = append(q.rectangles, rectangle)

	// We hit the objects limit and are still allowed to create more nodes,
	// split and push objects into sub-nodes
	if len(q.rectangles) > MAX_OBJECTS && q.level < MAX_LEVELS {
		// there are no sub nodes in this node
		if len(q.nodes) == 0 {
			q.split()
		}
		for i := len(q.rectangles) - 1; i >= 0; i-- {
			index := q.index(q.rectangles[i])
			if index != -1 {
				q.nodes[index].Insert(q.rectangles[i])
				q.rectangles = append(q.rectangles[:i], q.rectangles[i+1:]...)
			}
		}
	}
}

// Retrieve will find all intersecting rectangles
func (q *QuadTree) Retrieve(rect *Rectangle, result *[]*Rectangle) {
	*result = append(*result, q.rectangles...)
	for _, node := range q.nodes {
		if node.bounds.Intersects(rect) {
			node.Retrieve(rect, result)
		}
	}
}

func (q *QuadTree) Clear() {
	q.rectangles = make([]*Rectangle, 0)
	for _, node := range q.nodes {
		node.Clear()
	}
	q.nodes = make(map[int]*QuadTree, 4)
}

func (q *QuadTree) Boundaries() []*Rectangle {
	boundaries := make([]*Rectangle, 0)
	boundaries = append(boundaries, q.bounds)
	for _, node := range q.nodes {
		boundaries = append(boundaries, node.Boundaries()...)
	}
	return boundaries
}

func (q *QuadTree) Debug() {
	for _, obj := range q.rectangles {
		for i := 0; i < q.level; i++ {
			fmt.Print("\t")
		}
		fmt.Printf("object x:%d y:%d w:%d h:%d\n", int(obj.position.X), int(obj.position.Y), int(obj.halfWidth*2), int(obj.halfHeight*2))
	}
	for index, node := range q.nodes {
		for i := 0; i < node.level; i++ {
			fmt.Print("\t")
		}
		fmt.Printf("node %d.%d\n", node.level, index)
		node.Debug()
	}
}

/*
 * index determines which node the object belongs to. nil means
 * object cannot completely fit within a child node and is part
 * of the parent node
 */
func (q *QuadTree) index(r *Rectangle) int {
	// Object can completely fit within the left quadrants
	left := (r.position.X < q.bounds.position.X) && (r.position.X+r.halfWidth < q.bounds.position.X)
	// Object can completely fit within the top quadrants
	top := (r.position.Y < q.bounds.position.Y) && (r.position.Y+r.halfHeight < q.bounds.position.Y)
	if top && left {
		return TOP_LEFT
	}
	// Object can completely fit within the right quadrants
	right := (r.position.X > q.bounds.position.X) && (r.position.X-r.halfWidth > q.bounds.position.X)
	if top && right {
		return TOP_RIGHT
	}
	// Object can completely fit within the bottom quadrants
	bottom := (r.position.Y > q.bounds.position.Y) && (r.position.Y-r.halfHeight > q.bounds.position.Y)
	if bottom && left {
		return BOTTOM_LEFT
	}
	if bottom && right {
		return BOTTOM_RIGHT
	}
	// object can't fit in any of the sub nodes
	return -1
}

// split creates four sub-nodes for this node
func (q *QuadTree) split() {
	subWidth := q.bounds.halfWidth / 2
	subHeight := q.bounds.halfHeight / 2
	x := q.bounds.position.X
	y := q.bounds.position.Y
	q.nodes[TOP_LEFT] = NewQuadTree(q.level+1, NewRectangle(x-subWidth, y-subHeight, subWidth, subHeight))
	q.nodes[TOP_RIGHT] = NewQuadTree(q.level+1, NewRectangle(x+subWidth, y-subHeight, subWidth, subHeight))
	q.nodes[BOTTOM_LEFT] = NewQuadTree(q.level+1, NewRectangle(x-subWidth, y+subHeight, subWidth, subHeight))
	q.nodes[BOTTOM_RIGHT] = NewQuadTree(q.level+1, NewRectangle(x+subWidth, y+subHeight, subWidth, subHeight))
}
