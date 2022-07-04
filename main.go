package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"atomicgo.dev/cursor"
	"github.com/daniel-munoz/life/model"
)

func main() {
	var (
		w   model.World
		err error
	)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	if len(os.Args) > 1 {
		w, err = model.ReadWorld(&os.Args[1])
		if err != nil {
			fmt.Printf("Error reading sample: %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		w, _ = model.ReadWorld(nil)
	}

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
