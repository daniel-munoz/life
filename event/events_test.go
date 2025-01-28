package event

import (
	"testing"
)

func TestEventConstants(t *testing.T) {
	// Test that event constants are unique
	events := []Event{
		Up, Down, Left, Right,
		PageUp, PageDown, PageLeft, PageRight,
		Help, Stop, Pause, None,
	}

	seen := make(map[Event]bool)
	for _, e := range events {
		if seen[e] {
			t.Errorf("Duplicate event value found: %v", e)
		}
		seen[e] = true
	}
}

// mockListener implements Listener interface for testing
type mockListener struct {
	events  []Event
	current int
	running bool
}

func newMockListener(events []Event) *mockListener {
	return &mockListener{
		events:  events,
		current: 0,
		running: false,
	}
}

func (ml *mockListener) Start() {
	ml.running = true
}

func (ml *mockListener) Check() Event {
	if !ml.running || ml.current >= len(ml.events) {
		return None
	}
	event := ml.events[ml.current]
	ml.current++
	if event == Stop {
		ml.running = false
	}
	return event
}

func TestListener(t *testing.T) {
	tests := []struct {
		name     string
		events   []Event
		expected []Event
	}{
		{
			name:     "basic movement sequence",
			events:   []Event{Up, Right, Down, Left},
			expected: []Event{Up, Right, Down, Left, None},
		},
		{
			name:     "stop sequence",
			events:   []Event{Up, Stop, Right},
			expected: []Event{Up, Stop, None},
		},
		{
			name:     "empty sequence",
			events:   []Event{},
			expected: []Event{None},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listener := newMockListener(tt.events)

			// Test initial state
			if event := listener.Check(); event != None {
				t.Errorf("Initial Check() = %v, want None", event)
			}

			// Start listener
			listener.Start()

			// Check sequence of events
			for i, want := range tt.expected {
				got := listener.Check()
				if got != want {
					t.Errorf("Check() %d = %v, want %v", i, got, want)
				}
				if got == Stop {
					// Verify listener stops after Stop event
					next := listener.Check()
					if next != None {
						t.Errorf("After Stop, Check() = %v, want None", next)
					}
					break
				}
			}
		})
	}
}
