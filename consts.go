package clui

import (
	term "github.com/nsf/termbox-go"
)

const (
	// DoNotScale means 'never change size of the object when its parent resizes'
	DoNotScale int = 0
	// AutoSize is used only in constructors. It means that the constructor
	// should either calculate the size of an object, e.g. for Label it is its text
	// length, or use default intial values
	AutoSize int = -1
	// DoNotChange is used as a placeholder when you want to change only one
	// value and keep other ones untouched. Used in SetSize and SetConstraints
	// methods only
	// Example: control.SetConstraint(10, DoNotChange) changes only minimal width
	// of the control and do not change the current minimal control height
	DoNotChange int = -1
)

// Predefined types
type (
	// BorderStyle is a kind of frame: none, single, and double
	BorderStyle int
	// ViewButton is a set of buttons displayed in a view title
	ViewButton int
	// HitResult is a type of a view area that is under mouse cursor.
	// Used in mouse click events
	HitResult int
	// Align is text align: left, right and center
	Align int
	// EventType is a type of event fired by an object
	EventType int
	// Direction indicates the direction in which a control must draw its
	// content. At that moment it can be applied to Label (text output
	// direction and to ProgressBar (direction of bar filling)
	Direction int
	// PackType sets how to pack controls inside its parent. Can be Vertical or
	// Horizontal
	PackType int
	// SelectDialogType sets the way of choosing an item from a list for
	// SelectionDialog control: a list-based selections, or radio group one
	SelectDialogType uint
)

// Event is structure used by Views and controls to communicate with Composer
// and vice versa
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

// Alignment constants
const (
	AlignLeft = iota
	AlignRight
	AlignCenter
)

// Output direction
// Used for Label text output direction and for Radio items distribution,
// and for container controls
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

// EventType is event that window or control may process
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
	// dialog closed
	EventDialogClose
	// Close application
	EventQuit
)

const (
	DialogClosed  = -1
	DialogAlive   = 0
	DialogButton1 = 1
	DialogButton2 = 2
	DialogButton3 = 3
)

var (
	ButtonsOK          = []string{"OK"}
	ButtonsYesNo       = []string{"Yes", "No"}
	ButtonsYesNoCancel = []string{"Yes", "No", "Cancel"}
)

const (
	SelectDialogList = iota
	SelectDialogRadio
)
