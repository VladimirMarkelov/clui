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
	// Size returns size of the console(visible) buffer
	Size() (width int, height int)

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

// Theme is a theme manager: set of colors and object looks
type Theme interface {
	// SysObject returns object look by its id for the current
	// theme. E.g, border lines for frame or arrows for scrollbar
	SysObject(string) string
	// SysColor returns attribute by its id for the current theme
	SysColor(string) term.Attribute
	// SetCurrentTheme changes the current theme.
	// Returns false if changing failed - e.g, theme does not exist
	SetCurrentTheme(string) bool
	// ThemeNames returns the list of short theme names (file names)
	ThemeNames() []string
	// ThemeInfo returns detailed info about theme
	ThemeInfo(string) ThemeInfo
	// SetThemePath changes the directory that contains themes.
	// If new path does not equal old one, theme list reloads
	SetThemePath(string)
}

// View is an interface that every object that is managed by
// composer should implement
type View interface {
	// Title returns the current title or text of the control
	Title() string
	// SetTitle changes control text or title
	SetTitle(string)
	// Draw paints the view screen buffer to a canvas. It does not
	// repaint all view children
	Draw(Canvas)
	// Repaint draws the control on console surface
	Repaint()
	// Constraints return minimal control widht and height
	Constraints() (int, int)
	// Size returns current control width and height
	Size() (int, int)
	// SetSize changes control size. Constant DoNotChange can be
	// used as placeholder to indicate that the control attrubute
	// should be unchanged.
	// Method panics if new size is less than minimal size
	SetSize(int, int)
	// Pos returns the current control position: X and Y.
	// For View the position's origin is top left corner of console window,
	// for other controls the origin is top left corner of View that hold
	// the control
	Pos() (int, int)
	// SetPos changes contols position. Manual call of the method does not
	// make sense for any control except View because control positions
	// inside of container always recalculated after View resizes
	SetPos(int, int)
	// Canvas returns an internal graphic buffer to draw everything.
	// Used by children controls - they paint themselves on the canvas
	Canvas() Canvas
	// Active returns if a control is active. Only active controls can
	// process keyboard events. Parent View looks for active controls to
	// make sure that there is only one active control at a time
	Active() bool
	// SetActive activates and deactivates control
	SetActive(bool)
	/*
	   ProcessEvent processes all events come from the control parent. If a control
	   processes an event it should return true. If the method returns false it means
	   that the control do not want or cannot process the event and the caller sends
	   the event to the control parent
	*/
	ProcessEvent(Event) bool
	// ActivateControl make the control active and previously
	// focused control loses the focus. As a side effect the method
	// emits two events: deactivate for previously focused and
	// activate for new one if it is possible (EventActivate with
	// different X values)
	ActivateControl(Control)
	// RegisterControl adds a control to the view control list. It
	// a list of all controls visible on the view - used to
	// calculate the control under mouse when a user clicks, and
	// to calculate the next control after a user presses TAB key
	RegisterControl(Control)
	// Screen returns the composer that manages the view
	Screen() Screen
	// Parent return control's container or nil if there is no parent container
	Parent() Control
	// HitTest returns the area that corresponds to the clicked
	// position X, Y (absolute position in console window): title,
	// internal view area, title button, border or outside the view
	HitTest(int, int) HitResult
	// SetModal enables or disables modal mode
	SetModal(bool)
	// Modal returns if the view is in modal mode.In modal mode a
	// user cannot switch to any other view until the user closes
	// the modal view. Used by confirmation and select dialog to be
	// sure that the user has made a choice before continuing work
	Modal() bool
	// OnClose sets a callback that is called when view is closed.
	// For dialogs after windows is closed a user can check the
	// close result
	OnClose(func(Event))

	// Paddings returns a number of spaces used to auto-arrange children inside
	// a container: indent from left and right sides, indent from top and bottom
	// sides, horizontal space between controls, vertical space between controls.
	// Horizontal space is used in case of PackType is horizontal, and vertical
	// in other case
	Paddings() (int, int, int, int)
	// SetPaddings changes indents for the container. Use DoNotChange as a placeholder
	// if you do not want to touch a parameter
	SetPaddings(int, int, int, int)
	// AddChild add control to a list of view children. Minimal size
	// of the view calculated as a sum of sizes of its children.
	// Method panics if the same control is added twice
	AddChild(Control, int)
	// SetPack changes the direction of children packing
	SetPack(PackType)
	// Pack returns direction in which a container packs
	// its children: horizontal or vertical
	Pack() PackType
	// Children returns the list of container child controls
	Children() []Control
	// ChildExists returns true if the container already has
	// the control in its children list
	ChildExists(Control) bool
	// Scale return scale coefficient that is used to calculate
	// new control size after its parent resizes.
	// DoNotScale means the controls never changes its size.
	// Any positive value is a real coefficient of scaling.
	// How the scaling works: after resizing, parent control
	// calculates the difference between minimal and current sizes,
	// then divides the difference between controls that has
	// positive scale depending on a scale value. The more scale,
	// the larger control after resizing. Example: if you have
	// two controls with scales 1 and 2, then after every resizing
	// the latter controls expands by 100% more than the first one.
	Scale() int
	// SetScale sets a scale coefficient for the control.
	// See Scale method for details
	SetScale(int)
	// TabStop returns if a control can be selected by traversing
	// controls using TAB key
	TabStop() bool
	// Colors return the basic attrubutes for the controls: text
	// attribute and background one. Some controls inroduce their
	// own additional controls: see ProgressBar
	Colors() (term.Attribute, term.Attribute)
	// ActiveColors return the attrubutes for the controls when it
	// is active: text and background colors
	ActiveColors() (term.Attribute, term.Attribute)
	// SetBackColor changes background color of the control
	SetBackColor(term.Attribute)
	// SetActiveBackColor changes background color of the active control
	SetActiveBackColor(term.Attribute)
	// SetTextColor changes text color of the control
	SetTextColor(term.Attribute)
	// SetActiveTextColor changes text color of the active control
	SetActiveTextColor(term.Attribute)
	// RecalculateConstraints used by containers to recalculate new minimal size
	// depending on its children constraints after a new child is added
	RecalculateConstraints()
	// SetMaximized opens the view to full screen or restores its
	// previous size
	SetMaximized(maximize bool)
	// Maximized returns if the view is in full screen mode
	Maximized() bool

	Logger() *log.Logger
}

