package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atomicgo/cursor"
	"github.com/daniel-munoz/life/model"
)

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	/*
	w := model.NewWorld()
	w.AddCellIn(0,1,0)
	w.AddCellIn(1,2,0)
	w.AddCellIn(2,0,0)
	w.AddCellIn(2,1,0)
	w.AddCellIn(2,2,0)
	*/

	w := model.ReadWorld()

	cursor.Hide()
	defer cursor.Show()

	go display(w, -10, -10, 60, 100)

	<-sigs
	fmt.Println("Bye")
}

func display(w model.World, top, left, bottom, right int64) {
	topLeft, bottomRight := model.NewIndex(left,top), model.NewIndex(right,bottom)
	w.PrintWindow(topLeft, bottomRight)
	for true {
		time.Sleep(250 * time.Millisecond)
		w.Evolve()
		cursor.Up(int(bottom - top + 2))
		w.PrintWindow(topLeft, bottomRight)
	}
}
