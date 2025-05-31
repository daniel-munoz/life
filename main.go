package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/daniel-munoz/life/model"
	"github.com/daniel-munoz/life/types"
	"github.com/daniel-munoz/life/ui"
)

// listSamples lists all available samples in the samples directory
func listSamples() ([]string, error) {
	var samples []string

	// Get all .life files from samples directory
	files, err := filepath.Glob("./samples/*.life")
	if err != nil {
		return nil, err
	}

	// Extract sample names without extension
	for _, file := range files {
		base := filepath.Base(file)
		sampleName := strings.TrimSuffix(base, ".life")
		samples = append(samples, sampleName)
	}

	return samples, nil
}

// promptSampleSelection presents available samples and lets user select one
func promptSampleSelection() (string, error) {
	samples, err := listSamples()
	if err != nil {
		return "", fmt.Errorf("failed to list samples: %w", err)
	}

	// Display samples
	fmt.Println("Available samples:")
	for i, sample := range samples {
		fmt.Printf("%d. %s\n", i+1, sample)
	}

	// Prompt user for selection
	fmt.Print("\nEnter the number of the sample to display (or press Enter for default 'gliders'): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Trim newline and check if input is empty
	input = strings.TrimSpace(input)
	if input == "" {
		return "gliders", nil
	}

	// Convert input to number
	num, err := strconv.Atoi(input)
	if err != nil || num < 1 || num > len(samples) {
		return "", fmt.Errorf("invalid selection")
	}

	return samples[num-1], nil
}

func main() {
	var (
		w          types.World
		err        error
		sampleName string
	)

	// check if reading from a pipe, which does not work now
	inStat, _ := os.Stdin.Stat()
	if (inStat.Mode() & os.ModeCharDevice) != os.ModeCharDevice {
		fmt.Println("Sorry, using a redirected input is not supported")
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		// Use command line argument if provided
		sampleName = os.Args[1]
	} else {
		// Otherwise prompt user to select a sample
		sampleName, err = promptSampleSelection()
		if err != nil {
			fmt.Printf("Error selecting sample: %s\n", err.Error())
			os.Exit(1)
		}
	}

	fmt.Printf("Loading sample: %s\n", sampleName)
	w, err = model.ReadWorld(sampleName)
	if err != nil {
		fmt.Printf("Error reading sample: %s\n", err.Error())
		os.Exit(1)
	}

	ui.Show(w, -10, -10, 40, 80)
}
