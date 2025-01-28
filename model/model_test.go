package model

import (
	"testing"
)

func TestNewIndex(t *testing.T) {
	tests := []struct {
		name  string
		x, y  int64
		wantX int64
		wantY int64
	}{
		{
			name:  "create index with positive coordinates",
			x:     10,
			y:     20,
			wantX: 10,
			wantY: 20,
		},
		{
			name:  "create index with negative coordinates",
			x:     -5,
			y:     -15,
			wantX: -5,
			wantY: -15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := NewIndex(tt.x, tt.y)
			if idx.X() != tt.wantX {
				t.Errorf("NewIndex(%d, %d).X() = %d, want %d", tt.x, tt.y, idx.X(), tt.wantX)
			}
			if idx.Y() != tt.wantY {
				t.Errorf("NewIndex(%d, %d).Y() = %d, want %d", tt.x, tt.y, idx.Y(), tt.wantY)
			}
		})
	}
}
