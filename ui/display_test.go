package ui

import (
	"sync"
	"testing"
	"time"
)

// mockArea implements a test version of cursor.Area
type mockArea struct {
	content string
	updates int
	mu      sync.RWMutex
}

func (m *mockArea) Update(content string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.content = content
	m.updates++
}

func (m *mockArea) GetContent() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.content
}

func (m *mockArea) GetUpdates() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.updates
}

// mockDisplay implements Display interface for testing
type mockDisplay struct {
	area   *mockArea
	lock   chan struct{}
	closed bool
	mu     sync.RWMutex
}

func (d *mockDisplay) UpdateAndLock(content string, duration time.Duration) {
	d.mu.RLock()
	if d.closed {
		d.mu.RUnlock()
		return
	}
	d.mu.RUnlock()

	d.lock <- struct{}{}
	d.area.Update(content)
	go func() {
		time.Sleep(duration)
		<-d.lock
	}()
}

func (d *mockDisplay) Close() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return
	}
	d.closed = true
	close(d.lock)
}

func (d *mockDisplay) UpdateAndClose(finalMessage string) {
	d.mu.Lock()
	if d.closed {
		d.mu.Unlock()
		return
	}
	d.closed = true
	d.mu.Unlock()
	close(d.lock)
	d.area.Update(finalMessage)
}

func (d *mockDisplay) IsClosed() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.closed
}

func newMockDisplay() *mockDisplay {
	return &mockDisplay{
		area:   &mockArea{},
		lock:   make(chan struct{}, 1),
		closed: false,
	}
}

func TestDisplay(t *testing.T) {
	t.Run("basic update and lock", func(t *testing.T) {
		d := newMockDisplay()
		content := "test content"
		duration := 10 * time.Millisecond

		d.UpdateAndLock(content, duration)
		if d.area.GetContent() != content {
			t.Errorf("UpdateAndLock() content = %q, want %q", d.area.GetContent(), content)
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

		if !d.IsClosed() {
			t.Error("Display should be marked as closed")
		}

		// Updates after close should be ignored
		initialContent := d.area.GetContent()
		d.UpdateAndLock("after close", time.Millisecond)
		if d.area.GetContent() != initialContent {
			t.Error("Update after close should be ignored")
		}

		// Double close should not panic
		d.Close()
	})

	t.Run("update and close", func(t *testing.T) {
		d := newMockDisplay()
		finalMsg := "final message"
		d.UpdateAndClose(finalMsg)

		if !d.IsClosed() {
			t.Error("Display should be closed after UpdateAndClose")
		}
		if d.area.GetContent() != finalMsg {
			t.Errorf("Final content = %q, want %q", d.area.GetContent(), finalMsg)
		}

		// Further updates should be ignored
		d.UpdateAndLock("new content", time.Millisecond)
		if d.area.GetContent() != finalMsg {
			t.Error("Update after close should not change content")
		}

		// Double UpdateAndClose should not panic or change content
		d.UpdateAndClose("another message")
		if d.area.GetContent() != finalMsg {
			t.Error("Second UpdateAndClose should not change content")
		}
	})

	t.Run("multiple updates", func(t *testing.T) {
		d := newMockDisplay()
		updates := []string{"first", "second", "third"}
		duration := 10 * time.Millisecond

		for _, content := range updates {
			d.UpdateAndLock(content, duration)
			if d.area.GetContent() != content {
				t.Errorf("Content = %q, want %q", d.area.GetContent(), content)
			}
			time.Sleep(duration + 5*time.Millisecond) // Wait for lock to be released
		}

		if d.area.GetUpdates() != len(updates) {
			t.Errorf("Got %d updates, want %d", d.area.GetUpdates(), len(updates))
		}
	})

	t.Run("zero duration update", func(t *testing.T) {
		d := newMockDisplay()
		content := "zero duration"

		d.UpdateAndLock(content, 0)
		if d.area.GetContent() != content {
			t.Errorf("Content = %q, want %q", d.area.GetContent(), content)
		}

		// Verify lock is released immediately
		time.Sleep(time.Millisecond)
		select {
		case d.lock <- struct{}{}:
			// Expected behavior - lock should be released
		default:
			t.Error("Lock should be released immediately for zero duration")
		}
	})
}
