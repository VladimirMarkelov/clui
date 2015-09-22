package clui

import (
	"github.com/VladimirMarkelov/termbox-go"
)

const (
	// scale coefficient means never change of the object when its parent resizes
	DoNotScale int = 0
	// Used as a placeholder when width or height is not important, e.g, creating a control manually to use in packer later
	AutoSize int = -1
	// Used as a placeholder when you want to change only one value keeping others untouched. Example: ctrl.SetConstraint(10, DoNotChange)
	DoNotChange int = -1
)

// Predefined types
type (
	BorderStyle int
	WinId       int
	BorderIcon  int
	HitResult   int
	DragAction  int
	Anchor      int
	Align       int
	EventType   int
	Direction   int
	EditBoxMode int
	LayoutType  int
	PackType    int
	Color       termbox.Attribute
)

// Internal structure
type Coord struct {
	x, y, w, h int
}

// Additional properties that can be set while creating a control
type Props struct {
	// Foreground and background colors. Use ColorDefault if default color of current theme should be used
	Fg, Bg Color
	// Text alignment inside control bounding box. Default is AlignNone
	Alignment Align
	// Anchors for control. Defaut is AnchorNone. See Control.SetAnchor
	Anchors Anchor
	// Border style. Default is BorderNone. Only a few controls can have border
	Border BorderStyle
	// Direction of text output(e.g, for Label) or placing items inside control(e.g, Radio). Default is DirHorizontal
	Dir Direction
	// Used to assign additional control text(e.g, items of a ListBox)
	Text string
	// If control text is not editable. Default is false. Used by a few controls
	ReadOnly bool
}

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
	// Mouse cursor changes position. X and Y is new mouse cursor position inside console window
	EventMouseMove
	// Mouse scroll was used. X < 0 if mouse scrolls up and X > 0 if it scrolls down
	EventMouseScroll
	// Mouse button was pressed. X and Y are coordinates where mouse was pressed
	EventMousePress
	// Mouse button was released. X and Y are coordinates where mouse was released
	EventMouseRelease
	// If press and release event happen at the same coordinates then after EventMouseRelease event the library sends EventMouseClick event. X and Y are coordinates where mouse clicks. In Window build the event is the same as EventMouse for non-Windows builds
	EventMouseClick

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

// A structure to send any event inside the library
type Event struct {
	// Event type. See Event* constants
	Type EventType
	// Key code if any key was pressed (EventKey only)
	Key termbox.Key
	// Key modifier(Alt, Control or Shift) for EventKey event
	Mod termbox.Modifier
	// printable character for pressed key (EventKey only)
	Ch rune
	// multi-purpose fields. Most often they contains a coordinates where an event occured. See Event* descriptions
	X, Y int
	// extra information in case of a control content was changed (new text for EditField, value of selected item in ListBox etc)
	Msg string
	// internal ID of an event sender. Control can be retrived by calling GetControl of a Window interface
	Ctrl WinId
	// Error code if any sent by console manager library
	Err error
}

// Used as default value for ID
const NullWindow WinId = WinId(-1)

// Internal event structure. Used by Windows and controls to communicate with Composer
type InternalEvent struct {
	act    EventType
	sender WinId
	view   WinId
	msg    string
	x, y   int
}

// Internal representation of a cell of console window
type Symbol struct {
	ch rune
	fg Color
	bg Color
}

// BorderStyle constants
const (
	BorderNone = iota
	BorderSingle
	BorderDouble
)

// Color predefined values
const (
	ColorDefault       = Color(termbox.ColorDefault)
	ColorBlack         = Color(termbox.ColorBlack)
	ColorRed           = Color(termbox.ColorRed)
	ColorGreen         = Color(termbox.ColorGreen)
	ColorYellow        = Color(termbox.ColorYellow)
	ColorBlue          = Color(termbox.ColorBlue)
	ColorMagenta       = Color(termbox.ColorMagenta)
	ColorCyan          = Color(termbox.ColorCyan)
	ColorWhite         = Color(termbox.ColorWhite)
	ColorBrightBlack   = Color(termbox.ColorBlack | termbox.AttrBold)
	ColorBrightRed     = Color(termbox.ColorRed | termbox.AttrBold)
	ColorBrightGreen   = Color(termbox.ColorGreen | termbox.AttrBold)
	ColorBrightYellow  = Color(termbox.ColorYellow | termbox.AttrBold)
	ColorBrightBlue    = Color(termbox.ColorBlue | termbox.AttrBold)
	ColorBrightMagenta = Color(termbox.ColorMagenta | termbox.AttrBold)
	ColorBrightCyan    = Color(termbox.ColorCyan | termbox.AttrBold)
	ColorBrightWhite   = Color(termbox.ColorWhite | termbox.AttrBold)
)

// HitResult constants
const (
	HitOutside = iota
	HitShadow
	HitInside
	HitHeader
	HitLeftBorder
	HitRightBorder
	HitBottomBorder
	HitTopLeft
	HitTopRight
	HitBottomLeft
	HitBottomRight
	HitButtonClose
	HitButtonBottom
)

// Buttons available to use in Window title
const (
	// No button
	IconDefault = 0
	// Button to close Window
	IconClose = 1 << 0
	// Button to move Window to bottom
	IconBottom = 1 << 1
	// Button to maximize/restore window
	IconMaximize = 1 << 2
)

// InternalEvent types
const (
	// asks Composer to redraw the screen
	ActionRedraw = iota
	// asks application to close
	ActionQuit
)

// DragAction constants
const (
	// No dragging now
	DragNone = iota
	// An object is being moved with mouse
	DragMove
	// An objece is being resized with mouse
	DragResize
)

// Alignment constants
const (
	AlignLeft = iota
	AlignRight
	AlignCenter
)

// Anchor constants
// Sets the Window side which an object sticks to
const (
	AnchorNone   = 0
	AnchorLeft   = 1 << (iota - 1)
	AnchorRight  = 1 << (iota - 1)
	AnchorTop    = 1 << (iota - 1)
	AnchorBottom = 1 << (iota - 1)
	AnchorAll    = AnchorLeft | AnchorRight | AnchorTop | AnchorBottom
	AnchorWidth  = AnchorLeft | AnchorRight
	AnchorHeight = AnchorTop | AnchorBottom
)

// Output direction
// Used for Label text output direction and for Radio items distribution
const (
	DirHorizontal = iota
	DirVertical
)

// EditField modes
const (
	// Simple text edit field
	EditBoxSimple = iota
	// Text edit field with drop down ListBox and Button
	EditBoxCombo
)

// Control distribution mode
const (
	// All control positions are set by user while creating Window and its controls
	LayoutManual = iota
	// Pack mode is enabled. Create controls in required order and they are distributed automatically
	LayoutDynamic
)

// Direction of automatic control placement inside a Packer
const (
	// manual positioning - not used
	PackFixed = iota
	// Packer is 1-control high and new control is added to the right
	PackHorizontal
	// Packer is 1-control wide and new control is added to the bottom
	PackVertical
)
