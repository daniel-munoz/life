// Package event provides keyboard event handling for the Game of Life UI.
package event

// Event represents a user input event from the keyboard.
type Event int

// Event constants define all possible user input actions.
const (
	Up        Event = iota // Move view window up by 1
	Down                   // Move view window down by 1
	Left                   // Move view window left by 1
	Right                  // Move view window right by 1
	PageUp                 // Move view window up by page amount
	PageDown               // Move view window down by page amount
	PageLeft               // Move view window left by page amount
	PageRight              // Move view window right by page amount
	Help                   // Display help information
	Stop                   // Stop the simulation and exit
	Pause                  // Toggle pause state
	None                   // No event (default/empty state)
)

