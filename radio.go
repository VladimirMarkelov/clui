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
	BaseControl
	group    *RadioGroup
	selected bool

	onChange func(bool)
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
func CreateRadio(parent Control, width int, title string, scale int) *Radio {
	c := new(Radio)
	c.BaseControl = NewBaseControl()

	if width == AutoSize {
		width = xs.Len(title) + 4
	}

	c.parent = parent

	c.SetSize(width, 1) // TODO: only one line heigth is supported at that moment
	c.SetConstraints(width, 1)
	c.SetTitle(title)
	c.SetTabStop(true)
	c.SetScale(scale)

	c.onChange = nil

	if parent != nil {
		parent.AddChild(c)
	}

	return c
}

// Repaint draws the control on its View surface
func (c *Radio) Draw() {
	if c.hidden {
		return
	}

	PushAttributes()
	defer PopAttributes()

	x, y := c.Pos()
	w, h := c.Size()

	fg, bg := RealColor(c.fg, c.Style(), ColorControlText), RealColor(c.bg, c.Style(), ColorControlBack)
	if !c.Enabled() {
		fg, bg = RealColor(c.fg, c.Style(), ColorControlDisabledText), RealColor(c.bg, c.Style(), ColorControlDisabledBack)
	} else if c.Active() {
		fg, bg = RealColor(c.fg, c.Style(), ColorControlActiveText), RealColor(c.bg, c.Style(), ColorControlActiveBack)
	}

	parts := []rune(SysObject(ObjRadio))
	cOpen, cClose, cEmpty, cCheck := parts[0], parts[1], parts[2], parts[3]
	cState := cEmpty
	if c.selected {
		cState = cCheck
	}

	SetTextColor(fg)
	SetBackColor(bg)
	FillRect(x, y, w, h, ' ')
	if w < 3 {
		return
	}

	PutChar(x, y, cOpen)
	PutChar(x+2, y, cClose)
	PutChar(x+1, y, cState)

	if w < 5 {
		return
	}

	shift, text := AlignColorizedText(c.title, w-4, c.align)
	DrawText(x+4+shift, y, text)
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

	if (event.Type == EventKey && event.Key == term.KeySpace) || event.Type == EventClick {
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

	if c.onChange != nil {
		go c.onChange(val)
	}
}

// Selected returns if the radio is selected
func (c *Radio) Selected() bool {
	return c.selected
}

// SetGroup sets the radio group to which the radio belongs
func (c *Radio) SetGroup(group *RadioGroup) {
	c.group = group
}

// OnChange sets the callback that is called whenever the state
// of the Radio is changed. Argument of callback is the current
func (c *Radio) OnChange(fn func(bool)) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.onChange = fn
}
