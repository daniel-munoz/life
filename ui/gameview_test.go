package ui

import (
	"testing"

	"github.com/daniel-munoz/life/event"
)

func TestNewGameView(t *testing.T) {
	stopChan := make(chan struct{}, 1)
	gv := NewGameView(-5, -5, 5, 5, stopChan)

	if gv.top != -5 || gv.left != -5 || gv.bottom != 5 || gv.right != 5 {
		t.Errorf("NewGameView coordinates = (%d,%d,%d,%d), want (-5,-5,5,5)",
			gv.top, gv.left, gv.bottom, gv.right)
	}

	if gv.paused || gv.showHelp || gv.ended {
		t.Error("New GameView should start with flags set to false")
	}

	if len(gv.actions) != 11 {
		t.Errorf("Expected 11 actions, got %d", len(gv.actions))
	}
}

func TestGameView_Movement(t *testing.T) {
	tests := []struct {
		name      string
		event     event.Event
		wantDelta struct{ top, left, bottom, right int64 }
	}{
		{
			name:  "move up",
			event: event.Up,
			wantDelta: struct{ top, left, bottom, right int64 }{
				top: -1, bottom: -1,
			},
		},
		{
			name:  "move down",
			event: event.Down,
			wantDelta: struct{ top, left, bottom, right int64 }{
				top: 1, bottom: 1,
			},
		},
		{
			name:  "move left",
			event: event.Left,
			wantDelta: struct{ top, left, bottom, right int64 }{
				left: -1, right: -1,
			},
		},
		{
			name:  "move right",
			event: event.Right,
			wantDelta: struct{ top, left, bottom, right int64 }{
				left: 1, right: 1,
			},
		},
		{
			name:  "page up",
			event: event.PageUp,
			wantDelta: struct{ top, left, bottom, right int64 }{
				top: -10, bottom: -10,
			},
		},
		{
			name:  "page down",
			event: event.PageDown,
			wantDelta: struct{ top, left, bottom, right int64 }{
				top: 10, bottom: 10,
			},
		},
		{
			name:  "page left",
			event: event.PageLeft,
			wantDelta: struct{ top, left, bottom, right int64 }{
				left: -10, right: -10,
			},
		},
		{
			name:  "page right",
			event: event.PageRight,
			wantDelta: struct{ top, left, bottom, right int64 }{
				left: 10, right: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stopChan := make(chan struct{}, 1)
			gv := NewGameView(0, 0, 10, 10, stopChan)
			initialTop, initialLeft := gv.top, gv.left
			initialBottom, initialRight := gv.bottom, gv.right

			gv.Execute(tt.event)

			if gv.top != initialTop+tt.wantDelta.top {
				t.Errorf("top = %d, want %d", gv.top, initialTop+tt.wantDelta.top)
			}
			if gv.left != initialLeft+tt.wantDelta.left {
				t.Errorf("left = %d, want %d", gv.left, initialLeft+tt.wantDelta.left)
			}
			if gv.bottom != initialBottom+tt.wantDelta.bottom {
				t.Errorf("bottom = %d, want %d", gv.bottom, initialBottom+tt.wantDelta.bottom)
			}
			if gv.right != initialRight+tt.wantDelta.right {
				t.Errorf("right = %d, want %d", gv.right, initialRight+tt.wantDelta.right)
			}
		})
	}
}

func TestGameView_ToggleFlags(t *testing.T) {
	stopChan := make(chan struct{}, 1)
	gv := NewGameView(0, 0, 10, 10, stopChan)

	// Test pause toggle
	if gv.IsPaused() {
		t.Error("Game should start unpaused")
	}
	gv.Execute(event.Pause)
	if !gv.IsPaused() {
		t.Error("Game should be paused after Pause event")
	}
	gv.Execute(event.Pause)
	if gv.IsPaused() {
		t.Error("Game should be unpaused after second Pause event")
	}

	// Test help toggle
	if gv.ShowHelp() {
		t.Error("Help should start hidden")
	}
	gv.Execute(event.Help)
	if !gv.ShowHelp() {
		t.Error("Help should be shown after Help event")
	}
	gv.ToggleHelp()
	if gv.ShowHelp() {
		t.Error("Help should be hidden after ToggleHelp")
	}
}

func TestGameView_Stop(t *testing.T) {
	stopChan := make(chan struct{}, 1)
	gv := NewGameView(0, 0, 10, 10, stopChan)

	if gv.Ended() {
		t.Error("Game should not start in ended state")
	}

	gv.Execute(event.Stop)

	if !gv.Ended() {
		t.Error("Game should be ended after Stop event")
	}

	select {
	case <-stopChan:
		// Expected: stop signal received
	default:
		t.Error("Stop event should send signal through stop channel")
	}
}

func TestGameView_Coordinates(t *testing.T) {
	stopChan := make(chan struct{}, 1)
	gv := NewGameView(-5, -5, 5, 5, stopChan)

	topLeft := gv.TopLeft()
	if topLeft.X() != -5 || topLeft.Y() != -5 {
		t.Errorf("TopLeft = (%d,%d), want (-5,-5)", topLeft.X(), topLeft.Y())
	}

	bottomRight := gv.BottomRight()
	if bottomRight.X() != 5 || bottomRight.Y() != 5 {
		t.Errorf("BottomRight = (%d,%d), want (5,5)", bottomRight.X(), bottomRight.Y())
	}
}

func TestGameView_UnknownEvent(t *testing.T) {
	stopChan := make(chan struct{}, 1)
	gv := NewGameView(0, 0, 10, 10, stopChan)

	// Store initial state
	initialTop, initialLeft := gv.top, gv.left
	initialBottom, initialRight := gv.bottom, gv.right
	initialPaused := gv.paused
	initialHelp := gv.showHelp
	initialEnded := gv.ended

	// Execute unknown event (None)
	gv.Execute(event.None)

	// Verify no state changes
	if gv.top != initialTop || gv.left != initialLeft ||
		gv.bottom != initialBottom || gv.right != initialRight ||
		gv.paused != initialPaused || gv.showHelp != initialHelp ||
		gv.ended != initialEnded {
		t.Error("Unknown event should not change GameView state")
	}
}
