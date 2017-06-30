package clui

import (
	term "github.com/nsf/termbox-go"
)

const (
	// Fixed means 'never change size of the object when its parent resizes'
	Fixed int = 0
	// AutoSize is used only in constructors. It means that the constructor
	// should either calculate the size of an object, e.g. for Label it is its text
	// length, or use default intial values
	AutoSize int = -1
	// KeepSize is used as a placeholder when you want to change only one
	// value and keep other ones untouched. Used in SetSize and SetConstraints
	// methods only
	// Example: control.SetConstraint(10, KeepValue) changes only minimal width
	// of the control and do not change the current minimal control height
	KeepValue int = -1
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
	// TableAction is a type of user-generated event for TableView
	TableAction int
	// SortOrder is a way of sorting rows in TableView
	SortOrder int
	DragType  int
)

const (
	DragNone DragType = iota
	DragMove
	DragResizeLeft
	DragResizeRight
	DragResizeBottom
	DragResizeBottomLeft
	DragResizeBottomRight
	DragResizeTopLeft
	DragResizeTopRight
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
}

// BorderStyle constants
const (
	BorderNone BorderStyle = iota
	BorderThin
	BorderThick
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
	HitOutside HitResult = iota
	HitInside
	HitBorder
	HitTop
	HitBottom
	HitRight
	HitLeft
	HitTopLeft
	HitTopRight
	HitBottomRight
	HitBottomLeft
	HitButtonClose
	HitButtonBottom
	HitButtonMaximize
)

// VeiwButton values - list of buttons available for using in View title
const (
	// ButtonDefault - no button
	ButtonDefault ViewButton = 0
	// ButtonClose - button to close View
	ButtonClose = 1 << 0
	// ButtonBottom -  move Window to bottom of the View stack
	ButtonBottom = 1 << 1
	// ButtonMaximaize - maximize and restore View
	ButtonMaximize = 1 << 2
)

// Alignment constants
const (
	AlignLeft Align = iota
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

// Available object identifiers that can be used in themes
const (
	ObjSingleBorder = "SingleBorder"
	ObjDoubleBorder = "DoubleBorder"
	ObjEdit         = "Edit"
	ObjScrollBar    = "ScrollBar"
	ObjViewButtons  = "ViewButtons"
	ObjCheckBox     = "CheckBox"
	ObjRadio        = "Radio"
	ObjProgressBar  = "ProgressBar"
	ObjBarChart     = "BarChart"
	ObjSparkChart   = "SparkChart"
	ObjTableView    = "TableView"
)

// Available color identifiers that can be used in themes
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

	// button control
	ColorButtonBack         = "ButtonBack"
	ColorButtonText         = "ButtonText"
	ColorButtonActiveBack   = "ButtonActiveBack"
	ColorButtonActiveText   = "ButtonActiveText"
	ColorButtonShadow       = "ButtonShadowBack"
	ColorButtonDisabledBack = "ButtonDisabledBack"
	ColorButtonDisabledText = "ButtonDisabledText"

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
	ColorProgressTitleText  = "ProgressTitle"

	// barchart colors
	ColorBarChartBack = "BarChartBack"
	ColorBarChartText = "BarChartText"

	// sparkchart colors
	ColorSparkChartBack    = "SparkChartBack"
	ColorSparkChartText    = "SparkChartText"
	ColorSparkChartBarBack = "SparkChartBarBack"
	ColorSparkChartBarText = "SparkChartBarText"
	ColorSparkChartMaxBack = "SparkChartMaxBack"
	ColorSparkChartMaxText = "SparkChartMaxText"

	// tableview colors
	ColorTableText           = "TableText"
	ColorTableBack           = "TableBack"
	ColorTableSelectedText   = "TableSelectedText"
	ColorTableSelectedBack   = "TableSelectedBack"
	ColorTableActiveCellText = "TableActiveCellText"
	ColorTableActiveCellBack = "TableActiveCellBack"
	ColorTableLineText       = "TableLineText"
	ColorTableHeaderText     = "TableHeaderText"
	ColorTableHeaderBack     = "TableHeaderBack"
)

// EventType is event that window or control may process
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
)

// ConfirmationDialog and SelectDialog exit codes
const (
	// DialogClosed - a user clicked close button on the dialog title
	DialogClosed = -1
	// DialogAlive - a user does not close the dialog yet, exit code is unavailable
	DialogAlive = 0
	// DialogButton1 - a user clicked the first button in the dialog (by default, it is 'Yes' or 'OK')
	DialogButton1 = 1
	// DialogButton2 - a user clicked the second button in the dialog
	DialogButton2 = 2
	// DialogButton3 - a user clicked the third button in the dialog
	DialogButton3 = 3
)

// Predefined sets of the buttons for ConfirmationDialog and SelectDialog
var (
	ButtonsOK          = []string{"OK"}
	ButtonsYesNo       = []string{"Yes", "No"}
	ButtonsYesNoCancel = []string{"Yes", "No", "Cancel"}
)

// SelectDialogType constants
const (
	// SelectDialogList - all items are displayed in a ListBox
	SelectDialogList SelectDialogType = iota
	// SelectDialogList - all items are displayed in a RadioGroup
	SelectDialogRadio
)

// TableAction constants
const (
	// A user pressed F2 or Enter key in TableView
	TableActionEdit TableAction = iota
	// A user pressed Insert key in TableView
	TableActionNew
	// A user pressed Delete key in TableView
	TableActionDelete
	// A user clicked on a column header in TableView
	TableActionSort
)

// SortOrder constants
const (
	// Do not sort
	SortNone SortOrder = iota
	// Sort ascending
	SortAsc
	// Sort descending
	SortDesc
)
