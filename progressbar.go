package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	"strconv"
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
	titleFg          term.Attribute
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
	b.align = AlignCenter

	if parent != nil {
		parent.AddChild(b, scale)
	}

	return b
}

// Repaint draws the control on its View surface.
// Horizontal ProgressBar supports custom title over the bar.
// One can set title using method SetTitle. There are a few
// predefined variables that can be used inside title to
// show value or total progress. Variable must be enclosed
// with double curly brackets. Available variables:
// percent - the current percentage rounded to closest integer
// value - raw ProgressBar value
// min - lower ProgressBar limit
// max - upper ProgressBar limit
// Examples:
//      pb.SetTitle("{{value}} of {{max}}")
//      pb.SetTitle("{{percent}}%")
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

	var title string
	if b.direction == Horizontal && b.Title() != "" {
		title = b.Title()
		title = strings.Replace(title, "{{percent}}", strconv.Itoa(prc), -1)
		title = strings.Replace(title, "{{value}}", strconv.Itoa(b.value), -1)
		title = strings.Replace(title, "{{min}}", strconv.Itoa(b.min), -1)
		title = strings.Replace(title, "{{max}}", strconv.Itoa(b.max), -1)
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

		if title != "" {
			shift, str := AlignText(title, w, b.align)
			titleClr := RealColor(tm, b.titleFg, ColorProgressTitleText)
			var sOn, sOff string
			if filled == 0 || shift >= filled {
				sOff = str
			} else if w == filled || shift+xs.Len(str) < filled {
				sOn = str
			} else {
				r := filled - shift
				sOn = xs.Slice(str, 0, r)
				sOff = xs.Slice(str, r, -1)
			}
			if sOn != "" {
				canvas.PutText(x+shift, y, sOn, titleClr, bgOn)
			}
			if sOff != "" {
				canvas.PutText(x+shift+xs.Len(sOn), y, sOff, titleClr, bgOff)
			}
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

// TitleColor returns text color of ProgressBar's title. Title
// background color always equals background color of the
// part of the ProgressBar on which it is displayed. In other
// words, background color of title is transparent
func (b *ProgressBar) TitleColor() term.Attribute {
	return b.titleFg
}

// SetTitleColor sets text color of ProgressBar's title
func (b *ProgressBar) SetTitleColor(clr term.Attribute) {
	b.titleFg = clr
}
