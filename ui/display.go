// Package ui provides terminal-based display components for the Game of Life.
package ui

import (
	"time"

	"atomicgo.dev/cursor"
)

// Display manages terminal output with rate-limiting to prevent flickering.
type Display interface {
	// UpdateAndLock updates the display content and locks for the specified duration.
	UpdateAndLock(string, time.Duration)
	// UpdateAndClose displays a final message and closes the display.
	UpdateAndClose(string)
	// Close releases display resources.
	Close()
}

// defaultDisplay implements Display using atomicgo/cursor for terminal manipulation.
type defaultDisplay struct {
	area   cursor.Area
	lock   chan struct{}
	closed bool
}

// NewDisplay creates a new terminal display.
func NewDisplay() Display {
	return &defaultDisplay{
		area:   cursor.NewArea(),
		lock:   make(chan struct{}, 1),
		closed: false,
	}
}

// UpdateAndLock updates the display and prevents further updates for the specified duration.
func (d *defaultDisplay) UpdateAndLock(content string, duration time.Duration) {
	if d.closed {
		return
	}
	d.lock <- struct{}{}
	d.area.Update(content)
	go func() {
		time.Sleep(duration)
		<-d.lock
	}()
}

// Close releases display resources and prevents further updates.
func (d *defaultDisplay) Close() {
	if d.closed {
		return
	}
	d.closed = true
	close(d.lock)
}

// UpdateAndClose displays a final message and closes the display.
func (d *defaultDisplay) UpdateAndClose(finalMessage string) {
	if d.closed {
		return
	}
	d.Close()
	d.area.Update(finalMessage)
}