// Control is an interface that every visible control on the View must
// implement
type Control interface {
	// Title returns the current title or text of the control
	Title() string
	// SetTitle changes control text or title
	SetTitle(string)
	// Pos returns the current control position: X and Y.
	// For View the position's origin is top left corner of console window,
	// for other controls the origin is top left corner of View that hold
	// the control
	Pos() (int, int)
	// SetPos changes contols position. Manual call of the method does not
	// make sense for any control except View because control positions
	// inside of container always recalculated after View resizes
	SetPos(int, int)
	// Size returns current control width and height
	Size() (int, int)
	// SetSize changes control size. Constant DoNotChange can be
	// used as placeholder to indicate that the control attrubute
	// should be unchanged.
	// Method panics if new size is less than minimal size
	SetSize(int, int)
	// Scale return scale coefficient that is used to calculate
	// new control size after its parent resizes.
	// DoNotScale means the controls never changes its size.
	// Any positive value is a real coefficient of scaling.
	// How the scaling works: after resizing, parent control
	// calculates the difference between minimal and current sizes,
	// then divides the difference between controls that has
	// positive scale depending on a scale value. The more scale,
	// the larger control after resizing. Example: if you have
	// two controls with scales 1 and 2, then after every resizing
	// the latter controls expands by 100% more than the first one.
	Scale() int
	// SetScale sets a scale coefficient for the control.
	// See Scale method for details
	SetScale(int)
	// Constraints return minimal control widht and height
	Constraints() (int, int)
	// Paddings returns a number of spaces used to auto-arrange children inside
	// a container: indent from left and right sides, indent from top and bottom
	// sides, horizontal space between controls, vertical space between controls.
	// Horizontal space is used in case of PackType is horizontal, and vertical
	// in other case
	Paddings() (int, int, int, int)
	// SetPaddings changes indents for the container. Use DoNotChange as a placeholder
	// if you do not want to touch a parameter
	SetPaddings(int, int, int, int)
	// Repaint draws the control on its View surface
	Repaint()
	// AddChild adds a new child to a container. For the most
	// of controls the method is just a stub that panics
	// because not every control can be a container
	AddChild(Control, int)
	// SetPack changes the direction of children packing
	SetPack(PackType)
	// Pack returns direction in which a container packs
	// its children: horizontal or vertical
	Pack() PackType
	// Children returns the list of container child controls
	Children() []Control
	// Active returns if a control is active. Only active controls can
	// process keyboard events. Parent View looks for active controls to
	// make sure that there is only one active control at a time
	Active() bool
	// SetActive activates and deactivates control
	SetActive(bool)
	/*
	   ProcessEvent processes all events come from the control parent. If a control
	   processes an event it should return true. If the method returns false it means
	   that the control do not want or cannot process the event and the caller sends
	   the event to the control parent
	*/
	ProcessEvent(Event) bool
	// TabStop returns if a control can be selected by traversing
	// controls using TAB key
	TabStop() bool
	// Parent return control's container or nil if there is no parent container
	Parent() Control
	// Colors return the basic attrubutes for the controls: text
	// attribute and background one. Some controls inroduce their
	// own additional controls: see ProgressBar
	Colors() (term.Attribute, term.Attribute)
	// ActiveColors return the attrubutes for the controls when it
	// is active: text and background colors
	ActiveColors() (term.Attribute, term.Attribute)
	// SetBackColor changes background color of the control
	SetBackColor(term.Attribute)
	// SetActiveBackColor changes background color of the active control
	SetActiveBackColor(term.Attribute)
	// SetTextColor changes text color of the control
	SetTextColor(term.Attribute)
	// SetActiveTextColor changes text color of the active control
	SetActiveTextColor(term.Attribute)

	// RecalculateConstraints used by containers to recalculate new minimal size
	// depending on its children constraints after a new child is added
	RecalculateConstraints()

	Logger() *log.Logger
}
