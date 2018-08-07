package clui

import (
	xs "github.com/huandu/xstrings"
)

/*
Label is a decorative control that can display text in horizontal
or vertical direction. Other available text features are alignment
and multi-line ability. Text can be single- or multi-colored with
tags inside the text. Multi-colored strings have limited support
of alignment feature: if text is longer than Label width the text
is always left aligned
*/
type Label struct {
	BaseControl
	direction   Direction
	multiline   bool
	textDisplay Align
}

/*
NewLabel creates a new label.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
w and h - are minimal size of the control.
title - is Label title.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func CreateLabel(parent Control, w, h int, title string, scale int) *Label {
	c := new(Label)
	c.BaseControl = NewBaseControl()

	if w == AutoSize {
		w = xs.Len(title)
	}
	if h == AutoSize {
		h = 1
	}

	c.parent = parent

	c.SetTitle(title)
	c.SetSize(w, h)
	c.SetConstraints(w, h)
	c.SetScale(scale)
	c.tabSkip = true
	c.textDisplay = AlignLeft

	if parent != nil {
		parent.AddChild(c)
	}

	return c
}

// Direction returns direction of text output: vertical or horizontal
func (l *Label) Direction() Direction {
	return l.direction
}

// SetDirection sets the text output direction
func (l *Label) SetDirection(dir Direction) {
	l.direction = dir
}

func (l *Label) Draw() {
	if l.hidden {
		return
	}

	PushAttributes()
	defer PopAttributes()

	fg, bg := RealColor(l.fg, l.Style(), ColorText), RealColor(l.bg, l.Style(), ColorBack)
	if !l.Enabled() {
		fg = RealColor(l.fg, l.Style(), ColorDisabledText)
	}

	SetTextColor(fg)
	SetBackColor(bg)
	FillRect(l.x, l.y, l.width, l.height, ' ')

	if l.title == "" {
		return
	}

	if l.multiline {
		parser := NewColorParser(l.title, fg, bg)
		elem := parser.NextElement()
		xx, yy := l.x, l.y
		for elem.Type != ElemEndOfText {
			if xx >= l.x+l.width || yy >= l.y+l.height {
				break
			}

			if elem.Type == ElemLineBreak {
				xx = l.x
				yy += 1
			} else if elem.Type == ElemPrintable {
				SetTextColor(elem.Fg)
				SetBackColor(elem.Bg)
				putCharUnsafe(xx, yy, elem.Ch)

				if l.direction == Horizontal {
					xx += 1
					if xx >= l.x+l.width {
						xx = l.x
						yy += 1
					}
				} else {
					yy += 1
					if yy >= l.y+l.height {
						yy = l.y
						xx += 1
					}
				}
			}

			elem = parser.NextElement()
		}
	} else {
		if l.direction == Horizontal {
			shift, str := AlignColorizedText(l.title, l.width, l.align)
			if str != l.title && l.align != l.textDisplay {
				shift, str = AlignColorizedText(l.title, l.width, l.textDisplay)
			}
			DrawText(l.x+shift, l.y, str)
		} else {
			shift, str := AlignColorizedText(l.title, l.height, l.align)
			if str != l.title && l.align != l.textDisplay {
				shift, str = AlignColorizedText(l.title, l.width, l.textDisplay)
			}
			DrawTextVertical(l.x, l.y+shift, str)
		}
	}
}

// Multiline returns if text is displayed on several lines if the
// label title is longer than label width or title contains
// line breaks
func (l *Label) Multiline() bool {
	return l.multiline
}

// SetMultiline sets if the label should output text as one line
// or automatically display it in several lines
func (l *Label) SetMultiline(multi bool) {
	l.multiline = multi
}

// TextDisplay returns which part of the lable title is displayed in case of
// title is longer than the label:
// - AlignLeft - the head of the title is shown
// - AlignRight - the tail of the title is shown
// The property is used only by single line Label
func (l *Label) TextDisplay() Align {
	return l.textDisplay
}

// SetTextDisplay sets which part of the title is displayed in case of the title
// is longer than the lable. Only AlignLeft and AlignRigth are valid values
// for the property. Any other value does is skipped and does not affect
// displaying the title
func (l *Label) SetTextDisplay(align Align) {
	if align != AlignLeft && align != AlignRight {
		return
	}

	l.textDisplay = align
}
