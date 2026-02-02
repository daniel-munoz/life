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
		{
			name: "underpopulation - two cells die",
			initialCells: [][2]int64{
				{0, 0}, {1, 0},
			},
			wantSurvivors: nil,
		},
		{
			name: "overpopulation - center cell dies",
			initialCells: [][2]int64{
				{0, 0}, {1, 0}, {2, 0},
				{1, 1},
				{1, -1},
			},
			wantSurvivors: [][2]int64{
				{0, -1}, {1, -1}, {2, -1},
				{0, 0}, {2, 0},
				{0, 1}, {1, 1}, {2, 1},
			},
		},
		{
			name: "reproduction - new cell born",
			initialCells: [][2]int64{
				{0, 0}, {1, 0}, {0, 1},
			},
			wantSurvivors: [][2]int64{
				{0, 0}, {1, 0}, {0, 1}, {1, 1},
			},
		},
		{
			name: "glider - phase 1",
			initialCells: [][2]int64{
				{1, 0},
				{2, 1},
				{0, 2}, {1, 2}, {2, 2},
			},
			wantSurvivors: [][2]int64{
				{0, 1}, {2, 1},
				{1, 2}, {2, 2},
				{1, 3},
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

			// Verify no unexpected survivors
			for loc := range w.cells {
				found := false
				for _, want := range tt.wantSurvivors {
					if loc.x == want[0] && loc.y == want[1] {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Unexpected survivor at (%d, %d)", loc.x, loc.y)
				}
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

func TestWorld_ApplyChanges(t *testing.T) {
	tests := []struct {
		name            string
		initialCells    [][2]int64
		changes         map[index]Change
		wantCells       [][2]int64
		wantTopLeft     index
		wantBottomRight index
	}{
		{
			name:         "add new cells",
			initialCells: [][2]int64{{0, 0}},
			changes: map[index]Change{
				{1, 1}: {turn: 1, reason: BIRTH},
				{2, 2}: {turn: 1, reason: BIRTH},
			},
			wantCells:       [][2]int64{{0, 0}, {1, 1}, {2, 2}},
			wantTopLeft:     index{0, 0},
			wantBottomRight: index{2, 2},
		},
		{
			name:         "remove cells",
			initialCells: [][2]int64{{0, 0}, {1, 1}, {2, 2}},
			changes: map[index]Change{
				{1, 1}: {turn: 1, reason: DEATH},
			},
			wantCells:       [][2]int64{{0, 0}, {2, 2}},
			wantTopLeft:     index{0, 0},
			wantBottomRight: index{2, 2},
		},
		{
			name:         "mixed changes",
			initialCells: [][2]int64{{0, 0}, {1, 1}},
			changes: map[index]Change{
				{1, 1}:   {turn: 1, reason: DEATH},
				{-1, -1}: {turn: 1, reason: BIRTH},
				{2, 2}:   {turn: 1, reason: BIRTH},
			},
			wantCells:       [][2]int64{{-1, -1}, {0, 0}, {2, 2}},
			wantTopLeft:     index{-1, -1},
			wantBottomRight: index{2, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorld()

			// Set up initial state
			for _, cell := range tt.initialCells {
				w.AddCellIn(cell[0], cell[1], 0)
			}

			// Apply changes
			w.ApplyChanges(tt.changes)

			// Verify cells
			for _, want := range tt.wantCells {
				if cell := w.GetCellIn(want[0], want[1]); cell == nil {
					t.Errorf("Expected live cell at (%d, %d), got none", want[0], want[1])
				}
			}

			// Verify total cell count
			if len(w.cells) != len(tt.wantCells) {
				t.Errorf("Got %d cells, want %d", len(w.cells), len(tt.wantCells))
			}

			// Verify borders
			if w.topLeft != tt.wantTopLeft {
				t.Errorf("topLeft = %v, want %v", w.topLeft, tt.wantTopLeft)
			}
			if w.bottomRight != tt.wantBottomRight {
				t.Errorf("bottomRight = %v, want %v", w.bottomRight, tt.wantBottomRight)
			}
		})
	}
}

func TestWorld_RecalculateBorders(t *testing.T) {
	tests := []struct {
		name            string
		cells           [][2]int64
		wantTopLeft     index
		wantBottomRight index
	}{
		{
			name:            "single cell",
			cells:           [][2]int64{{5, 5}},
			wantTopLeft:     index{5, 5},
			wantBottomRight: index{5, 5},
		},
		{
			name:            "multiple cells - positive quadrant",
			cells:           [][2]int64{{1, 1}, {3, 4}, {5, 2}},
			wantTopLeft:     index{1, 1},
			wantBottomRight: index{5, 4},
		},
		{
			name:            "multiple cells - mixed quadrants",
			cells:           [][2]int64{{-2, -2}, {3, 4}, {-1, 2}, {1, -3}},
			wantTopLeft:     index{-2, -3},
			wantBottomRight: index{3, 4},
		},
		{
			name:            "empty world",
			cells:           [][2]int64{},
			wantTopLeft:     index{0, 0},
			wantBottomRight: index{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorld()

			// Add cells
			for _, cell := range tt.cells {
				w.AddCellIn(cell[0], cell[1], 0)
			}

			// Force recalculation
			w.recalculateBorders()

			// Verify borders
			if w.topLeft != tt.wantTopLeft {
				t.Errorf("topLeft = %v, want %v", w.topLeft, tt.wantTopLeft)
			}
			if w.bottomRight != tt.wantBottomRight {
				t.Errorf("bottomRight = %v, want %v", w.bottomRight, tt.wantBottomRight)
			}
		})
	}
}

func TestWorld_CountNeighborsOf(t *testing.T) {
	tests := []struct {
		name   string
		cells  [][2]int64
		target index
		offset int
		want   int
	}{
		{
			name:   "no neighbors",
			cells:  [][2]int64{{5, 5}},
			target: index{0, 0},
			offset: 0,
			want:   0,
		},
		{
			name: "all neighbors",
			cells: [][2]int64{
				{0, 0}, {1, 0}, {2, 0},
				{0, 1}, {1, 1}, {2, 1},
				{0, 2}, {1, 2}, {2, 2},
			},
			target: index{1, 1},
			offset: 1, // Target cell is alive
			want:   8,
		},
		{
			name: "some neighbors",
			cells: [][2]int64{
				{0, 0}, {2, 0},
				{1, 1},
				{0, 2}, {2, 2},
			},
			target: index{1, 1},
			offset: 1,
			want:   4,
		},
		{
			name: "edge neighbors",
			cells: [][2]int64{
				{-1, -1}, {0, -1}, {1, -1},
				{-1, 0}, {1, 0},
				{-1, 1}, {0, 1}, {1, 1},
			},
			target: index{0, 0},
			offset: 0,
			want:   8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorld()
			cache := make(map[index]int)

			// Add cells
			for _, cell := range tt.cells {
				w.AddCellIn(cell[0], cell[1], 0)
			}

			got := w.countNeighborsOf(tt.target, cache, tt.offset)
			if got != tt.want {
				t.Errorf("countNeighborsOf(%v) = %d, want %d", tt.target, got, tt.want)
			}

			// Test cache functionality
			cached := w.countNeighborsOf(tt.target, cache, tt.offset)
			if cached != tt.want {
				t.Errorf("cached countNeighborsOf(%v) = %d, want %d", tt.target, cached, tt.want)
			}
		})
	}
}

func TestWorld_Analize(t *testing.T) {
	tests := []struct {
		name       string
		cells      [][2]int64
		target     index
		turn       int64
		wantChange bool
		wantReason ChangeType
	}{
		{
			name: "underpopulation - cell dies",
			cells: [][2]int64{
				{1, 1},
				{2, 1},
			},
			target:     index{1, 1},
			turn:       1,
			wantChange: true,
			wantReason: DEATH,
		},
		{
			name: "survival - 2 neighbors",
			cells: [][2]int64{
				{0, 0},
				{1, 0},
				{2, 0},
			},
			target:     index{1, 0},
			turn:       1,
			wantChange: false,
		},
		{
			name: "survival - 3 neighbors",
			cells: [][2]int64{
				{0, 0}, {1, 0},
				{2, 0}, {2, 1},
			},
			target:     index{1, 0},
			turn:       1,
			wantChange: false,
		},
		{
			name: "overpopulation - cell dies",
			cells: [][2]int64{
				{0, 0}, {1, 0}, {2, 0},
				{1, 1},
				{1, -1},
			},
			target:     index{1, 0},
			turn:       1,
			wantChange: true,
			wantReason: DEATH,
		},
		{
			name: "reproduction - cell born",
			cells: [][2]int64{
				{0, 0}, {1, 0},
				{0, 1},
			},
			target:     index{1, 1},
			turn:       1,
			wantChange: true,
			wantReason: BIRTH,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorld()
			cache := make(map[index]int)
			changes := make(map[index]Change)

			// Add cells
			for _, cell := range tt.cells {
				w.AddCellIn(cell[0], cell[1], 0)
			}

			// Analyze target cell
			w.analyze(tt.target, tt.turn, cache, changes)

			// Check if change was recorded
			change, hasChange := changes[tt.target]
			if hasChange != tt.wantChange {
				t.Errorf("analize() hasChange = %v, want %v", hasChange, tt.wantChange)
			}

			if tt.wantChange && change.reason != tt.wantReason {
				t.Errorf("analize() change reason = %v, want %v", change.reason, tt.wantReason)
			}
		})
	}
}

func TestWorld_AnalizeNeighborsOf(t *testing.T) {
	tests := []struct {
		name        string
		cells       [][2]int64
		target      index
		turn        int64
		wantChanges map[index]ChangeType
	}{
		{
			name: "isolated cell - all neighbors dead",
			cells: [][2]int64{
				{1, 1},
			},
			target: index{1, 1},
			turn:   1,
			wantChanges: map[index]ChangeType{
				{1, 1}: DEATH, // Center cell dies from underpopulation
			},
		},
		{
			name: "block formation",
			cells: [][2]int64{
				{0, 0}, {1, 0},
				{0, 1},
			},
			target: index{0, 0},
			turn:   1,
			wantChanges: map[index]ChangeType{
				{1, 1}: BIRTH, // New cell born to complete block
			},
		},
		{
			name: "complex pattern",
			cells: [][2]int64{
				{0, 0}, {1, 0}, {2, 0},
				{1, 1},
			},
			target: index{1, 0},
			turn:   1,
			wantChanges: map[index]ChangeType{
				{1, -1}: BIRTH, // New cell born here
				{0, 1}:  BIRTH, // New cell born above
				{2, 1}:  BIRTH, // New cell born below
			},
		},
		{
			name: "edge of world",
			cells: [][2]int64{
				{-1, -1}, {0, -1},
				{-1, 0},
			},
			target: index{-1, -1},
			turn:   1,
			wantChanges: map[index]ChangeType{
				{0, 0}: BIRTH, // New cell born to complete block
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorld()
			cache := make(map[index]int)
			changes := make(map[index]Change)

			// Add cells
			for _, cell := range tt.cells {
				w.AddCellIn(cell[0], cell[1], 0)
			}

			// Analyze neighbors
			w.analyzeNeighborsOf(tt.target, tt.turn, cache, changes)

			// Verify changes
			for loc, wantReason := range tt.wantChanges {
				change, hasChange := changes[loc]
				if !hasChange {
					t.Errorf("Expected change at %v, got none", loc)
					continue
				}
				if change.reason != wantReason {
					t.Errorf("At %v: got change reason %v, want %v", loc, change.reason, wantReason)
				}
				if change.turn != tt.turn {
					t.Errorf("At %v: got turn %d, want %d", loc, change.turn, tt.turn)
				}
			}

			// Verify no unexpected changes
			for loc, change := range changes {
				if _, expected := tt.wantChanges[loc]; !expected {
					t.Errorf("Unexpected change at %v: %v", loc, change)
				}
			}
		})
	}
}

func TestWorld_WindowContent_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		cells       [][2]int64
		topLeft     index
		bottomRight index
		wantLines   []string
	}{
		{
			name:        "empty world",
			cells:       [][2]int64{},
			topLeft:     index{-1, -1},
			bottomRight: index{1, 1},
			wantLines: []string{
				"   ",
				"   ",
				"   ",
			},
		},
		{
			name:        "single point window",
			cells:       [][2]int64{{0, 0}},
			topLeft:     index{0, 0},
			bottomRight: index{0, 0},
			wantLines: []string{
				"x",
			},
		},
		{
			name:        "window outside living cells",
			cells:       [][2]int64{{0, 0}},
			topLeft:     index{5, 5},
			bottomRight: index{7, 7},
			wantLines: []string{
				"   ",
				"   ",
				"   ",
			},
		},
		{
			name:        "negative coordinates",
			cells:       [][2]int64{{-2, -2}, {-1, -1}},
			topLeft:     index{-3, -3},
			bottomRight: index{-1, -1},
			wantLines: []string{
				"   ",
				" x ",
				"  x",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorld()

			// Add cells
			for _, cell := range tt.cells {
				w.AddCellIn(cell[0], cell[1], 0)
			}

			// Get window content
			content := w.WindowContent(tt.topLeft, tt.bottomRight)

			// Split content and skip first line (turn info)
			lines := strings.Split(content, "\n")
			lines = lines[1 : len(lines)-1] // Skip first line and empty last line

			// Verify number of lines
			if len(lines) != len(tt.wantLines) {
				t.Errorf("Got %d lines, want %d", len(lines), len(tt.wantLines))
				return
			}

			// Verify each line
			for i, want := range tt.wantLines {
				if !strings.HasSuffix(lines[i], want) {
					t.Errorf("line %d = %q, want suffix %q", i, lines[i], want)
				}
			}
		})
	}
}
