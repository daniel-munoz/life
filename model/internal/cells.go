// Package internal provides the core implementation of Conway's Game of Life.
package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/daniel-munoz/life/types"
)

// Game of Life rules constants.
// A cell is born with exactly 3 neighbors, survives with 2-3 neighbors,
// and dies otherwise (underpopulation or overpopulation).
const (
	birthNeighborCount   = 3 // Cell is born with exactly 3 neighbors
	minSurvivalNeighbors = 2 // Cell dies with fewer than 2 neighbors
	maxSurvivalNeighbors = 3 // Cell dies with more than 3 neighbors
)

// index represents a 2D coordinate in the world grid.
type index struct {
	x, y int64
}

// X returns the x coordinate.
func (i index) X() int64 {
	return i.x
}

// Y returns the y coordinate.
func (i index) Y() int64 {
	return i.y
}

// Cell represents a living cell in the world with its birth turn recorded.
type Cell struct {
	birthTurn int64
}

// World represents the Game of Life universe containing all cells.
type World struct {
	cells                map[index]*Cell
	topLeft, bottomRight index
	turn                 int64
	changes              int
	start                time.Time
}

// newCell creates a new cell born at the specified turn.
func newCell(turn int64) *Cell {
	return &Cell{birthTurn: turn}
}

// NewIndex creates a new coordinate index.
func NewIndex(x, y int64) types.Index {
	return index{x: x, y: y}
}

// NewWorld creates an empty world ready for cells to be added.
func NewWorld() *World {
	return &World{
		cells:       make(map[index]*Cell),
		topLeft:     index{0, 0},
		bottomRight: index{0, 0},
		turn:        0,
		start:       time.Now(),
	}
}

// GetCellIn returns the cell at the specified coordinates, or nil if empty.
func (w World) GetCellIn(x, y int64) *Cell {
	return w.cells[index{x: x, y: y}]
}

// ChangeType indicates whether a cell is being born or dying.
type ChangeType int

// Change type constants.
const (
	BIRTH ChangeType = iota // A new cell is born
	DEATH                   // An existing cell dies
)

// Change represents a pending birth or death of a cell.
type Change struct {
	turn   int64
	reason ChangeType
}

// ApplyChanges applies all pending births and deaths to the world.
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

// AddCellIn adds a new cell at the specified coordinates.
func (w *World) AddCellIn(x, y, turn int64) {
	w.cells[index{x, y}] = newCell(turn)
	w.recalculateBorders()
}

// Print outputs the entire world to stdout.
func (w World) Print() {
	w.PrintWindow(w.topLeft, w.bottomRight)
}

// PrintWindow outputs the specified window of the world to stdout.
func (w World) PrintWindow(topLeft, bottomRight types.Index) {
	fmt.Print(w.WindowContent(topLeft, bottomRight))
}

// WindowContent returns a string representation of the world within the given bounds.
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

// recalculateBorders updates the world's bounding box based on current cells.
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

// countNeighborsOf counts living neighbors of a cell, using cache for efficiency.
// The offset is subtracted from the count (1 if the cell itself is alive, 0 otherwise).
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

// analize determines if a cell should be born or die based on Game of Life rules.
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

	if c == birthNeighborCount && cell == nil {
		changes[location] = Change{turn: turn, reason: BIRTH}
	}
	if (c < minSurvivalNeighbors || c > maxSurvivalNeighbors) && cell != nil {
		changes[location] = Change{turn: turn, reason: DEATH}
	}
}

// analizeNeighborsOf analyzes all 9 cells in the 3x3 grid centered on the given location.
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

// Evolve advances the world by one generation, applying Game of Life rules.
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
