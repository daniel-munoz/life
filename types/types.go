package types

type Index interface {
	X() int64
	Y() int64
}

type World interface {
	AddCellIn(x, y, turn int64)
	Evolve()
	WindowContent(topLeft, bottomRight Index) string
}
