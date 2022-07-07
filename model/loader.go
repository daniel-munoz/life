package model

import (
	"bufio"
	"fmt"
	"os"
)

func ReadWorld(sampleName string) (World, error) {
	var (
		x, y    int64
		scanner *bufio.Scanner
	)
	newWorld := NewWorld()

	filename := fmt.Sprintf("./samples/%s.life", sampleName)
	f, err := os.Open(filename)
	if err != nil {
		return newWorld, err
	}
	defer f.Close()

	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		x = 0
		for _, c := range line {
			if c != ' ' {
				newWorld.AddCellIn(x, y, 0)
			}
			x++
		}
		y++
	}
	return newWorld, nil
}
