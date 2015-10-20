package clui

import (
	_ "fmt"
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
)

type Label struct {
	ControlBase
	direction Direction
}

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

func (l *Label) Direction() Direction {
	return l.direction
}

func (l *Label) SetDirection(dir Direction) {
	l.direction = dir
}

func (l *Label) Repaint() {
	canvas := l.view.Canvas()
	tm := l.view.Screen().Theme()

	fg, bg := RealColor(tm, l.fg, ColorText), RealColor(tm, l.bg, ColorBack)
	if !l.Enabled() {
		fg = RealColor(tm, l.fg, ColorDisabledText)
	}

	canvas.FillRect(l.x, l.y, l.width, l.height, term.Cell{Ch: ' ', Fg: fg, Bg: bg})

	if l.direction == Horizontal {
		shift, text := AlignText(l.title, l.width, l.align)
		canvas.PutText(l.x+shift, l.y, text, fg, bg)
	} else {
		shift, text := AlignText(l.title, l.height, l.align)
		canvas.PutVerticalText(l.x, l.y+shift, text, fg, bg)
	}
}
