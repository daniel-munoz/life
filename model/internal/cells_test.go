package internal

import (
	"strings"
	"testing"
)

func TestNewIndex(t *testing.T) {
	tests := []struct {
		name  string
		x, y  int64
		wantX int64
		wantY int64
	}{
		{
			name:  "positive coordinates",
			x:     5,
			y:     10,
			wantX: 5,
			wantY: 10,
		},
		{
			name:  "negative coordinates",
			x:     -3,
			y:     -7,
			wantX: -3,
			wantY: -7,
		},
		{
			name:  "zero coordinates",
			x:     0,
			y:     0,
			wantX: 0,
			wantY: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := NewIndex(tt.x, tt.y)
			if idx.X() != tt.wantX {
				t.Errorf("NewIndex(%d, %d).X() = %d, want %d", tt.x, tt.y, idx.X(), tt.wantX)
			}
			if idx.Y() != tt.wantY {
				t.Errorf("NewIndex(%d, %d).Y() = %d, want %d", tt.x, tt.y, idx.Y(), tt.wantY)
			}
		})
	}
}

func TestWorld_AddCellIn(t *testing.T) {
	tests := []struct {
		name string
		x, y int64
	}{
		{
			name: "add cell at origin",
			x:    0,
			y:    0,
		},
		{
			name: "add cell in positive quadrant",
			x:    5,
			y:    5,
		},
		{
			name: "add cell in negative quadrant",
			x:    -3,
			y:    -3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorld()
			w.AddCellIn(tt.x, tt.y, 0)

			if cell := w.GetCellIn(tt.x, tt.y); cell == nil {
				t.Errorf("AddCellIn(%d, %d, 0) failed to add cell", tt.x, tt.y)
			}
		})
	}
}

func TestWorld_Evolve(t *testing.T) {
	tests := []struct {
		name          string
		initialCells  [][2]int64
		wantSurvivors [][2]int64
	}{
		{
			name: "single cell dies",
			initialCells: [][2]int64{
				{0, 0},
			},
			wantSurvivors: nil,
		},
		{
			name: "stable block",
			initialCells: [][2]int64{
				{0, 0}, {1, 0},
				{0, 1}, {1, 1},
			},
			wantSurvivors: [][2]int64{
				{0, 0}, {1, 0},
				{0, 1}, {1, 1},
			},
		},
		{
			name: "blinker - phase 1",
			initialCells: [][2]int64{
				{0, 0}, {1, 0}, {2, 0},
			},
			wantSurvivors: [][2]int64{
				{1, -1}, {1, 0}, {1, 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorld()

			// Add initial cells
			for _, cell := range tt.initialCells {
				w.AddCellIn(cell[0], cell[1], 0)
			}

			// Evolve world
			w.Evolve()

			// Check survivors
			if tt.wantSurvivors == nil {
				if len(w.cells) != 0 {
					t.Errorf("Expected no survivors, got %d cells", len(w.cells))
				}
				return
			}

			// Check each expected survivor
			for _, want := range tt.wantSurvivors {
				if cell := w.GetCellIn(want[0], want[1]); cell == nil {
					t.Errorf("Expected live cell at (%d, %d), got none", want[0], want[1])
				}
			}

			// Check total number of survivors
			if len(w.cells) != len(tt.wantSurvivors) {
				t.Errorf("Got %d survivors, want %d", len(w.cells), len(tt.wantSurvivors))
			}
		})
	}
}

func TestWorld_WindowContent(t *testing.T) {
	w := NewWorld()

	// Create a simple pattern (block)
	w.AddCellIn(0, 0, 0)
	w.AddCellIn(1, 0, 0)
	w.AddCellIn(0, 1, 0)
	w.AddCellIn(1, 1, 0)

	// Test window content
	topLeft := NewIndex(-1, -1)
	bottomRight := NewIndex(2, 2)
	content := w.WindowContent(topLeft, bottomRight)

	// Check if the content contains the expected pattern
	expected := []string{
		"    ",
		" xx ",
		" xx ",
		"    ",
	}

	// Split the content into lines and check each line
	lines := strings.Split(content, "\n")
	// Skip the first line as it contains turn information
	lines = lines[1 : len(lines)-1]

	if len(lines) != len(expected) {
		t.Errorf("WindowContent returned %d lines, want %d", len(lines), len(expected))
		return
	}

	for i, want := range expected {
		if !strings.HasSuffix(lines[i], want) {
			t.Errorf("line %d = %q, want suffix %q", i, lines[i], want)
		}
	}
}
