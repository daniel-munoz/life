package model

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadWorld(t *testing.T) {
	// Create samples directory
	err := os.MkdirAll("samples", 0755)
	if err != nil {
		t.Fatalf("Failed to create samples directory: %v", err)
	}
	defer os.RemoveAll("samples")

	tests := []struct {
		name       string
		pattern    string
		wantCells  [][2]int64
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "basic pattern",
			pattern: "x x\n xx\nx x",
			wantCells: [][2]int64{
				{0, 0}, {2, 0}, // First row
				{1, 1}, {2, 1}, // Second row
				{0, 2}, {2, 2}, // Third row
			},
		},
		{
			name:      "empty file",
			pattern:   "",
			wantCells: [][2]int64{},
		},
		{
			name:    "single cell",
			pattern: "x",
			wantCells: [][2]int64{
				{0, 0},
			},
		},
		{
			name:    "multiple lines with spaces",
			pattern: "  x  \n x x \n  x  ",
			wantCells: [][2]int64{
				{2, 0},
				{1, 1}, {3, 1},
				{2, 2},
			},
		},
		{
			name:    "irregular line lengths",
			pattern: "x x\nx x x\nx",
			wantCells: [][2]int64{
				{0, 0}, {2, 0},
				{0, 1}, {2, 1}, {4, 1},
				{0, 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join("samples", "test.life")
			err := os.WriteFile(testFile, []byte(tt.pattern), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Read world
			world, err := ReadWorld("test")
			if tt.wantErr {
				if err == nil {
					t.Error("ReadWorld() expected error")
				} else if err.Error() != tt.wantErrMsg {
					t.Errorf("ReadWorld() error = %v, want %v", err, tt.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("ReadWorld() unexpected error: %v", err)
			}

			// Check each expected cell
			for _, cell := range tt.wantCells {
				topLeft := NewIndex(cell[0], cell[1])
				bottomRight := NewIndex(cell[0], cell[1])
				content := world.WindowContent(topLeft, bottomRight)
				if content[len(content)-2] != 'x' { // Account for newline
					t.Errorf("Expected live cell at (%d, %d)", cell[0], cell[1])
				}
			}

			// Verify no unexpected cells in surrounding area
			for x := int64(-10); x <= 10; x++ {
				for y := int64(-10); y <= 10; y++ {
					topLeft := NewIndex(x, y)
					bottomRight := NewIndex(x, y)
					content := world.WindowContent(topLeft, bottomRight)
					hasCell := content[len(content)-2] == 'x'

					expected := false
					for _, want := range tt.wantCells {
						if x == want[0] && y == want[1] {
							expected = true
							break
						}
					}
					if hasCell && !expected {
						t.Errorf("Unexpected live cell at (%d, %d)", x, y)
					}
				}
			}
		})
	}

	// Test error cases
	t.Run("non-existent file", func(t *testing.T) {
		_, err := ReadWorld("non-existent")
		if err == nil {
			t.Error("ReadWorld() expected error for non-existent file")
		}
	})

	t.Run("directory instead of file", func(t *testing.T) {
		dirPath := filepath.Join("samples", "test_dir")
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
		defer os.RemoveAll(dirPath)

		_, err := ReadWorld("test_dir")
		if err == nil {
			t.Error("ReadWorld() expected error when reading directory")
		}
	})

	t.Run("file with no read permissions", func(t *testing.T) {
		testFile := filepath.Join("samples", "noperm.life")
		err := os.WriteFile(testFile, []byte("x"), 0000)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		_, err = ReadWorld("noperm")
		if err == nil {
			t.Error("ReadWorld() expected error for file with no permissions")
		}
	})
}
