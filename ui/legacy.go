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

func Show(w model.World, top, left, bottom, right int64) {
	stopChannel := make(chan struct{})

	cursor.Hide()
	defer cursor.Show()

	display := NewDisplay()

	listener := event.NewListener()
	listener.Start()
	showHelp := false

	paused := false
	go func() {
		for {
			if showHelp {
				display.UpdateAndLock(options, 4500 * time.Millisecond)
  				showHelp = false
			}
			if !paused {
				w.Evolve()
			}

			topLeft, bottomRight := model.NewIndex(left,top), model.NewIndex(right,bottom)
			display.UpdateAndLock(w.WindowContent(topLeft, bottomRight), 200 * time.Millisecond)

			check := listener.Check()
			switch(check) {
			case event.Stop:
				stopChannel <- struct{}{}
				return
			case event.Up:
				top--
				bottom--
			case event.Down:
				top++
				bottom++
			case event.Left:
				left--
				right--
			case event.Right:
				left++
				right++
			case event.PageUp:
				top-=10
				bottom-=10
			case event.PageDown:
				top+=10
				bottom+=10
			case event.PageLeft:
				left-=10
				right-=10
			case event.PageRight:
				left+=10
				right+=10
			case event.Help:
				showHelp = true
			case event.Pause:
				paused = !paused
			default:
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
