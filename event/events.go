package event

type Event int

const (
	Up Event = iota
	Down 
	Left
	Right
	PageUp
	PageDown
	PageLeft
	PageRight
	Help
	Stop
	Pause
	None
)

