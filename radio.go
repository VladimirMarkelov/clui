package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
)

/*
Radio button control. Unite a few radios in one radio group to
make a user select one of available choices.
*/
type Radio struct {
	ControlBase
	group    *RadioGroup
	selected bool
}

/*
NewRadio creates a new radio button.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width - is minimal width of the control.
title - radio title.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func NewRadio(view View, parent Control, width int, title string, scale int) *Radio {
	c := new(Radio)

	if width == AutoSize {
		width = xs.Len(title) + 4
	}

	c.view = view
	c.parent = parent

	c.SetSize(width, 1) // TODO: only one line heigth is supported at that moment
	c.SetConstraints(width, 1)
	c.SetTitle(title)
	c.SetTabStop(true)

	if parent != nil {
		parent.AddChild(c, scale)
	}

	return c
}

// Repaint draws the control on its View surface
func (c *Radio) Repaint() {
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

	parts := []rune(tm.SysObject(ObjRadio))

	cOpen, cClose, cEmpty, cCheck := parts[0], parts[1], parts[2], parts[3]
	cState := cEmpty
	if c.selected {
		cState = cCheck
	}

	canvas.FillRect(x, y, w, h, term.Cell{Ch: ' ', Bg: bg})
	if w < 3 {
		return
	}

	canvas.PutSymbol(x, y, term.Cell{Ch: cOpen, Fg: fg, Bg: bg})
	canvas.PutSymbol(x+2, y, term.Cell{Ch: cClose, Fg: fg, Bg: bg})
	canvas.PutSymbol(x+1, y, term.Cell{Ch: cState, Fg: fg, Bg: bg})

	if w < 5 {
		return
	}

	shift, text := AlignText(c.title, w-4, c.align)
	canvas.PutText(x+4+shift, y, text, fg, bg)
}

/*
ProcessEvent processes all events come from the control parent. If a control
processes an event it should return true. If the method returns false it means
that the control do not want or cannot process the event and the caller sends
the event to the control parent.
The control processes only space button and mouse clicks to make control selected. Deselecting control is not possible: one has to click another radio of the radio group to deselect this button
*/
func (c *Radio) ProcessEvent(event Event) bool {
	if (!c.Active() && event.Type == EventKey) || !c.Enabled() {
		return false
	}

	if (event.Type == EventKey && event.Key == term.KeySpace) || event.Type == EventMouse {
		if c.group == nil {
			c.SetSelected(true)
		} else {
			c.group.SelectItem(c)
		}
		return true
	}

	return false
}

// SetSelected makes the button selected. One should not use
// the method directly, it is for RadioGroup control
func (c *Radio) SetSelected(val bool) {
	c.selected = val
}

// Selected returns if the radio is selected
func (c *Radio) Selected() bool {
	return c.selected
}

// SetGroup sets the radio group to which the radio belongs
func (c *Radio) SetGroup(group *RadioGroup) {
	c.group = group
}
