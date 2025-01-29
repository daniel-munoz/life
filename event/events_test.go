package event

import (
	"testing"
	"time"
)

func TestEventConstants(t *testing.T) {
	// Test that event constants are unique and properly ordered
	events := []Event{
		Up, Down, Left, Right,
		PageUp, PageDown, PageLeft, PageRight,
		Help, Stop, Pause, None,
	}

	seen := make(map[Event]bool)
	for i, e := range events {
		if seen[e] {
			t.Errorf("Duplicate event value found: %v", e)
		}
		seen[e] = true

		// Verify iota ordering
		if int(e) != i {
			t.Errorf("Event %v has value %d, want %d", e, e, i)
		}
	}
}

// mockListener implements Listener interface for testing
type mockListener struct {
	events  []Event
	current int
	running bool
	queue   chan Event
}

func newMockListener(events []Event) *mockListener {
	return &mockListener{
		events:  events,
		current: 0,
		running: false,
		queue:   make(chan Event, 10),
	}
}

func (ml *mockListener) Start() {
	if ml.running {
		return
	}
	ml.running = true
	go func() {
		for ml.running && ml.current < len(ml.events) {
			ml.queue <- ml.events[ml.current]
			ml.current++
		}
	}()
	time.Sleep(10 * time.Millisecond)
}

func (ml *mockListener) Check() Event {
	if !ml.running {
		return None
	}
	select {
	case event := <-ml.queue:
		if event == Stop {
			ml.running = false
		}
		return event
	default:
		return None
	}
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
		{
			name:     "page navigation",
			events:   []Event{PageUp, PageDown, PageLeft, PageRight},
			expected: []Event{PageUp, PageDown, PageLeft, PageRight, None},
		},
		{
			name:     "help and pause",
			events:   []Event{Help, Pause, Help},
			expected: []Event{Help, Pause, Help, None},
		},
		{
			name:     "mixed events",
			events:   []Event{Up, Help, Pause, Stop, Down},
			expected: []Event{Up, Help, Pause, Stop, None},
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

			// Check sequence of events with timeout
			timeout := time.After(1 * time.Second)
			for i, want := range tt.expected {
				select {
				case <-timeout:
					t.Fatalf("Test timed out waiting for event %d", i)
				default:
					got := listener.Check()
					if got != want {
						t.Errorf("Check() %d = %v, want %v", i, got, want)
					}
					if got == None && want == None {
						return // End of sequence
					}
					if got == Stop {
						// Verify listener stops after Stop event
						time.Sleep(20 * time.Millisecond) // Wait for goroutine to process Stop
						next := listener.Check()
						if next != None {
							t.Errorf("After Stop, Check() = %v, want None", next)
						}
						return
					}
					if got != None {
						time.Sleep(5 * time.Millisecond) // Wait for next event
					}
				}
			}
		})
	}
}
