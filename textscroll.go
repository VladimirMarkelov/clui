package clui

import (
	xs "github.com/huandu/xstrings"
)

/*
A simple control to display text that is constantly changing by
adding lines to its end - example of similar behavior is 'tail'
utility. Maybe useful to display the tail of a log file or the
last application activities.

The control does not keep old text. After adding a new line the
content of the control is scrolled up and items above the top are
removed.

TextScroll can display its content in two modes: wordwrap mode
when a line of text can occupy a few lines and each line starting
from the second one has indication of wordwrap in the first
column; and simple mode when a line is truncated to control width.

TextScroll provides the following own methods:
AddItem, SetWordWrap, GetWordWrap
*/
type TextScroll struct {
	posX, posY    int
	width, height int
	title         string
	anchor        Anchor
	id            WinId
	enabled       bool
	align         Align
	active        bool
	parent        Window
	visible       bool
	tabStop       bool
	wordWrap      bool
	textColor     Color
	backColor     Color
	scale         int

	minW, minH int

	// own listbox members
	items []string
}

func NewTextScroll(parent Window, id WinId, x, y, width, height int, props Props) *TextScroll {
	l := new(TextScroll)
	l.SetEnabled(true)
	l.SetPos(x, y)
	l.SetSize(width, height)
	l.items = make([]string, 0)
	l.parent = parent
	l.visible = true
	l.tabStop = false
	l.id = id
	l.anchor = props.Anchors

	l.minW, l.minH = 5, 3

	return l
}

func (l *TextScroll) SetText(title string) {
	l.title = title
}

func (l *TextScroll) GetText() string {
	return l.title
}

func (l *TextScroll) GetId() WinId {
	return l.id
}

func (l *TextScroll) GetSize() (int, int) {
	return l.width, l.height
}

func (l *TextScroll) GetConstraints() (int, int) {
	return l.minW, l.minH
}

func (l *TextScroll) SetConstraints(minW, minH int) {
	if minW >= 5 {
		l.minW = minW
	}
	if minH >= 3 {
		l.minH = minH
	}
}

func (l *TextScroll) SetSize(width, height int) {
	width, height = ApplyConstraints(l, width, height)
	l.width = width
	l.height = height
}

func (l *TextScroll) GetPos() (int, int) {
	return l.posX, l.posY
}

func (l *TextScroll) SetPos(x, y int) {
	l.posX = x
	l.posY = y
}

func (l *TextScroll) calculateTopItem() int {
	if len(l.items) == 0 {
		return 0
	}

	idx := len(l.items) - 1
	hgt := l.height

	for {
		itemHeight := 1
		length := xs.Len(l.items[idx])
		if l.wordWrap && length > l.width {
			length -= l.width
			itemHeight += length / (l.width - 1)
			if length%(l.width-1) != 0 {
				itemHeight++
			}
		}

		if hgt <= itemHeight {
			if hgt < itemHeight && idx < len(l.items)-1 {
				idx++
			}
			break
		}

		if idx == 0 {
			break
		}

		hgt -= itemHeight
		idx--
	}

	return idx
}

func (l *TextScroll) redrawItems(canvas Canvas, tm *ThemeManager) {
	fg, bg := l.textColor, l.backColor
	if fg == ColorDefault {
		fg = tm.GetSysColor(ColorEditText)
	}
	if bg == ColorDefault {
		bg = tm.GetSysColor(ColorEditBack)
	}

	y := l.posY
	ch := string(tm.GetSysObject(ObjEditWordWrap))

	top := l.calculateTopItem()

	for idx := top; idx < len(l.items); idx++ {
		s := l.items[idx]
		w := []rune(s)

		if !l.wordWrap || len(w) <= l.width {
			canvas.DrawText(l.posX, y, l.width, s, fg, bg)
			y++
			continue
		}

		firstLine := true
		for {
			if y >= l.posY+l.height {
				break
			}

			var str string
			if firstLine {
				str = string(w[:l.width])
				w = w[l.width:]
				firstLine = false
			} else {
				if len(w) <= l.width-1 {
					str = ch + string(w)
				} else {
					str = ch + string(w[:l.width-1])
				}
				if len(w) > l.width-1 {
					w = w[l.width-1:]
				} else {
					w = make([]rune, 0)
				}
			}

			canvas.DrawText(l.posX, y, l.width, str, fg, bg)
			y++

			if len(w) == 0 {
				break
			}
		}

		if y >= l.posY+l.height {
			break
		}
	}
}

func (l *TextScroll) Redraw(canvas Canvas) {
	x, y := l.GetPos()
	w, h := l.GetSize()

	tm := canvas.Theme()
	bg := l.backColor
	if bg == ColorDefault {
		bg = tm.GetSysColor(ColorEditBack)
	}

	canvas.ClearRect(x, y, w, h, bg)
	l.redrawItems(canvas, tm)
}

func (l *TextScroll) GetEnabled() bool {
	return l.enabled
}

func (l *TextScroll) SetEnabled(active bool) {
	l.enabled = active
}

func (l *TextScroll) SetAlign(align Align) {
	l.align = align
}

func (l *TextScroll) GetAlign() Align {
	return l.align
}

func (l *TextScroll) SetAnchors(anchor Anchor) {
	l.anchor = anchor
}

func (l *TextScroll) GetAnchors() Anchor {
	return l.anchor
}

func (l *TextScroll) GetActive() bool {
	return l.active
}

func (l *TextScroll) SetActive(active bool) {
	l.active = active
}

func (l *TextScroll) GetTabStop() bool {
	return l.tabStop
}

func (l *TextScroll) SetTabStop(tab bool) {
	l.tabStop = tab
}

func (l *TextScroll) Clear() {
	l.items = make([]string, 0)
}

func (l *TextScroll) ProcessEvent(event Event) bool {
	if !l.active || !l.enabled {
		return false
	}

	return false
}

func (l *TextScroll) SetVisible(visible bool) {
	l.visible = visible
}

func (l *TextScroll) GetVisible() bool {
	return l.visible
}

// ----- own methods -------------------

// Add a new text line to the bottom of list. The
// content of the control is scrolled up if the new
// item is outside of display area
func (l *TextScroll) AddItem(item string) bool {
	if len(l.items) >= l.height {
		l.items = l.items[len(l.items)-l.height+1:]
	}
	l.items = append(l.items, item)

	go func() {
		ev := InternalEvent{act: EventRedraw, sender: l.id}
		l.parent.SendEvent(ev)
	}()

	return true
}

// Returns if word wrap mode is on
func (l *TextScroll) GetWordWrap() bool {
	return l.wordWrap
}

// Enables or disables word wrap mode
func (l *TextScroll) SetWordWrap(wrap bool) {
	l.wordWrap = wrap
}

func (l *TextScroll) GetColors() (Color, Color) {
	return l.textColor, l.backColor
}

func (l *TextScroll) SetTextColor(clr Color) {
	l.textColor = clr
}

func (l *TextScroll) SetBackColor(clr Color) {
	l.backColor = clr
}

func (l *TextScroll) HideChildren() {
	// nothing to do
}

func (l *TextScroll) GetScale() int {
	return l.scale
}

func (l *TextScroll) SetScale(scale int) {
	l.scale = scale
}
