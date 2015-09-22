package clui

import (
	"log"
)

/*
A top level object. It is the only object that can contain and manage other
controls (a window cannot be a child of other window). A window is responsible
to process all the events(both user - mouse, keyboard - and internal ones).
Window always has a border - active one has double border, inactive one has
single border. A user can resize and move a window with a mouse or using hot
keys. Every windows uses its own coordinate system starting from the window
top left corner. So a window has an active clip box that starts at point 0,0
and ends at width-3,height-3(window size without borders).
*/
type Window interface {
	// Change window title
	SetTitle(title string)
	// Retuns current window title
	GetTitle() string
	// Retuns internal window identifier. Internal method used by parent object Composer
	GetId() WinId
	// Retuns current window size including borders
	GetSize() (int, int)
	// Sets window size (including borders). If a new size is less than current window constraint then the constraint is used
	SetSize(int, int)
	// Returns window position on the screen. Top left corner of console window is is 0,0
	GetPos() (int, int)
	// Moves window to a new position. Top left corner of console window is is 0,0
	SetPos(int, int)
	// Returns border style. For a window it can be only BorderSingle or BorderDouble depending on window activity flag
	GetBorderStyle() BorderStyle
	// Returns a set of icons available for the window: IconClose - window can be closed with mouse click (the last window cannot be closed), IconBotton - a user can move window to bottom (if there are more than one window), IconMaximize - to maximize and restore window
	GetBorderIcons() BorderIcon
	// Set what icons should be shown inside window title
	SetBorderIcons(icons BorderIcon)
	// Repaints a window, its decorations, and all its controls
	Redraw()
	// Returns a symbol that window displays at screen coordinates screenX and screenY. If the coordinates are outside window the ok value is false and sym has default value. The method is used internally by composer to repaint all info on the screen
	GetScreenSymbol(screenX, screenY int) (sym Symbol, ok bool)
	// Changes activity flag of a window. Internal method. It does not make the window active really, just makes window to believe it. To activate window you should use composer methods
	SetActive(active bool)
	// Returns if a window is active
	GetActive() bool
	// Returns what part of a window is the point with screen coordinates screenX,screenY. Check constants Hit*. Used by composer in mouse related methods
	HitTest(screenX, screenY int) HitResult
	// Returns minimal possible width and height of a window (negative value means no restriction). You cannot set window size less than constraints. Used by SetSize method and by composer in mouse related operations
	GetConstraints() (int, int)
	// Set new constraints for a window - the window cannot be less than these constraints. If a constraint is greater than current size then the window is resized automaticaly
	SetConstraints(minW, minH int)
	// Registers a new child inside a window. Use this method only in case of manual window layout, if you use dynamic one - then similar methods of packers must be applied.
	AddControl(control Control) WinId
	// Destroys control and frees all its memory. Do not use it if dynamic layout is turned on
	RemoveControl(control Control)
	// Make a control active. Only one control can be active and process keyboard and some mouse events at a time. A contols activates automatically when a user clicks it. Returns false if activation failed. If activation completes successfully then previously activated control gets event about deactivation and activated control gets activation event
	ActivateControl(control Control) bool
	// Processes any event(keyboard, mouse, resize, activate etc). If a window cannot process this kind of event or it does not want to do it now then the method returns false. If no control has processed the event, that event passed to up level. Example: if a ComboBox has its ListBox open then the ComboBox consumes all arrow up and down events to move cursor inside its ListBox and returns true, but if the ListBox is hidden then the methon returns false if arrow keys are pressed.
	ProcessEvent(Event) bool
	// Sends event to up level (to Window in this case). E.g, useful if window needs repainting, so it sends a redraw event to composer
	SendEvent(InternalEvent)
	// Internal method to get ID of the next control. Used by composer to assign unique IDs to all created controls
	GetNextControlId() WinId
	// Start dynamic control layout mode. A window can contain only one packer - Horizontal or Vertical one. That packer always has the same size as the window
	AddPack(PackType) Packer
	// Ends placing controls in dynamic mode. After calling this method adding new controls with dynamic layout is not possible - you can place control only at fixed position. The method returns the calculated mininal size of windows depending on controls - the numbers are used to set initial constraints
	PackEnd() (int, int)
	// Returns a child by its ID. If child is not found then nil is returned
	GetControl(WinId) Control

	// for cental debugging purpose - will be removed later
	Logger() *log.Logger
}

