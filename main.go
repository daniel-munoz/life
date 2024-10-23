package main

import (
	"fmt"
	"os"

	"github.com/daniel-munoz/life/model"
	"github.com/daniel-munoz/life/types"
	"github.com/daniel-munoz/life/ui"
)

func main() {
	var (
		w   types.World
		err error
		sampleName = "gliders"
	)

	// check if reading from a pipe, which does not work now
	inStat, _ := os.Stdin.Stat()
	if (inStat.Mode() & os.ModeCharDevice) != os.ModeCharDevice {
		fmt.Println("Sorry, using a redirected input is not supported")
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		sampleName = os.Args[1]
	}

	w, err = model.ReadWorld(sampleName)
	if err != nil {
		fmt.Printf("Error reading sample: %s\n", err.Error())
		os.Exit(1)
	}

	ui.Show(w, -10, -10, 40, 80)
}
