package clui

import (
	term "github.com/nsf/termbox-go"
	"strings"
)

/*
ProgressBar control visualizes the progression of extended operation.

The control has two sets of colors(almost all other controls have only
one set: foreground and background colors): for filled part and for
empty one. By default colors are the same.
*/
type ProgressBar struct {
	ControlBase
	direction        Direction
	min, max         int
	value            int
	emptyFg, emptyBg term.Attribute
}

/*
NewProgressBar creates a new ProgressBar.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width and heigth - are minimal size of the control.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func NewProgressBar(view View, parent Control, width, height int, scale int) *ProgressBar {
	b := new(ProgressBar)

	if height == AutoSize {
		height = 1
	}
	if width == AutoSize {
		width = 10
	}

	b.SetSize(width, height)
	b.SetConstraints(width, height)
	b.SetTabStop(false)
	b.min = 0
	b.max = 10
	b.direction = Horizontal
	b.parent = parent
	b.view = view

	if parent != nil {
		parent.AddChild(b, scale)
	}

	return b
}

// Repaint draws the control on its View surface
func (b *ProgressBar) Repaint() {
	if b.max <= b.min {
		return
	}

	canvas := b.view.Canvas()
	tm := b.view.Screen().Theme()

	fgOff, fgOn := RealColor(tm, b.fg, ColorProgressText), RealColor(tm, b.fgActive, ColorProgressActiveText)
	bgOff, bgOn := RealColor(tm, b.bg, ColorProgressBack), RealColor(tm, b.bgActive, ColorProgressActiveBack)

	parts := []rune(tm.SysObject(ObjProgressBar))
	cFilled, cEmpty := parts[0], parts[1]

	prc := 0
	if b.value >= b.max {
		prc = 100
	} else if b.value < b.max && b.value > b.min {
		prc = (100 * (b.value - b.min)) / (b.max - b.min)
	}

	x, y := b.Pos()
	w, h := b.Size()

	if b.direction == Horizontal {
		filled := prc * w / 100
		sFilled := strings.Repeat(string(cFilled), filled)
		sEmpty := strings.Repeat(string(cEmpty), w-filled)

		for yy := y; yy < y+h; yy++ {
			canvas.PutText(x, yy, sFilled, fgOn, bgOn)
			canvas.PutText(x+filled, yy, sEmpty, fgOff, bgOff)
		}
	} else {
		filled := prc * h / 100
		sFilled := strings.Repeat(string(cFilled), w)
		sEmpty := strings.Repeat(string(cEmpty), w)
		for yy := y; yy < y+h-filled; yy++ {
			canvas.PutText(x, yy, sEmpty, fgOff, bgOff)
		}
		for yy := y + h - filled; yy < y+h; yy++ {
			canvas.PutText(x, yy, sFilled, fgOn, bgOn)
		}
	}
}

//----------------- own methods -------------------------

// SetValue sets new progress value. If value exceeds ProgressBar
// limits then the limit value is used
func (b *ProgressBar) SetValue(pos int) {
	if pos < b.min {
		b.value = b.min
	} else if pos > b.max {
		b.value = b.max
	} else {
		b.value = pos
	}
}

// Value returns the current ProgressBar value
func (b *ProgressBar) Value() int {
	return b.value
}

// Limits returns current minimal and maximal values of ProgressBar
func (b *ProgressBar) Limits() (int, int) {
	return b.min, b.max
}

// SetLimits set new ProgressBar limits. The current value
// is adjusted if it exceeds new limits
func (b *ProgressBar) SetLimits(min, max int) {
	b.min = min
	b.max = max

	if b.value < b.min {
		b.value = min
	}
	if b.value > b.max {
		b.value = max
	}
}

// Step increases ProgressBar value by 1 if the value is less
// than ProgressBar high limit
func (b *ProgressBar) Step() int {
	b.value++

	if b.value > b.max {
		b.value = b.max
	}

	return b.value
}

// SecondaryColors returns text and background colors for empty
// part of the ProgressBar
func (b *ProgressBar) SecondaryColors() (term.Attribute, term.Attribute) {
	return b.emptyFg, b.emptyBg
}

// SetSecondaryColors sets new text and background colors for
// empty part of the ProgressBar
func (b *ProgressBar) SetSecondaryColors(fg, bg term.Attribute) {
	b.emptyFg, b.emptyBg = fg, bg
}
