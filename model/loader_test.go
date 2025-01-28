package model

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadWorld(t *testing.T) {
	// Create a temporary test pattern
	pattern := []byte("x x\n xx\nx x")
	err := os.MkdirAll("samples", 0755)
	if err != nil {
		t.Fatalf("Failed to create samples directory: %v", err)
	}
	defer os.RemoveAll("samples")

	testFile := filepath.Join("samples", "test.life")
	err = os.WriteFile(testFile, pattern, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test successful read
	t.Run("successful pattern load", func(t *testing.T) {
		world, err := ReadWorld("test")
		if err != nil {
			t.Fatalf("ReadWorld() error = %v", err)
		}

		// Check pattern was loaded correctly
		expectedCells := [][2]int64{
			{0, 0}, {2, 0}, // First row
			{1, 1}, {2, 1}, // Second row
			{0, 2}, {2, 2}, // Third row
		}

		for _, cell := range expectedCells {
			topLeft := NewIndex(cell[0], cell[1])
			bottomRight := NewIndex(cell[0], cell[1])
			content := world.WindowContent(topLeft, bottomRight)
			if content[len(content)-2] != 'x' { // Account for newline
				t.Errorf("Expected live cell at (%d, %d)", cell[0], cell[1])
			}
		}
	})

	// Test error handling
	t.Run("non-existent file", func(t *testing.T) {
		_, err := ReadWorld("non-existent")
		if err == nil {
			t.Error("ReadWorld() expected error for non-existent file")
		}
	})
}
