package clui

import (
	_ "fmt"
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	"strings"
)

type Label struct {
	ControlBase
	direction Direction
	multiline bool
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
		if l.direction == Horizontal {
			shift, text := AlignText(l.title, l.width, l.align)
			canvas.PutText(l.x+shift, l.y, text, fg, bg)
		} else {
			shift, text := AlignText(l.title, l.height, l.align)
			canvas.PutVerticalText(l.x, l.y+shift, text, fg, bg)
		}
	}
}

func (l *Label) Multiline() bool {
	return l.multiline
}

func (l *Label) SetMultiline(multi bool) {
	l.multiline = multi
}
