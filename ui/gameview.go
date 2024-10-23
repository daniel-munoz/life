package ui

import (
	"github.com/daniel-munoz/life/event"
	"github.com/daniel-munoz/life/model"
	"github.com/daniel-munoz/life/types"
)

// GameView is the view of the game. It shows the world in a view window, defined
// by the top, left, bottom and right coordinates. It also keeps the status of the
// pause and help flags.
type GameView struct {
	top, left, bottom, right int64
	paused, showHelp, ended  bool
	actions                  map[event.Event]Action
}

// NewGameView creates a new GameView.
func NewGameView(top, left, bottom, right int64, stopChannel chan struct{}) *GameView {
	gv := &GameView{
		top:    top,
		left:   left,
		bottom: bottom,
		right:  right,
	}
	gv.actions = map[event.Event]Action{
		event.Stop: func() {
			gv.ended = true
			stopChannel <- struct{}{}
		},
		event.Up: func() {
			gv.top--
			gv.bottom--
		},
		event.Down: func() {
			gv.top++
			gv.bottom++
		},

		event.Left: func() {
			gv.left--
			gv.right--
		},
		event.Right: func() {
			gv.left++
			gv.right++
		},
		event.PageUp: func() {
			gv.top -= 10
			gv.bottom -= 10
		},
		event.PageDown: func() {
			gv.top += 10
			gv.bottom += 10
		},
		event.PageLeft: func() {
			gv.left -= 10
			gv.right -= 10
		},
		event.PageRight: func() {
			gv.left += 10
			gv.right += 10
		},
		event.Help: func() {
			gv.showHelp = true
		},
		event.Pause: func() {
			gv.paused = !gv.paused
		},
	}
	return gv
}

// TopLeft returns the top and left coordinates of the view window.
func (gv *GameView) TopLeft() types.Index {
	return model.NewIndex(gv.left, gv.top)
}

// BottomRight returns the bottom and right coordinates of the view window.
func (gv *GameView) BottomRight() types.Index {
	return model.NewIndex(gv.right, gv.bottom)
}

// IsPaused returns true if the game is paused.
func (gv *GameView) IsPaused() bool {
	return gv.paused
}

// ShowHelp returns true if the help is being shown.
func (gv *GameView) ShowHelp() bool {
	return gv.showHelp
}

// ToggleHelp toggles the help flag.
func (gv *GameView) ToggleHelp() {
	gv.showHelp = !gv.showHelp
}

// Execute executes the action associated to the given event.
func (gv *GameView) Execute(e event.Event) {
	action, ok := gv.actions[e]
	if ok {
		action()
	}
}

// Ended returns true if the game has ended.
func (gv *GameView) Ended() bool {
	return gv.ended
}
