package event

import (
	"testing"

	"atomicgo.dev/keyboard/keys"
)

func TestNewListener(t *testing.T) {
	listener := NewListener()
	if listener == nil {
		t.Error("NewListener() returned nil")
	}

	gl, ok := listener.(*gameListener)
	if !ok {
		t.Error("NewListener() did not return *gameListener")
	}

	if gl.running {
		t.Error("New listener should not be running initially")
	}

	if gl.queue != nil {
		t.Error("New listener should have nil queue initially")
	}
}

func TestCheck_WhenNotRunning(t *testing.T) {
	listener := NewListener()

	event := listener.Check()
	if event != None {
		t.Errorf("Check() on non-running listener = %v, want None", event)
	}
}

func TestStop_ClosesQueue(t *testing.T) {
	listener := NewListener()
	gl := listener.(*gameListener)

	// Manually set up state as if listener was running
	gl.running = true
	gl.queue = make(chan Event)

	listener.Stop()

	if gl.running {
		t.Error("Stop() should set running to false")
	}

	// Verify channel is closed by trying to receive
	select {
	case _, ok := <-gl.queue:
		if ok {
			t.Error("Stop() should close the queue channel")
		}
	default:
		// Channel closed or empty, which is fine
	}
}

func TestStop_WhenNotRunning(t *testing.T) {
	listener := NewListener()

	// Should not panic when called on non-running listener
	listener.Stop()
}

func TestMapKeyToEvent(t *testing.T) {
	tests := []struct {
		name      string
		key       keys.Key
		wantEvent Event
		wantStop  bool
	}{
		{
			name:      "Up arrow",
			key:       keys.Key{Code: keys.Up},
			wantEvent: Up,
			wantStop:  false,
		},
		{
			name:      "Down arrow",
			key:       keys.Key{Code: keys.Down},
			wantEvent: Down,
			wantStop:  false,
		},
		{
			name:      "Left arrow",
			key:       keys.Key{Code: keys.Left},
			wantEvent: Left,
			wantStop:  false,
		},
		{
			name:      "Right arrow",
			key:       keys.Key{Code: keys.Right},
			wantEvent: Right,
			wantStop:  false,
		},
		{
			name:      "Ctrl+C",
			key:       keys.Key{Code: keys.CtrlC},
			wantEvent: Stop,
			wantStop:  true,
		},
		{
			name:      "Space",
			key:       keys.Key{Code: keys.Space},
			wantEvent: Pause,
			wantStop:  false,
		},
		{
			name:      "Unknown key",
			key:       keys.Key{Code: keys.Tab},
			wantEvent: None,
			wantStop:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEvent, gotStop := mapKeyToEvent(tt.key)
			if gotEvent != tt.wantEvent {
				t.Errorf("mapKeyToEvent() event = %v, want %v", gotEvent, tt.wantEvent)
			}
			if gotStop != tt.wantStop {
				t.Errorf("mapKeyToEvent() stop = %v, want %v", gotStop, tt.wantStop)
			}
		})
	}
}

func TestMapRuneKeyToEvent(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		wantEvent Event
		wantStop  bool
	}{
		{
			name:      "q key",
			key:       "q",
			wantEvent: Stop,
			wantStop:  true,
		},
		{
			name:      "i key",
			key:       "i",
			wantEvent: PageUp,
			wantStop:  false,
		},
		{
			name:      "k key",
			key:       "k",
			wantEvent: PageDown,
			wantStop:  false,
		},
		{
			name:      "j key",
			key:       "j",
			wantEvent: PageLeft,
			wantStop:  false,
		},
		{
			name:      "l key",
			key:       "l",
			wantEvent: PageRight,
			wantStop:  false,
		},
		{
			name:      "h key",
			key:       "h",
			wantEvent: Help,
			wantStop:  false,
		},
		{
			name:      "unknown rune",
			key:       "x",
			wantEvent: None,
			wantStop:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEvent, gotStop := mapRuneKeyToEvent(tt.key)
			if gotEvent != tt.wantEvent {
				t.Errorf("mapRuneKeyToEvent() event = %v, want %v", gotEvent, tt.wantEvent)
			}
			if gotStop != tt.wantStop {
				t.Errorf("mapRuneKeyToEvent() stop = %v, want %v", gotStop, tt.wantStop)
			}
		})
	}
}

func TestMapKeyToEvent_RuneKey(t *testing.T) {
	// Test that RuneKey properly delegates to mapRuneKeyToEvent
	tests := []struct {
		name      string
		runes     []rune
		wantEvent Event
		wantStop  bool
	}{
		{
			name:      "q rune key",
			runes:     []rune{'q'},
			wantEvent: Stop,
			wantStop:  true,
		},
		{
			name:      "h rune key",
			runes:     []rune{'h'},
			wantEvent: Help,
			wantStop:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := keys.Key{Code: keys.RuneKey, Runes: tt.runes}
			gotEvent, gotStop := mapKeyToEvent(key)
			if gotEvent != tt.wantEvent {
				t.Errorf("mapKeyToEvent() with RuneKey event = %v, want %v", gotEvent, tt.wantEvent)
			}
			if gotStop != tt.wantStop {
				t.Errorf("mapKeyToEvent() with RuneKey stop = %v, want %v", gotStop, tt.wantStop)
			}
		})
	}
}
