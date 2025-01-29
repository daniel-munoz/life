package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/daniel-munoz/life/types"
)

type index struct {
	x, y int64
}

func (i index) X() int64 {
	return i.x
}

func (i index) Y() int64 {
	return i.y
}

type Cell struct {
	birthTurn int64
}

type World struct {
	cells                map[index]*Cell
	topLeft, bottomRight index
	turn                 int64
	changes              int
	start                time.Time
}

func newCell(turn int64) *Cell {
	return &Cell{birthTurn: turn}
}

func NewIndex(x, y int64) types.Index {
	return index{x: x, y: y}
}

func NewWorld() *World {
	return &World{
		cells:       make(map[index]*Cell),
		topLeft:     index{0, 0},
		bottomRight: index{0, 0},
		turn:        0,
		start:       time.Now(),
	}
}

func (w World) GetCellIn(x, y int64) *Cell {
	return w.cells[index{x: x, y: y}]
}

type ChangeType int

const (
	BIRTH ChangeType = iota
	DEATH
)

type Change struct {
	turn   int64
	reason ChangeType
}

func (w *World) ApplyChanges(changes map[index]Change) {
	for location, c := range changes {
		switch c.reason {
		case BIRTH:
			w.cells[location] = newCell(c.turn)
		case DEATH:
			delete(w.cells, location)
		}
	}
	w.recalculateBorders()
}

func (w *World) AddCellIn(x, y, turn int64) {
	w.cells[index{x, y}] = newCell(turn)
	w.recalculateBorders()
}

func (w World) Print() {
	w.PrintWindow(w.topLeft, w.bottomRight)
}

func (w World) PrintWindow(topLeft, bottomRight types.Index) {
	fmt.Print(w.WindowContent(topLeft, bottomRight))
}

func (w World) WindowContent(topLeft, bottomRight types.Index) string {
	buffer := &strings.Builder{}
	fmt.Fprintf(buffer, "Turn: %d  Live Cells: %d  Limits: (%d,%d) -> (%d, %d) Changes: %d Age: %s    \n",
		w.turn,
		len(w.cells),
		w.topLeft.x,
		w.topLeft.y,
		w.bottomRight.x,
		w.bottomRight.y,
		w.changes,
		time.Since(w.start))
	x, y := topLeft.X(), topLeft.Y()
	for y <= bottomRight.Y() {
		for x <= bottomRight.X() {
			c := w.GetCellIn(x, y)
			if c != nil {
				fmt.Fprint(buffer, "x")
			} else {
				fmt.Fprint(buffer, " ")
			}
			x++
		}
		x = topLeft.X()
		y++
		fmt.Fprintln(buffer)
	}
	return buffer.String()
}

func (w *World) recalculateBorders() {
	var minX, maxX, minY, maxY int64
	first := true
	for location := range w.cells {
		if first {
			minX, minY, maxX, maxY = location.x, location.y, location.x, location.y
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
	w.topLeft = index{x: minX, y: minY}
	w.bottomRight = index{x: maxX, y: maxY}
}

func (w World) countNeighborsOf(location index, cache map[index]int, offset int) int {
	count, found := cache[location]
	if found {
		return count
	}
	x := location.x - 1
	for x <= location.x+1 {
		y := location.y - 1
		for y <= location.y+1 {
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

func (w World) analize(location index, turn int64, cache map[index]int, changes map[index]Change) {
	_, cellHasChange := changes[location]
	if cellHasChange {
		return
	}

	cell := w.GetCellIn(location.x, location.y)
	offset := 0
	if cell != nil {
		offset = 1
	}
	c := w.countNeighborsOf(location, cache, offset)

	if c == 3 && cell == nil {
		changes[location] = Change{turn: turn, reason: BIRTH}
	}
	if (c < 2 || c > 3) && cell != nil {
		changes[location] = Change{turn: turn, reason: DEATH}
	}
}

func (w World) analizeNeighborsOf(location index, turn int64, cache map[index]int, changes map[index]Change) {
	x := location.x - 1
	for x <= location.x+1 {
		y := location.y - 1
		for y <= location.y+1 {
			w.analize(index{x: x, y: y}, turn, cache, changes)
			y++
		}
		x++
	}
}

func (w *World) Evolve() {
	countCache := make(map[index]int)
	changes := make(map[index]Change)
	for cellLocation := range w.cells {
		w.analizeNeighborsOf(cellLocation, w.turn+1, countCache, changes)
	}
	w.turn++
	w.changes = len(changes)
	w.ApplyChanges(changes)
}