/*
All controls must have use output methods of this object and must
not implement their own screen output procedures.
Every canvas is attached to its own window. Only inner area of a window
(the area of window excluding its border) can be used. Coordinate system
starts at 0,0 (top left corner - marked 'A') and ends at width-2,height-2
(right bottom corner - marked 'Z'):
   +--Window title--+
   |A               |
   |                |
   |               Z|
   +----------------+
At that moment it is not possible to draw anything on the frame. It is system area for title and icons.
*/
type Canvas interface {
	// Draws a horizontal text line
	// x starting point, y starting point, maximum length (if a string is longer than max the text is truncated), text line, foreground color, background color
	DrawText(int, int, int, string, Color, Color)
	// Draws a vertical text line
	// x starting point, y starting point, maximum height (if a string is longer than max the text is truncated), text line, foreground color, background color
	DrawVerticalText(int, int, int, string, Color, Color)
	// Draws an aligned text line within a box
	// x starting point, y starting point, box width (if a string is longer than max the text is truncated otherwise the string is aligned inside the box according to the last argument value), text line, foreground color, background color, alignment
	DrawAlignedText(int, int, int, string, Color, Color, Align)
	// Draws a character
	// x position, y position, symbol, foreground color, background color
	DrawRune(int, int, rune, Color, Color)
	// Draws a frame. It draws only frame and does not clear area inside the frame
	// x starting point, y starting point, width, height, border style: double or single(if border style is none the method do nothing), foreground color, background color
	DrawFrame(int, int, int, int, BorderStyle, Color, Color)
	// Fills an area with space character
	// x starting point, y starting point, width, height, background(fill) color
	ClearRect(int, int, int, int, Color)
	// Sets text cursor to a given position inside control. Top left corner of any control is 0, 0. (Used by EditField and similar controls)
	// control that wants to show cursor inside it, x cursor position, y cursor position
	SetCursorPos(Control, int, int)
	// Returns current theme. Used to get default colors and characters to display all objects
	Theme() *ThemeManager
}

