package пакИнтерфейсы

import (
	term "github.com/nsf/termbox-go"
)

//HitResult -- Used in mouse click events
type HitResult int

// EventType is a type of event fired by an object
// EventType is event that window or control may process
type EventType int

// Note: Do not change events from EventKey to EventNone - they correspond to the same named events in termbox library
const (
	// a key pressed
	EventKey EventType = iota
	// an object or console size changed. X and Y are new width and height
	EventResize
	// Mouse button clicked. X and Y are coordinates of mouse click
	EventMouse
	// Something bad happened
	EventError
	EventInterrupt
	EventRaw
	EventNone

	// Asks an object to redraw. A library can ask a control to redraw and control can send the event to its parent to ask for total repaint, e.g, button sends redraw event after to its parent it depressed after a while to imitate real button
	EventRedraw = iota + 100
	// an object that receives the event should close and destroys itself
	EventClose
	// Notify an object when it is activated or deactivated. X determines whether the object is activated or deactivated(0 - deactivated, 1 - activated)
	EventActivate
	// An object changes its position. X and Y are new coordinates of the object
	EventMove

	/*
	   control events
	*/
	// Content of a control changed. E.g, EditField text changed, selected item of ListBox changed etc
	// X defines how the content was changed: 0 - by pressing any key, 1 - by clicking mouse. This is used by compound controls, e.g, child ListBox of ComboBox should change its parent EditField text when a user selects a new item an ListBox with arrow keys and the ListBox should be closed if a user clicks on ListBox item
	EventChanged
	// Button event - button was clicked
	EventClick
	// dialog closed
	EventDialogClose
	// Close application
	EventQuit
	// Close top window - or application is there is only one window
	EventCloseWindow
	// Make a control (Target field of Event structure) to recalculate and reposition all its children
	EventLayout
	// A scroll-able control's child has been activated, then notify its parent to handle
	// the scrolling
	EventActivateChild
)

//ИСобытие -- интерфейс для событий
type ИСобытие interface {
	Type() EventType
	TypeSet(EventType)
	Mod() term.Modifier
	Msg() string
	MsgSet(string)
	X() int
	SetX(int)
	Y() int
	SetY(int)
	Err() error
	Key() term.Key
	KeySet(term.Key)
	Ch() rune
	Width() int
	Height() int
	Target() ИВиджет
	TargetSet(ИВиджет)
}
