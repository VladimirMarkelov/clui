package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	мКнст "./пакКонстанты"
	мИнт "./пакИнтерфейсы"
)

/*
CheckBox control. It can be two-state one(on and off) - it is default mode - or three-state.
State values are 0=off, 1=on, 2=third state
*/
type CheckBox struct {
	*BaseControl
	state       int
	allow3state bool

	onChange func(int)
}

/*
CreateCheckBox creates a new CheckBox control.
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width - is minimal width of the control.
title - button title.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
CheckBox state can be changed using mouse or pressing space on keyboard while the control is active
*/
func CreateCheckBox(parent мИнт.ИВиджет, width int, title string, scale int) *CheckBox {
	c := new(CheckBox)
	c.BaseControl = NewBaseControl()
	c.parent = parent

	if width == мКнст.AutoSize {
		width = xs.Len(title) + 4
	}

	c.SetSize(width, 1) // TODO: only one line checkboxes are supported at that moment
	c.SetConstraints(width, 1)
	c.state = 0
	c.SetTitle(title)
	c.SetTabStop(true)
	c.allow3state = false
	c.onChange = nil
	c.SetScale(scale)

	if parent != nil {
		parent.AddChild(c)
	}

	return c
}

//Draw Repaint draws the control on its View surface
func (c *CheckBox) Draw() {
	if c.hidden {
		return
	}

	c.mtx.RLock()
	defer c.mtx.RUnlock()

	PushAttributes()
	defer PopAttributes()

	x, y := c.Pos()
	w, h := c.Size()

	fg, bg := RealColor(c.fg, c.Style(), мКнст.ColorControlText), RealColor(c.bg, c.Style(), мКнст.ColorControlBack)
	if !c.Enabled() {
		fg, bg = RealColor(c.fg, c.Style(), мКнст.ColorControlDisabledText), RealColor(c.bg, c.Style(), мКнст.ColorControlDisabledBack)
	} else if c.Active() {
		fg, bg = RealColor(c.fg, c.Style(), мКнст.ColorControlActiveText), RealColor(c.bg, c.Style(), мКнст.ColorControlActiveBack)
	}

	parts := []rune(SysObject(мКнст.ObjCheckBox))

	cOpen, cClose, cEmpty, cCheck, cUnknown := parts[0], parts[1], parts[2], parts[3], parts[4]
	cState := []rune{cEmpty, cCheck, cUnknown}

	SetTextColor(fg)
	SetBackColor(bg)
	FillRect(x, y, w, h, ' ')
	if w < 3 {
		return
	}

	PutChar(x, y, cOpen)
	PutChar(x+2, y, cClose)
	PutChar(x+1, y, cState[c.state])

	if w < 5 {
		return
	}

	shift, text := AlignColorizedText(c.title, w-4, c.align)
	DrawText(x+4+shift, y, text)
}

//ProcessEvent processes all events come from the control parent. If a control
//   processes an event it should return true. If the method returns false it means
//   that the control do not want or cannot process the event and the caller sends
//   the event to the control parent
func (c *CheckBox) ProcessEvent(event мИнт.ИСобытие) bool {
	if (!c.Active() && event.Type() == мИнт.EventKey) || !c.Enabled() {
		return false
	}

	if (event.Type() == мИнт.EventKey && event.Key() == term.KeySpace) || (event.Type() == мИнт.EventClick) {
		if c.state == 0 {
			c.SetState(1)
		} else if c.state == 2 {
			c.SetState(0)
		} else {
			if c.allow3state {
				c.SetState(2)
			} else {
				c.SetState(0)
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
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if val == c.state {
		return
	}

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

	if c.onChange != nil {
		go c.onChange(val)
	}
}

// State returns current state of CheckBox
func (c *CheckBox) State() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

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

// SetSize changes control size. Constant DoNotChange can be
// used as placeholder to indicate that the control attrubute
// should be unchanged.
// Method does nothing if new size is less than minimal size
// CheckBox height cannot be changed - it equals 1 always
func (c *CheckBox) SetSize(width, height int) {
	if width != мКнст.KeepValue && (width > 1000 || width < c.minW) {
		return
	}
	if height != мКнст.KeepValue && (height > 200 || height < c.minH) {
		return
	}

	if width != мКнст.KeepValue {
		c.width = width
	}

	c.height = 1
}

// OnChange sets the callback that is called whenever the state
// of the CheckBox is changed. Argument of callback is the current
// CheckBox state: 0 - off, 1 - on, 2 - third state
func (c *CheckBox) OnChange(fn func(int)) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.onChange = fn
}
