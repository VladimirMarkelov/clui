package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
)

/*
CheckBox control. It can be two-state one(on and off) - it is default mode - or tree-state.
State values are 0=off, 1=on, 2=third state
*/
type CheckBox struct {
	ControlBase
	state       int
	allow3state bool
}

/*
NewCheckBox creates a new CheckBox control.
view - is a View that manages the control.
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width and heith - are minimal size of the control.
title - button title.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
CheckBox state can be changed using mouse or pressing space on keyboard while the control is active
*/
func NewCheckBox(view View, parent Control, width int, title string, scale int) *CheckBox {
	c := new(CheckBox)
	c.view = view
	c.parent = parent

	if width == AutoSize {
		width = xs.Len(title) + 4
	}

	c.SetSize(width, 1) // TODO: only one line checkboxes are supported at that moment
	c.SetConstraints(width, 1)
	c.state = 0
	c.SetTitle(title)
	c.SetTabStop(true)
	c.allow3state = false

	if parent != nil {
		parent.AddChild(c, scale)
	}

	return c
}

// Repaint draws the control on its View surface
func (c *CheckBox) Repaint() {
	x, y := c.Pos()
	w, h := c.Size()
	canvas := c.view.Canvas()
	tm := c.view.Screen().Theme()

	fg, bg := RealColor(tm, c.fg, ColorControlText), RealColor(tm, c.bg, ColorControlBack)
	if !c.Enabled() {
		fg, bg = RealColor(tm, c.fg, ColorControlDisabledText), RealColor(tm, c.bg, ColorControlDisabledBack)
	} else if c.Active() {
		fg, bg = RealColor(tm, c.fg, ColorControlActiveText), RealColor(tm, c.bg, ColorControlActiveBack)
	}

	parts := []rune(tm.SysObject(ObjCheckBox))

	cOpen, cClose, cCheck, cEmpty, cUnknown := parts[0], parts[1], parts[2], parts[3], parts[4]
	cState := []rune{cEmpty, cCheck, cUnknown}

	canvas.FillRect(x, y, w, h, term.Cell{Ch: ' ', Bg: bg})
	if w < 3 {
		return
	}

	canvas.PutSymbol(x, y, term.Cell{Ch: cOpen, Fg: fg, Bg: bg})
	canvas.PutSymbol(x+2, y, term.Cell{Ch: cClose, Fg: fg, Bg: bg})
	canvas.PutSymbol(x+1, y, term.Cell{Ch: cState[c.state], Fg: fg, Bg: bg})

	if w < 5 {
		return
	}

	shift, text := AlignText(c.title, w-4, c.align)
	canvas.PutText(x+4+shift, y, text, fg, bg)
}

//ProcessEvent processes all events come from the control parent. If a control
//   processes an event it should return true. If the method returns false it means
//   that the control do not want or cannot process the event and the caller sends
//   the event to the control parent
func (c *CheckBox) ProcessEvent(event Event) bool {
	if (!c.Active() && event.Type == EventKey) || !c.Enabled() {
		return false
	}

	if (event.Type == EventKey && event.Key == term.KeySpace) || event.Type == EventMouse {
		if c.state == 0 {
			c.state = 1
		} else if c.state == 2 {
			c.state = 0
		} else {
			if c.allow3state {
				c.state = 2
			} else {
				c.state = 0
			}
		}
		return true
	}

	return false
}

// SetState changes the current state of CheckBox
// Value must be 0 or 1 if Allow3State is off,
// and 0, 1, or 2 if Allow3State is on
func (c *CheckBox) SetState(val int) {
	if val < 0 {
		val = 0
	}
	if val > 1 && !c.allow3state {
		val = 1
	}
	if val > 2 {
		val = 2
	}

	c.state = val
}

// State returns current state of CheckBox
func (c *CheckBox) State() int {
	return c.state
}

// SetAllow3State sets if ComboBox should use 3 states. If the current
// state is unknown and one disables Allow3State option then the current
// value resets to off
func (c *CheckBox) SetAllow3State(enable bool) {
	if !enable && c.state == 2 {
		c.state = 0
	}
	c.allow3state = enable
}

// Allow3State returns true if ComboBox uses 3 states
func (c *CheckBox) Allow3State() bool {
	return c.allow3state
}
