package model

import (
	"bufio"
	"os"
)

func ReadWorld() World {
	var x, y int64
	newWorld := NewWorld()

	scanner := bufio.NewScanner(os.Stdin)
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
	return newWorld
}
