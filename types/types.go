// Package types defines the core interfaces for the Game of Life simulation.
package types

// Index represents a 2D coordinate in the world grid.
type Index interface {
	X() int64
	Y() int64
}

// World represents the Game of Life universe and its operations.
type World interface {
	// AddCellIn adds a new cell at the specified coordinates.
	AddCellIn(x, y, turn int64)
	// Evolve advances the world by one generation.
	Evolve()
	// WindowContent returns a string representation of the world within the given bounds.
	WindowContent(topLeft, bottomRight Index) string
}
