package ui

import (
	"time"

	"atomicgo.dev/cursor"
)

type Display interface {
	UpdateAndLock(string, time.Duration)
	UpdateAndClose(string)
	Close()
}

type defaultDisplay struct {
	area   cursor.Area
	lock   chan struct{}
	closed bool
}

func NewDisplay() Display {
	return &defaultDisplay{
		area:   cursor.NewArea(),
		lock:   make(chan struct{}, 1),
		closed: false,
	}
}

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

func (d *defaultDisplay) Close() {
	if d.closed {
		return
	}
	d.closed = true
	close(d.lock)
}

func (d *defaultDisplay) UpdateAndClose(finalMessage string) {
	if d.closed {
		return
	}
	d.Close()
	d.area.Update(finalMessage)
}
