package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
)

/*
Text edit field contol. Can be simple edit field(default mode) and edit field
with drop down list. The EditField mode is set during creation and cannot be changed on the fly.
Edit field consumes some keyboard events when it is active: all printable charaters;
Delete, BackSpace, Home, End, left and right arrows; Enter, up and down arrows if EditField
in combobox mode and drop down list is visible; F5 to open drop down list in combobox mode;
Ctrl+R to clear EditField.
Edit text can be limited. By default a user can enter text of any lenght. Use SetMaxWidth to limit the maximum text length. If the text is longer than maximun then the text is automatically truncated.
EditField call funtion onChage in case of its text is changed. Event field Msg contains the new text
*/
type EditField struct {
	ControlBase
	// cursor position in edit text
	cursorPos int
	// the number of the first displayed text character - it is used in case of text is longer than edit width
	offset   int
	readonly bool
	maxWidth int

	onChange func(Event)
}

func NewEditField(view View, parent Control, width int, text string, scale int) *EditField {
	e := new(EditField)
	e.onChange = nil
	e.SetTitle(text)
	e.SetEnabled(true)

	if width == AutoSize {
		width = xs.Len(text) + 1
	}

	e.SetSize(width, 1)
	e.cursorPos = xs.Len(text)
	e.offset = 0
	e.parent = parent
	e.view = view
	e.parent = parent
	e.readonly = false

	e.SetConstraints(width, 1)

	e.end()

	if parent != nil {
		parent.AddChild(e, scale)
	}

	return e
}

func (e *EditField) OnChange(fn func(Event)) {
	e.onChange = fn
}

func (e *EditField) SetTitle(title string) {
	if e.title != title {
		e.title = title
		if e.onChange != nil {
			ev := Event{Msg: title, Sender: e}
			go e.onChange(ev)
		}
	}
}

func (e *EditField) Repaint() {
	canvas := e.view.Canvas()

	x, y := e.Pos()
	w, _ := e.Size()

	tm := e.view.Screen().Theme()
	parts := []rune(tm.SysObject(ObjEdit))
	chLeft, chRight := string(parts[0]), string(parts[1])

	var textOut string
	curOff := 0
	if e.offset == 0 && xs.Len(e.title) < e.width {
		textOut = e.title
	} else {
		fromIdx := 0
		toIdx := 0
		if e.offset == 0 {
			toIdx = e.width - 1
			textOut = xs.Slice(e.title, 0, toIdx) + chRight
			curOff = -e.offset
		} else {
			curOff = 1 - e.offset
			fromIdx = e.offset
			if e.width-1 <= xs.Len(e.title)-e.offset {
				toIdx = e.offset + e.width - 2
				textOut = chLeft + xs.Slice(e.title, fromIdx, toIdx) + chRight
			} else {
				textOut = chLeft + xs.Slice(e.title, fromIdx, -1)
			}
		}
	}

	fg, bg := RealColor(tm, e.fg, ColorEditText), RealColor(tm, e.bg, ColorEditBack)
	if !e.Enabled() {
		fg, bg = RealColor(tm, e.fg, ColorDisabledText), RealColor(tm, e.fg, ColorDisabledBack)
	} else if e.Active() {
		fg, bg = RealColor(tm, e.fg, ColorEditActiveText), RealColor(tm, e.bg, ColorEditActiveBack)

	}

	canvas.FillRect(x, y, w, 1, term.Cell{Ch: ' ', Bg: bg})
	canvas.PutText(x, y, textOut, fg, bg)
	if e.active {
		wx, wy := e.view.Pos()
		canvas.SetCursorPos(e.cursorPos+curOff+wx+e.x, wy+e.y)
	}
}

func (e *EditField) insertRune(ch rune) {
	if e.readonly {
		return
	}

	if e.maxWidth > 0 && xs.Len(e.title) >= e.maxWidth {
		return
	}

	idx := e.cursorPos

	if idx == 0 {
		e.SetTitle(string(ch) + e.title)
	} else if idx >= xs.Len(e.title) {
		e.SetTitle(e.title + string(ch))
	} else {
		e.SetTitle(xs.Slice(e.title, 0, idx) + string(ch) + xs.Slice(e.title, idx, -1))
	}

	e.cursorPos++

	if e.cursorPos >= e.width {
		if e.offset == 0 {
			e.offset = 2
		} else {
			e.offset++
		}
	}
}

func (e *EditField) backspace() {
	if e.title == "" || e.cursorPos == 0 || e.readonly {
		return
	}

	length := xs.Len(e.title)
	if e.cursorPos >= length {
		e.cursorPos--
		e.SetTitle(xs.Slice(e.title, 0, length-1))
	} else if e.cursorPos == 1 {
		e.cursorPos = 0
		e.SetTitle(xs.Slice(e.title, 1, -1))
		e.offset = 0
	} else {
		e.cursorPos--
		e.SetTitle(xs.Slice(e.title, 0, e.cursorPos) + xs.Slice(e.title, e.cursorPos+1, -1))
	}

	if length-1 < e.width {
		e.offset = 0
	}
}

func (e *EditField) del() {
	length := xs.Len(e.title)

	if e.title == "" || e.cursorPos == length || e.readonly {
		return
	}

	if e.cursorPos == length-1 {
		e.SetTitle(xs.Slice(e.title, 0, length-1))
	} else {
		e.SetTitle(xs.Slice(e.title, 0, e.cursorPos) + xs.Slice(e.title, e.cursorPos+1, -1))
	}

	if length-1 < e.width {
		e.offset = 0
	}
}

func (e *EditField) charLeft() {
	if e.cursorPos == 0 || e.title == "" {
		return
	}

	if e.cursorPos == e.offset {
		e.offset--
	}

	e.cursorPos--
}

func (e *EditField) charRight() {
	length := xs.Len(e.title)
	if e.cursorPos == length || e.title == "" {
		return
	}

	e.cursorPos++
	if e.cursorPos != length && e.cursorPos >= e.offset+e.width-2 {
		e.offset++
	}
}

func (e *EditField) home() {
	e.offset = 0
	e.cursorPos = 0
}

func (e *EditField) end() {
	length := xs.Len(e.title)
	e.cursorPos = length

	if length < e.width {
		return
	}

	e.offset = length - (e.width - 2)
}

func (e *EditField) Clear() {
	e.home()
	e.SetTitle("")
}

func (e *EditField) ProcessEvent(event Event) bool {
	if !e.Active() || !e.Enabled() {
		return false
	}

	if event.Type == EventActivate && event.X == 0 {
		term.HideCursor()
	}

	if event.Type == EventKey && event.Key != term.KeyTab && event.Key != term.KeyEnter {
		switch event.Key {
		case term.KeySpace:
			e.insertRune(' ')
			return true
		case term.KeyBackspace:
			e.backspace()
			return true
		case term.KeyDelete:
			e.del()
			return true
		case term.KeyArrowLeft:
			e.charLeft()
			return true
		case term.KeyHome:
			e.home()
			return true
		case term.KeyEnd:
			e.end()
			return true
		case term.KeyCtrlR:
			if !e.readonly {
				e.Clear()
			}
			return true
		case term.KeyArrowRight:
			e.charRight()
			return true
		default:
			if event.Ch != 0 {
				e.insertRune(event.Ch)
				return true
			}
		}
		return false
	}

	return false
}

func (e *EditField) SetMaxWidth(w int) {
	e.maxWidth = w
	if w > 0 && xs.Len(e.title) > w {
		e.title = xs.Slice(e.title, 0, w)
		e.end()
	}
}

func (e *EditField) GetMaxWidth() int {
	return e.maxWidth
}
