package clui

import (
	term "github.com/nsf/termbox-go"
	мКнст "./пакКонстанты"
)

// Control is an interface that every visible control should implement
type Control interface {
	// Title returns the current title or text of the control
	Title() string
	// SetTitle changes control text or title
	SetTitle(title string)
	// Size returns current control width and height
	Size() (widht int, height int)
	// SetSize changes control size. Constant KeepValue can be
	// used as placeholder to indicate that the control attrubute
	// should be unchanged.
	SetSize(width, height int)
	// Pos returns the current absolute control position: X and Y.
	Pos() (x int, y int)
	// SetPos changes contols position. Manual call of the method does not
	// make sense for any control except for Window because control positions
	// inside of container always recalculated after its parent resizes
	SetPos(x, y int)
	// Constraints return minimal control widht and height
	Constraints() (minw int, minh int)
	SetConstraints(minw, minh int)
	// Active returns if a control is active. Only active controls can
	// process keyboard events. Parent looks for active controls to
	// make sure that there is only one active control at a time
	Active() bool
	// SetActive activates and deactivates control
	SetActive(active bool)
	// TabStop returns if a control can be selected by traversing
	// controls using TAB key
	TabStop() bool
	SetTabStop(tabstop bool)
	// Enable return if a control can process keyboard and mouse events
	Enabled() bool
	SetEnabled(enabled bool)
	// Visible return if a control is visible
	Visible() bool
	SetVisible(enabled bool)
	// Parent return control's container or nil if there is no parent container
	// that is true for Windows
	Parent() Control
	// The function should not be called manually. It is for internal use by
	// library
	SetParent(parent Control)
	// Modal returns if a control is always on top and does not allow to
	// change the current control. Used only by Windows, for other kind of
	// controls it does nothing
	Modal() bool
	SetModal(modal bool)
	// Paddings returns a number of spaces used to auto-arrange children inside
	// a container: indent from left and right sides, indent from top and bottom
	// sides.
	Paddings() (px int, py int)
	// SetPaddings changes indents for the container. Use KeepValue as a placeholder
	// if you do not want to touch a parameter
	SetPaddings(px, py int)
	// Gaps returns number of spaces inserted between child controls. dx is used
	// by horizontally-packed parents and dy by vertically-packed ones
	Gaps() (dx int, dy int)
	SetGaps(dx, dy int)
	// Pack returns direction in which a container packs
	// its children: horizontal or vertical
	Pack() мКнст.PackType
	// SetPack changes the direction of children packing
	SetPack(pack мКнст.PackType)
	// Scale return scale coefficient that is used to calculate
	// new control size after its parent resizes.
	// Fixed means the controls never changes its size.
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
	SetScale(scale int)
	// Align returns alignment of title in control
	Align() мКнст.Align
	SetAlign(align мКнст.Align)

	TextColor() term.Attribute
	// SetTextColor changes text color of the control.
	// Use ColorDefault to apply theme default color for the control
	SetTextColor(clr term.Attribute)
	BackColor() term.Attribute
	// SetBackColor changes background color of the control.
	// Use ColorDefault to apply theme default color for the control
	SetBackColor(clr term.Attribute)
	// ActiveColors return the attrubutes for the controls when it
	// is active: text and background colors.
	// Use ColorDefault to apply theme default color for the control
	ActiveColors() (term.Attribute, term.Attribute)
	// SetActiveBackColor changes background color of the active control.
	// Use ColorDefault to apply theme default color for the control
	SetActiveBackColor(term.Attribute)
	// SetActiveTextColor changes text color of the active control.
	// Use ColorDefault to apply theme default color for the control
	SetActiveTextColor(term.Attribute)

	// AddChild adds a new child to a container
	// The method should not be called manually. It is automatically called
	// if parent is not nil in Create* function
	AddChild(control Control)
	// Children returns the copy of the list of container child controls
	Children() []Control
	// ChildExists returns true if a control has argument as one of its
	// children or child of one of the children
	ChildExists(control Control) bool
	// MinimalSize returns the minimal size required by a control to show
	// it and all its children.
	MinimalSize() (w int, h int)
	// ChildrenScale returns the sum of all scales of all control decendants
	ChildrenScale() int
	// ResizeChildren recalculates new size of all control's children. Calling
	// the function manually is useless because the library calls this method
	// after any size change automatically(including call after adding a new
	// child)
	ResizeChildren()
	// PlaceChildren arranges all children inside a control. Useful to be called
	// after ResizeChildren, but manual call of the method is mostly useless.
	// The function is used by the library internally
	PlaceChildren()

	// Draw repaints the control on its parent surface
	Draw()
	// DrawChildren repaints all control children.
	// Method is added to avoid writing repetetive code for any parent control.
	// Just call the method at the end of your Draw method and all children
	// repaints automatically
	DrawChildren()

	// HitTest returns the area that corresponds to the clicked
	// position X, Y (absolute position in console window): title,
	// internal view area, title button, border or outside the control
	HitTest(x, y int) мКнст.HitResult
	// ProcessEvent processes all events come from the control parent. If a control
	// processes an event it should return true. If the method returns false it means
	// that the control do not want or cannot process the event and the caller sends
	// the event to the control parent
	ProcessEvent(ev мКнст.Event) bool
	// RefID returns the controls internal reference id
	RefID() int64
	// removeChild removes a child from a container
	// It's used to "destroy" controls whenever a control is no longer used
	// by the user
	removeChild(control Control)
	// Destroy is the public interface to remove an object from its parental chain
	// it implies this control will stop receiving events and will not be drawn nor
	// will impact on other objects position and size calculation
	Destroy()
	// SetStyle sets a control's custom style grouper/modifier, with a style set
	// the control will prefix the control theme with style, i.e if a button is modified
	// and set style to "MyCustom" then the theme will engine will first attempt to apply
	// MyCustomButtonBack and MyCustomButtonText if not present then apply the default
	// and standard ButtonBack and ButtonText
	SetStyle(style string)
	// Style returns the custom style grouper/modifier
	Style() string
	// SetClipped marks a control as clip-able, meaning the children components will not
	// affect the control's size - i.e will not make it expand
	SetClipped(clipped bool)
	// Clipped returns the current control's clipped flag
	Clipped() bool
	// Clipper if the component is clipped then return the clipper geometry, however
	// the size and pos is returned
	Clipper() (int, int, int, int)
}
