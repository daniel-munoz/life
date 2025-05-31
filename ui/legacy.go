package ui

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"atomicgo.dev/cursor"
	"github.com/daniel-munoz/life/event"
	"github.com/daniel-munoz/life/types"
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

// resetTerminal forces a terminal reset using stty to restore normal input mode
func resetTerminal() {
	// Use stty to reset terminal to sane state
	cmd := exec.Command("stty", "sane")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run() // Ignore errors as this is best-effort cleanup
}

// Show displays the world in a terminal window.
func Show(w types.World, top, left, bottom, right int64) {
	stopChannel := make(chan struct{})
	
	// Set up signal handling for proper cleanup
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Ensure proper cleanup of terminal state
	defer func() {
		cursor.Show()
		resetTerminal()
	}()

	cursor.Hide()

	display := NewDisplay()
	defer display.Close()

	listener := event.NewListener()
	listener.Start()
	defer listener.Stop()

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
			cursor.Show()
			resetTerminal()
			return
		case sig := <-sigChannel:
			// Handle signals for proper cleanup
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP:
				display.UpdateAndClose("Program interrupted")
				cursor.Show()
				resetTerminal()
				return
			}
		}
	}
}
