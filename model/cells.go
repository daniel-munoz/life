package model

import "fmt"

type Cell struct {
	birthTurn int64
}

type Index struct {
	x, y      int64
}

type World struct {
	cells                map[Index]*Cell
	topLeft, bottomRight Index
	turn                 int64
	changes              int
}

func newCell(turn int64) *Cell {
	return &Cell{birthTurn: turn}
}

func NewIndex(x, y int64) Index {
	return Index{x:x,y:y}
}

func NewWorld() World {
	return World{
		cells:       make(map[Index]*Cell),
		topLeft:     Index{0,0},
		bottomRight: Index{0,0},
		turn:        0,
	}
}

func (w World) GetCellIn(x, y int64) *Cell {
	return w.cells[Index{x:x,y:y}]
}

type ChangeType int

const (
	BIRTH ChangeType = iota
	DEATH
)

type Change struct {
	location Index
	turn     int64
	reason   ChangeType
}

func (w *World) ApplyChanges(changes []Change) {
	for _, c := range changes {
		switch c.reason {
			case BIRTH:
				w.cells[c.location] = newCell(c.turn)
			case DEATH:
				delete(w.cells, c.location)
		}
	}
	w.recalculateBorders()
}

func (w *World) AddCellIn(x, y, turn int64) {
	w.cells[Index{x,y}] = newCell(turn)
	w.recalculateBorders()
}

func (w World) Print() {
	w.PrintWindow(w.topLeft, w.bottomRight)
}

func (w World) PrintWindow(topLeft, bottomRight Index) {
	fmt.Printf("Turn: %d   Live Cells: %d   Limits: (%d,%d) -> (%d, %d)   Changes: %d     \n", w.turn, len(w.cells), w.topLeft.x, w.topLeft.y, w.bottomRight.x, w.bottomRight.y, w.changes)
	x, y := topLeft.x, topLeft.y
	for y <= bottomRight.y {
		for x <= bottomRight.x {
			c := w.GetCellIn(x, y)
			if c != nil {
				fmt.Print("x")
			} else {
				fmt.Print(" ")
			}
			x++
		}
		x = topLeft.x
		y++
		fmt.Println()
	}
}

func (w *World) recalculateBorders() {
    var minX, maxX, minY, maxY int64
    first := true
    for location,_ := range w.cells {
	    if first {
		    minX,minY,maxX,maxY = location.x,location.y,location.x,location.y
		    first = false
		    continue
	    }
	    if location.x < minX {
		    minX = location.x
	    }
	    if location.y < minY {
		    minY = location.y
	    }
	    if location.x > maxX {
		    maxX = location.x
	    }
	    if location.y > maxY {
		    maxY = location.y
	    }
    }
    w.topLeft = Index{x:minX,y:minY}
    w.bottomRight = Index{x:maxX,y:maxY}
}

func (w World) countNeighborsOf(location Index, cache map[Index]int, offset int) int {
	count, found := cache[location]
	if found {
		return count
	}
	x := location.x - 1
	for x <= location.x + 1 {
		y := location.y - 1
		for y <= location.y + 1 {
			if w.GetCellIn(x, y) != nil {
				count++
			}
			y++
		}
		x++
	}
	cache[location] = count - offset
	return count - offset
}

func (w World) analize(location Index, turn int64, cache map[Index]int) []Change {
	cell := w.GetCellIn(location.x, location.y) 
	offset := 0
	if cell != nil {
		offset = 1
	}
	c := w.countNeighborsOf(location, cache, offset)

	if c == 3 && cell == nil {
		return []Change{Change{location:location,turn:turn,reason:BIRTH}}
	}
	if (c <2 || c > 3) && cell != nil {
		return []Change{Change{location:location,turn:turn,reason:DEATH}}
	}
	return nil
}

func (w World) analizeNeighborsOf(location Index, turn int64, cache map[Index]int) (changes []Change) {
	x := location.x - 1
	for x <= location.x + 1 {
		y := location.y - 1
		for y <= location.y + 1 {
			changes = append(changes, w.analize(NewIndex(x,y), turn, cache)...)
			y++
		}
		x++
	}
	return
}

func (w *World) Evolve() {
	countCache := make(map[Index]int)
	changes := make([]Change, 0)
	for cellLocation, _ := range(w.cells) {
		changes = append(changes, w.analizeNeighborsOf(cellLocation, w.turn+1, countCache)...)
	}
	w.turn++
	w.changes = len(changes)
	w.ApplyChanges(changes)
}
