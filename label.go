package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	"strings"
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
	ControlBase
	direction  Direction
	multiline  bool
	multicolor bool
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
func NewLabel(view View, parent Control, w, h int, title string, scale int) *Label {
	c := new(Label)

	if w == AutoSize {
		w = xs.Len(title)
	}
	if h == AutoSize {
		h = 1
	}

	c.view = view
	c.parent = parent

	c.SetTitle(title)
	c.SetSize(w, h)
	c.SetConstraints(w, h)
	c.tabSkip = true

	if parent != nil {
		parent.AddChild(c, scale)
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

// Repaint draws the control on its View surface
func (l *Label) Repaint() {
	canvas := l.view.Canvas()
	tm := l.view.Screen().Theme()

	fg, bg := RealColor(tm, l.fg, ColorText), RealColor(tm, l.bg, ColorBack)
	if !l.Enabled() {
		fg = RealColor(tm, l.fg, ColorDisabledText)
	}

	canvas.FillRect(l.x, l.y, l.width, l.height, term.Cell{Ch: ' ', Fg: fg, Bg: bg})

	if l.multiline {
		lineCnt, lineLen := l.height, l.width
		if l.direction == Vertical {
			lineCnt, lineLen = l.width, l.height
		}

		lines := strings.Split(l.title, "\n")

		var realLines []string
		for _, s := range lines {
			curr := s
			for xs.Len(curr) > lineLen {
				realLines = append(realLines, xs.Slice(curr, 0, lineLen))
				curr = xs.Slice(curr, lineLen, -1)
			}
			realLines = append(realLines, curr)
		}

		idx := 0
		for idx < lineCnt && idx < len(realLines) {
			if l.direction == Horizontal {
				shift, text := AlignText(realLines[idx], l.width, l.align)
				canvas.PutText(l.x+shift, l.y+idx, text, fg, bg)
			} else {
				shift, text := AlignText(realLines[idx], l.height, l.align)
				canvas.PutVerticalText(l.x+idx, l.y+shift, text, fg, bg)
			}
			idx++
		}
	} else {
		if l.multicolor {
			max := l.width
			if l.direction == Vertical {
				max = l.height
			}
			shift, text := AlignColorizedText(l.title, max, l.align)
			if l.direction == Vertical {
				canvas.PutColorizedText(l.x, l.y+shift, max, text, fg, bg, l.direction)
			} else {
				canvas.PutColorizedText(l.x+shift, l.y, max, text, fg, bg, l.direction)
			}
		} else {
			if l.direction == Horizontal {
				shift, text := AlignText(l.title, l.width, l.align)
				canvas.PutText(l.x+shift, l.y, text, fg, bg)
			} else {
				shift, text := AlignText(l.title, l.height, l.align)
				canvas.PutVerticalText(l.x, l.y+shift, text, fg, bg)
			}
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

// MultiColored returns if the label checks and applies any
// color related tags inside its title. If MultiColores is
// false then title is displayed as is. In multicolor mode
// label has some limitations for alignment.
// To read about available color tags, please see ColorParser
func (l *Label) MultiColored() bool {
	return l.multicolor
}

// SetMultiColored changes how the label output its title: as is
// or parse and apply all internal color tags
func (l *Label) SetMultiColored(multi bool) {
	l.multicolor = multi
}
