// Package model provides the Game of Life world model and file loading.
package model

import (
	"github.com/daniel-munoz/life/model/internal"
	"github.com/daniel-munoz/life/types"
)

// NewIndex creates a new coordinate index at the specified position.
func NewIndex(x, y int64) types.Index {
	return internal.NewIndex(x, y)
}
