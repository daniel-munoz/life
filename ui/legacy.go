package ui

import (
	"time"

	"atomicgo.dev/cursor"
	"github.com/daniel-munoz/life/event"
	"github.com/daniel-munoz/life/model"
)

const (
	options string = `Keys:
  Up   : moves window 1 space up       Down : moves window 1 space down
  Left : moves window 1 space left     Right: moves window 1 space right
  I    : moves window 10 spaces up     K    : moves window 10 spaces down
  J    : moves window 10 spaces left   L    : moves window 10 spaces right
  Q    : ends the program              H    : displays this help
`
)

// Action is a function that updates the status of the world.
type Action func()

// Show displays the world in a terminal window.
func Show(w model.World, top, left, bottom, right int64) {
	stopChannel := make(chan struct{})

	cursor.Hide()
	defer cursor.Show()

	display := NewDisplay()

	listener := event.NewListener()
	listener.Start()

	gameView := NewGameView(top, left, bottom, right, stopChannel)

	go func() {
		for {
			if gameView.ShowHelp() {
				display.UpdateAndLock(options, 4500*time.Millisecond)
				gameView.ToggleHelp()
			}
			if !gameView.IsPaused() {
				w.Evolve()
			}

			topLeft, bottomRight := gameView.TopLeft(), gameView.BottomRight()
			display.UpdateAndLock(w.WindowContent(topLeft, bottomRight), 200*time.Millisecond)

			check := listener.Check()
			gameView.Execute(check)
			if gameView.Ended() {
				return
			}
		}
	}()
	for {
		select {
		case <-stopChannel:
			display.UpdateAndClose("Time to stop")
			return
		}
	}
}
