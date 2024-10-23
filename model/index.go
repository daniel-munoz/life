package model

import (
	"github.com/daniel-munoz/life/types"
	"github.com/daniel-munoz/life/model/internal"
)

func NewIndex(x, y int64) types.Index {
	return internal.NewIndex(x, y)
}
