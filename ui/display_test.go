package ui

import (
	"testing"
	"time"
)

// mockArea implements a test version of cursor.Area
type mockArea struct {
	content string
	updates int
}

func (m *mockArea) Update(content string) {
	m.content = content
	m.updates++
}

// mockDisplay implements Display interface for testing
type mockDisplay struct {
	area    *mockArea
	lock    chan struct{}
	closed  bool
	content string
}

func (d *mockDisplay) UpdateAndLock(content string, duration time.Duration) {
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

func (d *mockDisplay) Close() {
	if d.closed {
		return
	}
	d.closed = true
	close(d.lock)
}

func (d *mockDisplay) UpdateAndClose(finalMessage string) {
	if d.closed {
		return
	}
	d.Close()
	d.area.Update(finalMessage)
}

func newMockDisplay() *mockDisplay {
	return &mockDisplay{
		area:    &mockArea{},
		lock:    make(chan struct{}, 1),
		closed:  false,
		content: "",
	}
}

func TestDisplay(t *testing.T) {
	t.Run("basic update and lock", func(t *testing.T) {
		d := newMockDisplay()
		content := "test content"
		duration := 10 * time.Millisecond

		d.UpdateAndLock(content, duration)
		if d.area.content != content {
			t.Errorf("UpdateAndLock() content = %q, want %q", d.area.content, content)
		}

		// Verify lock is held
		select {
		case d.lock <- struct{}{}:
			t.Error("Lock channel should be blocked")
		case <-time.After(5 * time.Millisecond):
			// Expected behavior - channel is blocked
		}

		// Wait for lock to be released
		time.Sleep(duration + 5*time.Millisecond)
		select {
		case d.lock <- struct{}{}:
			// Expected behavior - lock is released
		case <-time.After(5 * time.Millisecond):
			t.Error("Lock should be released after duration")
		}
	})

	t.Run("close behavior", func(t *testing.T) {
		d := newMockDisplay()
		d.Close()

		if !d.closed {
			t.Error("Display should be marked as closed")
		}

		// Updates after close should be ignored
		content := "after close"
		d.UpdateAndLock(content, time.Millisecond)
		if d.area.content == content {
			t.Error("Update after close should be ignored")
		}
	})

	t.Run("update and close", func(t *testing.T) {
		d := newMockDisplay()
		finalMsg := "final message"
		d.UpdateAndClose(finalMsg)

		if !d.closed {
			t.Error("Display should be closed after UpdateAndClose")
		}
		if d.area.content != finalMsg {
			t.Errorf("Final content = %q, want %q", d.area.content, finalMsg)
		}

		// Further updates should be ignored
		d.UpdateAndLock("new content", time.Millisecond)
		if d.area.content != finalMsg {
			t.Error("Update after close should not change content")
		}
	})

	t.Run("multiple updates", func(t *testing.T) {
		d := newMockDisplay()
		updates := []string{"first", "second", "third"}
		duration := 10 * time.Millisecond

		for _, content := range updates {
			d.UpdateAndLock(content, duration)
			if d.area.content != content {
				t.Errorf("Content = %q, want %q", d.area.content, content)
			}
			time.Sleep(duration + 5*time.Millisecond) // Wait for lock to be released
		}

		if d.area.updates != len(updates) {
			t.Errorf("Got %d updates, want %d", d.area.updates, len(updates))
		}
	})
}
