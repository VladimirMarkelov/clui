package пакСобытия

import (
	term "github.com/nsf/termbox-go"
	мИнт "../пакИнтерфейсы"
)

// Event is structure used by Views and controls to communicate with Composer
// and vice versa
type Event struct {
	// Event type - the first events are mapped to termbox Event and then a few
	// own events added to the end
	Type EventType
	// Mod - is a key modifier. Only Alt modifier is supported
	Mod term.Modifier
	// Msg is a text part of the event. Used by few events: e.g, ListBox click
	// sends a value of clicked item
	Msg string
	// X and Y are multi-purpose fields: mouse coordinated for click event,
	// X is used to indicate on/off for events like Activate
	// Y is used for vertical-based events like ListBox item selection - id of the item
	X, Y int
	// Err is error got from termbox library
	Err error
	// Key is a pressed key
	Key term.Key
	// Ch is a printable representation of pressed key combinaton
	Ch rune
	// For resize event - new terminal size
	Width  int
	Height int
	Target мИнт.ИВиджет
}