/*
Control is an UI element that cannot exist outside a window. A window
manages control position while a control just process events and draws
itself using Canvas methods on redraw event. Control is a basic element
of UI and it cannot be a parent of any other control.
A control can provide more methods that interface has but it must provide
the minimal set of methods below. See Canvas description for coordinate system
*/
type Control interface {
	// Changes title of non-interactive control(Frame, Label etc) or editable text(EditField etc)
	// if a control does not have any of above (e.g, ProgressBar) then leave the method implementation empty
	SetText(title string)
	// Returns title or modified text of a control
	// Return empty string (not nil) in case of control does not have title or editable text (e.g, ProgressBar)
	GetText() string
	// Every control has its own ID that is unique for the window that managess the control. ID is assigned to control internally when the control is created. Do not change the ID manually. The method should just return internal control ID.
	GetId() WinId
	// Returns width and height of a control
	GetSize() (int, int)
	// Sets the width and height of a control. If a size is less than existing constraint the constraint value must be used.
	SetSize(int, int)
	// Returns control position inside parent window
	GetPos() (int, int)
	// Sets a new position inside parent window. Method is useful only if window layout is fixed or displaying of temporarily child of control is required, e.g, you need to show a drop down ListBox for ComboBox control. In dynamic layout, when packers are used, positions of controls are recalculated every time the window size changed
	SetPos(int, int)
	// Repaints control
	Redraw(Canvas)
	// Returns if a control is enabled - can process events
	GetEnabled() bool
	// Enables or disables a control. Disabled control does not process any keyboard or mouse event
	SetEnabled(bool)
	// Sets alignment for control text inside a bounding control box. Some controls do not apply this property - e.g, Button text is always centered
	SetAlign(Align)
	// Get control text alignment
	GetAlign() Align
	// Sets the way of resizing and moving a control after its parent
	// window is resized. Use it only in fixed layout mode, in dynamic
	// layout it is not applied. By default anchor is none - it means
	// that control never changes its position and size (it works the
	// same way as ancor is set to left). Anchor is combination of Anchor*
	// constants, each of them determines to which side of its parent
	// window the control is stuck. There are some shorthands for common
	// types of achors, e.g, AnchorAll - sticks a control to all parent
	// sides, so it resizes when its parent resizes, useful for space
	// fillers; AnchorWidth - sticks a control to left and right sides of
	// its parent, so the control changes its width when its parent is
	//resized but it keeps its height constant
	SetAnchors(Anchor)
	// Returns the current way of auto-resizing and -moving control on its parent resize. The method does not make sense if dynamic layout is used
	GetAnchors() Anchor
	// Returns if a control is active - in other words, if a control consumes all keyboard and mouse events
	GetActive() bool
	// Makes a control active. An active control consumes all mouse and keyboard events. Only one control can be active at a time
	SetActive(bool)
	// Returns true if a control can be activated from keyboard by pressing TAB
	GetTabStop() bool
	// Sets if a control can be accessed from keyboard by pressing TAB key.
	// Note: some controls can never be activated, so they drop the value of SetTabStop argument: examples are Label, Frame etc
	SetTabStop(bool)
	// Asks a control to process any event (mouse, keyboard or internal one). If control processes the event it returns true and the caller stops the event processing. If a control cannot process event(e.g, the control is disabled) or it is an event that control does not want(e.g, mouse scroll event for Button) the control returns false and caller should send the event up by level to the control's parent
	ProcessEvent(Event) bool
	// Returns if a control is visible and should be displayed
	GetVisible() bool
	// Hide or show a contol
	SetVisible(bool)
	// Returns minimal width and height of a control. SetSize checks constraints when appliying new widht and height.
	GetConstraints() (int, int)
	// Sets the minimal width and heigh of a control. Some controls does not allow to use constraints less than certain minimum values, e.g, listbox constraints cannot be less than 3(width),5(height). If you want to change only one constraint at a time then use predefined constant DoNotChange as a placeholder for old value
	SetConstraints(int, int)
	// Sets foreground color
	SetTextColor(Color)
	// Sets background color
	SetBackColor(Color)
	// Returns foreground and background colors
	GetColors() (Color, Color)
	// Makes a contol to hide all its children, e.g, ComboBox must destroy its ListBox if it is shown
	HideChildren()
	// Returns a scale coefficient for a control. It is used only in dynamic mode
	GetScale() int
	// Sets a control scale coefficient for dynamic mode. When parent window is resized and dynamic mode is on then every control is resized according to its scale coefficient. Zero value means that the control is not resizable
	SetScale(int)
}

