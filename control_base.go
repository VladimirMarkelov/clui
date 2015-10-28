package clui

import (
	"fmt"
	term "github.com/nsf/termbox-go"
	"log"
)

// ControlBase is a base for all visible controls.
// Every new control must inherit it or implement
// the same set of methods
type ControlBase struct {
	title           string
	x, y            int
	width, height   int
	minW, minH      int
	scale           int
	fg, bg          term.Attribute
	fgActive        term.Attribute
	bgActive        term.Attribute
	tabSkip         bool
	disabled        bool
	align           Align
	parent          Control
	view            View
	active          bool
	padX, padY      int
	padTop, padSide int
}

// Title returns the current title or text of the control
func (c *ControlBase) Title() string {
	return c.title
}

// SetTitle changes control text or title
func (c *ControlBase) SetTitle(title string) {
	c.title = title
}

// Size returns current control width and height
func (c *ControlBase) Size() (int, int) {
	return c.width, c.height
}

// SetSize changes control size. Constant DoNotChange can be
// used as placeholder to indicate that the control attrubute
// should be unchanged.
// Method panics if new size is less than minimal size
func (c *ControlBase) SetSize(width, height int) {
	if width != DoNotChange && (width > 1000 || width < c.minW) {
		panic(fmt.Sprintf("Invalid width: %v", width))
	}
	if height != DoNotChange && (height > 200 || height < c.minH) {
		panic(fmt.Sprintf("Invalid height: %v", height))
	}

	if width != DoNotChange {
		c.width = width
	}
	if height != DoNotChange {
		c.height = height
	}
}

// Pos returns the current control position: X and Y.
// For View the position's origin is top left corner of console window,
// for other controls the origin is top left corner of View that hold
// the control
func (c *ControlBase) Pos() (int, int) {
	return c.x, c.y
}

// SetPos changes contols position. Manual call of the method does not
// make sense for any control except View because control positions
// inside of container always recalculated after View resizes
func (c *ControlBase) SetPos(x, y int) {
	c.x = x
	c.y = y
}

// applyConstraints checks if the current size fits minimal size.
// Contol size is increased if its size is less than the current
// contol minimal size
func (c *ControlBase) applyConstraints() {
	w, h := c.Size()
	wM, hM := c.Constraints()

	newW, newH := w, h
	if w < wM {
		newW = wM
	}
	if h < hM {
		newH = hM
	}

	if newW != w || newH != h {
		c.SetSize(newW, newH)
	}
}

// SetConstraints sets new minimal size of control.
// If minimal size of the control is greater than the current
// control size then the control size is changed to fit minimal values
func (c *ControlBase) SetConstraints(width, height int) {
	if width >= 1 {
		c.minW = width
	}
	if height >= 1 {
		c.minH = height
	}

	c.applyConstraints()
}

// Constraints return minimal control widht and height
func (c *ControlBase) Constraints() (int, int) {
	return c.minW, c.minH
}

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
func (c *ControlBase) Scale() int {
	return c.scale
}

// SetScale sets a scale coefficient for the control.
// See Scale method for details
func (c *ControlBase) SetScale(scale int) {
	c.scale = scale
}

// Pack returns direction in which a container packs
// its children: horizontal or vertical
func (c *ControlBase) Pack() PackType {
	return Vertical
}

// SetPack changes the direction of children packing
func (c *ControlBase) SetPack(pk PackType) {
}

// AddChild adds a new child to a container. For the most
// of controls the method is just a stub that panics
// because not every control can be a container
func (c *ControlBase) AddChild(ctrl Control, scale int) {
	panic("This control cannot have children")
}

// Children returns the list of container child controls
func (c *ControlBase) Children() []Control {
	return make([]Control, 0)
}

// Colors return the basic attrubutes for the controls: text
// attribute and background one. Some controls inroduce their
// own additional controls: see ProgressBar
func (c *ControlBase) Colors() (term.Attribute, term.Attribute) {
	return c.fg, c.bg
}

// SetTextColor changes text color of the control
func (c *ControlBase) SetTextColor(clr term.Attribute) {
	c.fg = clr
}

// SetBackColor changes background color of the control
func (c *ControlBase) SetBackColor(clr term.Attribute) {
	c.bg = clr
}

// ActiveColors return the attrubutes for the controls when it
// is active: text and background colors
func (c *ControlBase) ActiveColors() (term.Attribute, term.Attribute) {
	return c.fg, c.bg
}

// SetActiveTextColor changes text color of the active control
func (c *ControlBase) SetActiveTextColor(clr term.Attribute) {
	c.fg = clr
}

// SetActiveBackColor changes background color of the active control
func (c *ControlBase) SetActiveBackColor(clr term.Attribute) {
	c.bg = clr
}

// TabStop returns if a control can be selected by traversing
// controls using TAB key
func (c *ControlBase) TabStop() bool {
	return !c.tabSkip
}

// SetTabStop sets if a control can be selected by pressing
// TAB key
func (c *ControlBase) SetTabStop(skip bool) {
	c.tabSkip = !skip
}

// Enabled returns if controls is enabled. Disabled controls
// do not process events and usually have different look
func (c *ControlBase) Enabled() bool {
	return !c.disabled
}

// SetEnabled enables or disables control
func (c *ControlBase) SetEnabled(enable bool) {
	c.disabled = !enable
}

// SetAlign sets text alignment for some controls(Label, CheckBox etc)
func (c *ControlBase) SetAlign(align Align) {
	c.align = align
}

// Align return text alignment
func (c *ControlBase) Align() Align {
	return c.align
}

// Active returns if a control is active. Only active controls can
// process keyboard events. Parent View looks for active controls to
// make sure that there is only one active control at a time
func (c *ControlBase) Active() bool {
	return c.active
}

// SetActive activates and deactivates control
func (c *ControlBase) SetActive(active bool) {
	c.active = active
}

// ProcessEvent processes all events come from the control parent. If a control
// processes an event it should return true. If the method returns false it means
// that the control do not want or cannot process the event and the caller sends
// the event to the control parent
func (c *ControlBase) ProcessEvent(ev Event) bool {
	return false
}

// Parent return control's container or nil if there is no parent container
func (c *ControlBase) Parent() Control {
	return c.parent
}

// RecalculateConstraints used by containers to recalculate new minimal size
// depending on its children constraints after a new child is added
func (c *ControlBase) RecalculateConstraints() {
}

// Paddings returns a number of spaces used to auto-arrange children inside
// a container: indent from left and right sides, indent from top and bottom
// sides, horizontal space between controls, vertical space between controls.
// Horizontal space is used in case of PackType is horizontal, and vertical
// in other case
func (c *ControlBase) Paddings() (int, int, int, int) {
	return c.padSide, c.padTop, c.padX, c.padY
}

// SetPaddings changes indents for the container. Use DoNotChange as a placeholder
// if you do not want to touch a parameter
func (c *ControlBase) SetPaddings(side, top, dx, dy int) {
	if side >= 0 {
		c.padSide = side
	}
	if top >= 0 {
		c.padTop = top
	}
	if dx >= 0 {
		c.padX = dx
	}
	if dy >= 0 {
		c.padY = dy
	}
}

//---------- debug ----------------
func (c *ControlBase) Logger() *log.Logger {
	if c.parent == nil {
		return nil
	}
	return c.parent.Logger()
}
