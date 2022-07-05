package keyboard

import (
    "atomicgo.dev/keyboard"
    "atomicgo.dev/keyboard/keys"
)

type Listener interface {
	Start()
	Check() Event
}

func NewListener() Listener {
	return &gameListener{}
}

type gameListener struct {
	queue   chan Event
	running bool
}

func (gl *gameListener) Start() {
	if gl.running {
		return
	}
	gl.running = true
	gl.queue = make(chan Event)

	go keyboard.Listen(func (k keys.Key) (bool, error) {
		event, stop := None, false
		switch k.Code {
		case keys.Up:
			event = Up
		case keys.Down:
			event = Down
		case keys.Left:
			event = Left
		case keys.Right:
			event = Right
		case keys.RuneKey:
			switch k.String() {
			case "q":
				event = Stop
				stop = true
			case "i":
				event = PageUp
			case "k":
				event = PageDown
			case "j":
				event = PageLeft
			case "l":
				event = PageRight
			case "h":
				event = Help
			default:
		}
		default:
			event = None
		}
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
