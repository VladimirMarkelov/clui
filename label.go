package clui

import (
	_ "fmt"
	term "github.com/nsf/termbox-go"
)

type Label struct {
	ControlBase
	direction Direction
}

func NewLabel(view View, parent Control, w, h int, title string, scale int) *Label {
	c := new(Label)

	c.view = view
	c.parent = parent
	c.minW, c.minH = 1, 1

	c.SetTitle(title)
	c.SetSize(w, h)
	c.SetConstraints(w, h)
	c.tabSkip = true

	c.fg = ColorWhite
	c.bg = ColorBlackBold

	if parent != nil {
		parent.AddChild(c, scale)
	}

	return c
}

func (l *Label) Repaint() {
	canvas := l.view.Canvas()
	tm := l.view.Screen().Theme()
	fg, bg := RealColor(tm, l.fg, ColorText), RealColor(tm, l.bg, ColorBack)

	canvas.FillRect(l.x, l.y, l.width, l.height, term.Cell{Ch: ' ', Fg: fg, Bg: bg})

	shift, text := AlignText(l.title, l.width, l.align)
	canvas.PutText(l.x+shift, l.y, text, fg, bg)
}
