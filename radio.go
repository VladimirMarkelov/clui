package clui

import (
	term "github.com/nsf/termbox-go"
)

/*
CheckBox control. It can be two-state one(on and off) - it is default mode - or tree-state.
State values are 0=off, 1=on, 2=third state
Minimal width of a checkbox cannot be less than 3
Own methods:
Get3State, Set3State, GetState, SetState
*/
type Radio struct {
	ControlBase
	group    *RadioGroup
	selected bool
}

func NewRadio(view View, parent Control, width int, title string, scale int) *Radio {
	c := new(Radio)
	c.view = view
	c.parent = parent

	c.SetSize(width, 1) // TODO: only one line checkboxes are supported at that moment
	c.SetConstraints(width, 1)
	c.SetTitle(title)
	c.SetTabStop(true)

	c.bg = ColorBlack
	c.fg = ColorWhite

	if parent != nil {
		parent.AddChild(c, scale)
	}

	return c
}

func (c *Radio) Repaint() {
	x, y := c.Pos()
	w, h := c.Size()
	canvas := c.view.Canvas()
	tm := c.view.Screen().Theme()

	fg, bg := RealColor(tm, c.fg, ColorControlText), RealColor(tm, c.bg, ColorControlBack)
	if !c.Enabled() {
		fg, bg = RealColor(tm, c.fg, ColorControlDisabledText), RealColor(tm, c.bg, ColorControlDisabledBack)
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

// Sets the current state of CheckBox
// Value must be 0/1 if 3State is off
// or 0/1/2 if 3State is on
func (c *Radio) SetSelected(val bool) {
	c.selected = val
}

func (c *Radio) Selected() bool {
	return c.selected
}

func (c *Radio) SetGroup(group *RadioGroup) {
	c.group = group
}
