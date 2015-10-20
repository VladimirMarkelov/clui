package clui

import (
	term "github.com/nsf/termbox-go"
)

const (
	// scale coefficient means never change of the object when its parent resizes
	DoNotScale int = 0
	// Used as a placeholder when width or height is not important, e.g, creating a control manually to use in packer later
	AutoSize int = -1
	// Used as a placeholder when you want to change only one value keeping others untouched. Example: ctrl.SetConstraint(10, DoNotChange)
	DoNotChange int = -1
	// NullWindow  View = nil
)

type Box struct {
	X, Y int
	W, H int
}

// Predefined types
type (
	BorderStyle int
	ViewButton  int
	HitResult   int
	Align       int
	EventType   int
	Direction   int
	PackType    int
)

// Internal event structure. Used by Windows and controls to communicate with Composer
type Event struct {
	Type   EventType
	Mod    term.Modifier
	Sender Control
	View   View
	Msg    string
	X, Y   int
	Err    error
	Key    term.Key
	Ch     rune
}

// BorderStyle constants
const (
	BorderNone = iota
	BorderSingle
	BorderDouble
)

// Color predefined values
const (
	ColorDefault     = term.ColorDefault
	ColorBlack       = term.ColorBlack
	ColorRed         = term.ColorRed
	ColorGreen       = term.ColorGreen
	ColorYellow      = term.ColorYellow
	ColorBlue        = term.ColorBlue
	ColorMagenta     = term.ColorMagenta
	ColorCyan        = term.ColorCyan
	ColorWhite       = term.ColorWhite
	ColorBlackBold   = term.ColorBlack | term.AttrBold
	ColorRedBold     = term.ColorRed | term.AttrBold
	ColorGreenBold   = term.ColorGreen | term.AttrBold
	ColorYellowBold  = term.ColorYellow | term.AttrBold
	ColorBlueBold    = term.ColorBlue | term.AttrBold
	ColorMagentaBold = term.ColorMagenta | term.AttrBold
	ColorCyanBold    = term.ColorCyan | term.AttrBold
	ColorWhiteBold   = term.ColorWhite | term.AttrBold
)

// HitResult constants
const (
	HitOutside = iota
	HitInside
	HitBorder
	HitButtonClose
	HitButtonBottom
	HitButtonMaximize
)

// Buttons available to use in Window title
const (
	// No button
	ButtonDefault = 0
	// Button to close Window
	ButtonClose = 1 << 0
	// Button to move Window to bottom
	ButtonBottom = 1 << 1
	// Button to maximize/restore window
	ButtonMaximize = 1 << 2
)

// InternalEvent types
const (
	// asks Composer to redraw the screen
	ActionRedraw = iota
	// asks application to close
	ActionQuit
)

// Alignment constants
const (
	AlignLeft = iota
	AlignRight
	AlignCenter
)

// Output direction
// Used for Label text output direction and for Radio items distribution
const (
	Horizontal = iota
	Vertical
)

// EditField modes
const (
	// Simple text edit field
	EditBoxSimple = iota
	// Text edit field with drop down ListBox and Button
	EditBoxCombo
)

const (
	ObjSingleBorder = "SingleBorder"
	ObjDoubleBorder = "DoubleBorder"
	ObjEdit         = "Edit"
	ObjScrollBar    = "ScrollBar"
	ObjViewButtons  = "ViewButtons"
	ObjCheckBox     = "CheckBox"
	ObjRadio        = "Radio"
	ObjProgressBar  = "ProgressBar"
)

const (
	// Window back and fore colors (inner area & border)
	ColorViewBack = "ViewBack"
	ColorViewText = "ViewText"

	// general colors
	ColorBack         = "Back"
	ColorText         = "Text"
	ColorDisabledText = "GrayText"
	ColorDisabledBack = "GrayBack"

	// editable & listbox-like controls
	ColorEditBack       = "EditBack"
	ColorEditText       = "EditText"
	ColorEditActiveBack = "EditActiveBack"
	ColorEditActiveText = "EditActiveText"
	ColorSelectionText  = "SelectionText"
	ColorSelectionBack  = "SelectionBack"

	// scroll control
	ColorScrollText = "ScrollText"
	ColorScrollBack = "ScrollBack"
	ColorThumbText  = "ThumbText"
	ColorThumbBack  = "ThumbBack"

	// window-like controls (button, radiogroup...)
	ColorControlText         = "ControlText"
	ColorControlBack         = "ControlBack"
	ColorControlActiveBack   = "ControlActiveBack"
	ColorControlActiveText   = "ControlActiveText"
	ColorControlDisabledBack = "ControlDisabledBack"
	ColorControlDisabledText = "ControlDisabledText"
	ColorControlShadow       = "ControlShadowBack"

	// progressbar colors
	ColorProgressBack       = "ProgressBack"
	ColorProgressText       = "ProgressText"
	ColorProgressActiveBack = "ProgressActiveBack"
	ColorProgressActiveText = "ProgressActiveText"
)

// EventType
// Event that window or control may process
// Note: Do not change events from EventKey to EventNone - they correspond to the same named events in termbox library
const (
	// a key pressed
	EventKey EventType = iota
	// an object or console size changed. X and Y are new width and height
	EventResize
	// Mouse button clicked. X and Y are coordinates of mouse click. Note: used only for non-Windows builds
	EventMouse
	// Something bad happened
	EventError
	EventInterrupt
	EventRaw
	EventNone

	// Asks an object to redraw. A library can ask a control to redraw and control can send the event to its parent to ask for total repaint, e.g, button sends redraw event after to its parent it depressed after a while to imitate real button
	EventRedraw
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
	EventClicked
	// Close application
	EventQuit
)
