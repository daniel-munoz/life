package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListSamples(t *testing.T) {
	// Change to the project directory to find samples
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Find project root by looking for samples directory
	projectDir := findProjectRoot(t)
	if projectDir == "" {
		t.Skip("Could not find project root with samples directory")
	}

	err = os.Chdir(projectDir)
	if err != nil {
		t.Fatalf("Failed to change to project directory: %v", err)
	}

	samples, err := listSamples()
	if err != nil {
		t.Fatalf("listSamples() returned error: %v", err)
	}

	if len(samples) == 0 {
		t.Error("listSamples() returned empty list, expected at least one sample")
	}

	// Verify samples don't have .life extension
	for _, sample := range samples {
		if filepath.Ext(sample) == ".life" {
			t.Errorf("Sample %q should not have .life extension", sample)
		}
	}
}

func TestListSamples_NoSamples(t *testing.T) {
	// Create a temporary directory without samples
	tmpDir, err := os.MkdirTemp("", "life-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create empty samples directory
	samplesDir := filepath.Join(tmpDir, "samples")
	err = os.Mkdir(samplesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create samples directory: %v", err)
	}

	// Change to the temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	_, err = listSamples()
	if err == nil {
		t.Error("listSamples() should return error when no samples found")
	}
}

// findProjectRoot looks for the project root by finding the samples directory
func findProjectRoot(t *testing.T) string {
	t.Helper()

	// Start from current working directory
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Check if samples directory exists in current directory
	samplesPath := filepath.Join(wd, "samples")
	if info, err := os.Stat(samplesPath); err == nil && info.IsDir() {
		return wd
	}

	// Try parent directories (up to 5 levels)
	dir := wd
	for i := 0; i < 5; i++ {
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent

		samplesPath := filepath.Join(dir, "samples")
		if info, err := os.Stat(samplesPath); err == nil && info.IsDir() {
			return dir
		}
	}

	return ""
}
