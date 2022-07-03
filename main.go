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

	w := model.ReadWorld()

	cursor.Hide()
	defer cursor.Show()

	go display(w, -10, -10, 60, 100)

	s := <- sigs
	switch s {
	case syscall.SIGINT:
		fmt.Println("\nAs you wish")
	case syscall.SIGTERM:
		fmt.Println("I've been told to stop")
	}
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