/*
Base block of window layout in dynamic mode. Dynamic mode is
enabled after a window AddPack is called. A window can have
only one pack and that pack always has the same size as its
parent window. But a packer can have unlimited number on child
packers. Because one cannot mix Packers and Controls in the same
Window, it is usually easier to build simple forms by manual
Control positioning and using Anchors - see demo to compare two
ways of creating layout
*/
type Packer interface {
	// Every control has its own ID that is unique for the window that managess the control. ID is assigned to control internally when the control is created. Do not change the ID manually. The method should just return internal control ID.
	GetId() WinId
	// Add a new child packer
	// A packer type(Vertical or Horizontal), packer scale coefficient
	// Returns created packer
	AddPack(PackType, int) Packer
	// Sets type of a packer: Vertical or Horizontal
	// Note: the method can be called only if the packer does not have any child yet
	SetPackType(PackType)
	// Returns type of packer: Vertical or Horizontal
	GetPackType() PackType
	// Add a control to a packer
	// Contol to add, control scale coefficient
	// Returns packed control
	PackControl(Control, int) Control
	// The method called internally after a parent window was resized to move child controls to a new calculated positons
	RepositionChildren()
	// Repant a packer and its children
	Redraw(Canvas)
	// Sets padding values used when distributing children inside a packer
	// Default values are zeroes
	// indent from left and right sides, indent from top and bottom side, horizontal space between children, vertical space between children
	SetPaddings(int, int, int, int)
	// Returns padding values: indent from left and right sides, indent from top and bottom side, horizontal space between children, vertical space between children
	GetPaddings() (int, int, int, int)
	// Calculate new children sizes after packer resizing
	// Arguments are packer size changes - difference between new packer width/height and its constraints
	ResizeChidren(int, int)
	// Returns packer position inside its parent window
	GetPos() (int, int)
	// Internal method to sets packer position. Every parent window resize repositions all children
	SetPos(int, int)
	// Returns current packer width and height
	GetSize() (int, int)
	// Internal method to change packer sizes. Every parent window resize recalculates children sizes
	SetSize(int, int)
	// Returns packer's parent. If there is no parent(e.g, the parent is window) then nil is returned
	GetContainer() Packer
	// Internal method called after all children are added to all packers. Calculates the minimal size of a packer depending on children constraints and packer paddings. After calculation the width and height are used as the packer constraints and initial packer sizes
	CalculateSize() (int, int)
	// Returns a scale coefficient for a packer
	GetScale() int
	// Sets a packer scale coefficient. When parent window or packer is resized then every packer is resized according to its scale coefficient. Zero value means that the packer is not resizable
	SetScale(int)
	// Returns minimal width and height
	GetConstraints() (int, int)
	// Sets minimal width and height. The method should not be used manually
	SetConstraints(int, int)
	// Returns top and left coordinate for the control. Used internally by PackControl and AddPack methods
	GetNextPosition() (int, int)
	// Sets the calculated position of the next control. Used internally by PackControl and AddPack methods
	SetNextPosition(int, int)
	// Set border style for a packer. By default a packer has no border. If border style is not none and current padding are zeroes then top/bottom and left/right indents changed automatically to 1 to avoid drawing children over the packer borders
	SetBorderStyle(BorderStyle)
	// Returns packer border style
	GetBorderStyle() BorderStyle
	// Returns parent window
	View() Window

	/*
	   shorthands to create a standard control and add it to a packer at the same time
	   Note: PackControl always returns Control, so you have to cast it to required object type before using unique control methods(e.g, SetValue for ProgressBar). All methods Pack* below returns an object already cast to requested control type
	*/
	// minimal width, label title, scale coefficient, additional properties
	PackLabel(int, string, int, Props) *Label
	// minimal width and height, frame title, scale coefficient, additional properties
	PackFrame(int, int, string, int, Props) *Frame
	// minimal width and height, button title, scale coefficient, additional properties
	PackButton(int, int, string, int, Props) *Button
	// minimal width, edit text, scale coefficient, additional properties
	PackEditField(int, string, int, Props) *EditField
	// minimal width and height, scale coefficient, additional properties
	PackListBox(int, int, int, Props) *ListBox
	// minimal width and height, radio group title, scale coefficient, additional properties
	// set radio group items in Props.Text = a string of items separated with '|' symbol
	PackRadioGroup(int, int, string, int, Props) *Radio
	// minimal width and height, minimal and maximum values, scale coefficient, additional properties
	PackProgressBar(int, int, int, int, int, Props) *ProgressBar
	// minimal width, checkbox title, scale coefficient, additional properties
	PackCheckBox(int, string, int, Props) *CheckBox
	// minimal width, combobox text, scale coefficient, additional properties
	PackComboBox(int, string, int, Props) *EditField
	// minimal width and height, scale coefficient, additional properties
	PackTextScroll(int, int, int, Props) *TextScroll
}
