package event

import (
	"os"
	"os/signal"
	"syscall"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

// Listener provides an interface for capturing and processing keyboard events.
// It runs in the background and queues events for retrieval via Check().
type Listener interface {
	// Start begins listening for keyboard input and OS signals.
	Start()
	// Check returns the next queued event, or None if no event is available.
	Check() Event
	// Stop terminates the listener and releases resources.
	Stop()
}

// NewListener creates a new keyboard event listener.
func NewListener() Listener {
	return &gameListener{}
}

// gameListener implements the Listener interface using atomicgo/keyboard.
type gameListener struct {
	queue   chan Event
	running bool
}

// mapKeyToEvent maps a keyboard key to an Event and returns whether to stop listening.
func mapKeyToEvent(k keys.Key) (event Event, stop bool) {
	switch k.Code {
	case keys.Up:
		return Up, false
	case keys.Down:
		return Down, false
	case keys.Left:
		return Left, false
	case keys.Right:
		return Right, false
	case keys.CtrlC:
		return Stop, true
	case keys.Space:
		return Pause, false
	case keys.RuneKey:
		return mapRuneKeyToEvent(k.String())
	default:
		return None, false
	}
}

// mapRuneKeyToEvent maps a rune key string to an Event and returns whether to stop listening.
func mapRuneKeyToEvent(key string) (event Event, stop bool) {
	switch key {
	case "q":
		return Stop, true
	case "i":
		return PageUp, false
	case "k":
		return PageDown, false
	case "j":
		return PageLeft, false
	case "l":
		return PageRight, false
	case "h":
		return Help, false
	default:
		return None, false
	}
}

// startSignalHandler starts a goroutine to handle OS signals.
func (gl *gameListener) startSignalHandler(sigs chan os.Signal) {
	go func() {
		select {
		case s := <-sigs:
			switch s {
			case syscall.SIGTERM:
				keyboard.SimulateKeyPress(rune('q'))
				gl.queue <- Stop
			}
		}
	}()
}

func (gl *gameListener) Start() {
	if gl.running {
		return
	}
	gl.running = true
	gl.queue = make(chan Event)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	gl.startSignalHandler(sigs)

	go keyboard.Listen(func(k keys.Key) (bool, error) {
		event, stop := mapKeyToEvent(k)
		if event != None {
			gl.queue <- event
		}
		return stop, nil
	})
}

func (gl *gameListener) Check() Event {
	var event Event
	if !gl.running {
		return None
	}
	select {
	case event = <-gl.queue:
		if event == Stop {
			gl.running = false
		}
		return event
	default:
		return None
	}
}

func (gl *gameListener) Stop() {
	if gl.running {
		gl.running = false
		if gl.queue != nil {
			close(gl.queue)
		}
	}
}
