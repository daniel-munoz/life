package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"atomicgo.dev/cursor"
	"github.com/daniel-munoz/life/keyboard"
	"github.com/daniel-munoz/life/model"
)

func main() {
	var (
		w   model.World
		err error
	)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGTERM)

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

	area := cursor.NewArea()
	go display(&area, w, -10, -10, 40, 80)

	s := <-sigs
	switch s {
	case syscall.SIGTERM:
		area.Update("I've been told to stop")
	}
}

func display(area *cursor.Area, w model.World, top, left, bottom, right int64) {
	listener := keyboard.NewListener()
	listener.Start()
	helpTimer := 8
	for {
		time.Sleep(250 * time.Millisecond)
		if helpTimer > 0 {
			area.Update(`Keys:
  Up   : moves window 1 space up       Down : moves window 1 space down
  Left : moves window 1 space left     Right: moves window 1 space right
  I    : moves window 10 spaces up     K    : moves window 10 spaces down
  J    : moves window 10 spaces left   L    : moves window 10 spaces right
  Q    : ends the program              H    : displays this help
  `)
			helpTimer--
			continue
		}
		w.Evolve()

		topLeft, bottomRight := model.NewIndex(left,top), model.NewIndex(right,bottom)
		area.Update(w.WindowContent(topLeft, bottomRight))

		check := listener.Check()
		switch(check) {
		case keyboard.Stop:
			area.Update("Time to stop")
			os.Exit(0)
		case keyboard.Up:
			top--
			bottom--
		case keyboard.Down:
			top++
			bottom++
		case keyboard.Left:
			left--
			right--
		case keyboard.Right:
			left++
			right++
		case keyboard.PageUp:
			top-=10
			bottom-=10
		case keyboard.PageDown:
			top+=10
			bottom+=10
		case keyboard.PageLeft:
			left-=10
			right-=10
		case keyboard.PageRight:
			left+=10
			right+=10
		case keyboard.Help:
			helpTimer = 8
		default:
		}
	}
}
