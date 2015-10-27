package clui

import (
	term "github.com/nsf/termbox-go"
	"log"
)

/*
Screen is a core of the library. It dispatches keyboard and mouse messages, supports theming, and manages views
*/
type Screen interface {
	// Theme returns the current color theme
	Theme() Theme
	// PutEvent sends an event directly to the the event loop. It is used by some controls to ask Screen to repaint console
	PutEvent(event Event)
	// DestroyView removes view from the view list and makes the next view in the view stack active. It is not possible to destroy the last view - Screen must have at least one visible view
	DestroyView(view View)

	Logger() *log.Logger
}

/*
Canvas is a 'graphical' buffer than represents a View or Screen. Its size equal size of parent object and supports a full set of painting methods
*/
type Canvas interface {
	// SetSize sets the new Canvas size. If new size does not equal old size then Canvas is recreated and cleared with default colors. Both Canvas width and height must be greater than 2
	SetSize(width int, height int)
	// Size returns current Canvas size
	Size() (width int, height int)
	// PutSymbol sets value for the Canvas cell: rune and its colors. Returns result of operation: e.g, if the symbol position is outside Canvas the operation fails and the function returns false
	PutSymbol(x int, y int, symbol term.Cell) bool
	// PutText draws horizontal string on Canvas clipping by Canvas boundaries. x and y are starting point, text is a string to display, fg and bg are text and background attributes
	PutText(x int, y int, text string, fg term.Attribute, bg term.Attribute)
	// PutVerticalText draws vertical string on Canvas clipping by Canvas boundaries. x and y are starting point, text is a string to display, fg and bg are text and background attributes
	PutVerticalText(x int, y int, text string, fg term.Attribute, bg term.Attribute)
	/*
	   PutColorizedText draws multicolor string on Canvas clipping by Canvas boundaries.
	   Multiline is not supported. Align feature is limited: the text is aligned only if it is
	   shorter than maximum width, and displayed left aligned otherwise
	*/
	PutColorizedText(x int, y int, maxWidth int, text string, fg term.Attribute, bg term.Attribute, dir Direction, align Align)
	// Symbol returns current Canvas cell value at given coordinates. If coordinates are outside Canvas ok is false
	Symbol(x int, y int) (symbol term.Cell, ok bool)
	// Clear fills Canvas with given background color
	Clear(bg term.Attribute)
	// FillRect fills area of Canvas with user-defined rune and colors
	FillRect(x int, y int, width int, height int, symbol term.Cell)
	// DrawFrame paints a frame inside Canvas with optional border rune set(by default, in case of border is empty string, the rune set equals "─│┌┐└┘" - single border). The inner area of frame is not filled - in other words it is transparent
	DrawFrame(x int, y int, width int, height int, fg term.Attribute, bg term.Attribute, border string)
	// SetCursorPos sets text caret position. Used by controls like EditField
	SetCursorPos(x int, y int)
}

type Theme interface {
	SysObject(string) string
	SysColor(string) term.Attribute
	SetCurrentTheme(string) bool
	ThemeNames() []string
	ThemeInfo(string) ThemeInfo
	SetThemePath(string)
}

type View interface {
	Title() string
	SetTitle(string)
	Draw(Canvas)
	// Repaint draws the control on console surface
	Repaint()
	Constraints() (int, int)
	Size() (int, int)
	SetSize(int, int)
	Pos() (int, int)
	SetPos(int, int)
	Canvas() Canvas
	Active() bool
	SetActive(bool)
	/*
	   ProcessEvent processes all events come from the control parent. If a control
	   processes an event it should return true. If the method returns false it means
	   that the control do not want or cannot process the event and the caller sends
	   the event to the control parent
	*/
	ProcessEvent(Event) bool
	ActivateControl(Control)
	RegisterControl(Control)
	Screen() Screen
	Parent() Control
	HitTest(int, int) HitResult
	SetModal(bool)
	Modal() bool
	OnClose(func(Event))

	Paddings() (int, int, int, int)
	SetPaddings(int, int, int, int)
	AddChild(Control, int)
	SetPack(PackType)
	Pack() PackType
	Children() []Control
	ChildExists(Control) bool
	Scale() int
	SetScale(int)
	TabStop() bool
	Colors() (term.Attribute, term.Attribute)
	ActiveColors() (term.Attribute, term.Attribute)
	SetBackColor(term.Attribute)
	SetActiveBackColor(term.Attribute)
	SetTextColor(term.Attribute)
	SetActiveTextColor(term.Attribute)
	RecalculateConstraints()

	Logger() *log.Logger
}

type Control interface {
	Title() string
	SetTitle(string)
	Pos() (int, int)
	SetPos(int, int)
	Size() (int, int)
	SetSize(int, int)
	Scale() int
	SetScale(int)
	Constraints() (int, int)
	Paddings() (int, int, int, int)
	SetPaddings(int, int, int, int)
	// Repaint draws the control on its View surface
	Repaint()
	AddChild(Control, int)
	SetPack(PackType)
	Pack() PackType
	Children() []Control
	Active() bool
	SetActive(bool)
	/*
	   ProcessEvent processes all events come from the control parent. If a control
	   processes an event it should return true. If the method returns false it means
	   that the control do not want or cannot process the event and the caller sends
	   the event to the control parent
	*/
	ProcessEvent(Event) bool
	TabStop() bool
	Parent() Control
	Colors() (term.Attribute, term.Attribute)
	ActiveColors() (term.Attribute, term.Attribute)
	SetBackColor(term.Attribute)
	SetActiveBackColor(term.Attribute)
	SetTextColor(term.Attribute)
	SetActiveTextColor(term.Attribute)

	RecalculateConstraints()

	Logger() *log.Logger
}
