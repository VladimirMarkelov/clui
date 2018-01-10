package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	"strings"
)

// OnChange sets the callback that is called when EditField content is changed
func (e *EditField) OnChange(fn func(Event)) {
	e.onChange = fn
}

// OnKeyPress sets the callback that is called when a user presses a Key while
// the controls is active. If a handler processes the key it should return
// true. If handler returns false it means that the default handler will
// process the key
func (e *EditField) OnKeyPress(fn func(term.Key) bool) {
	e.onKeyPress = fn
}

// SetTitle changes the EditField content and emits OnChage eventif the new value does not equal to old one
func (e *EditField) SetTitle(title string) {
	if e.title != title {
		e.title = title
		if e.onChange != nil {
			ev := Event{Msg: title}
			go e.onChange(ev)
		}
	}
}

// Repaint draws the control on its View surface
func (e *EditField) Draw() {
	PushAttributes()
	defer PopAttributes()

	x, y := e.Pos()
	w, _ := e.Size()

	parts := []rune(SysObject(ObjEdit))
	chLeft, chRight := string(parts[0]), string(parts[1])
	chStar := "*"
	if len(parts) > 3 {
		chStar = string(parts[3])
	}

	var textOut string
	curOff := 0
	if e.offset == 0 && xs.Len(e.title) < e.width {
		if e.showStars {
			textOut = strings.Repeat(chStar, xs.Len(e.title))
		} else {
			textOut = e.title
		}
	} else {
		fromIdx := 0
		toIdx := 0
		if e.offset == 0 {
			toIdx = e.width - 1
			if e.showStars {
				textOut = strings.Repeat(chStar, toIdx) + chRight
			} else {
				textOut = xs.Slice(e.title, 0, toIdx) + chRight
			}
			curOff = -e.offset
		} else {
			curOff = 1 - e.offset
			fromIdx = e.offset
			if e.width-1 <= xs.Len(e.title)-e.offset {
				toIdx = e.offset + e.width - 2
				if e.showStars {
					textOut = chLeft + strings.Repeat(chStar, toIdx-fromIdx) + chRight
				} else {
					textOut = chLeft + xs.Slice(e.title, fromIdx, toIdx) + chRight
				}
			} else {
				if e.showStars {
					textOut = chLeft + strings.Repeat(chStar, xs.Len(e.title)-fromIdx)
				} else {
					textOut = chLeft + xs.Slice(e.title, fromIdx, -1)
				}
			}
		}
	}

	fg, bg := RealColor(e.fg, ColorEditText), RealColor(e.bg, ColorEditBack)
	if !e.Enabled() {
		fg, bg = RealColor(e.fg, ColorDisabledText), RealColor(e.fg, ColorDisabledBack)
	} else if e.Active() {
		fg, bg = RealColor(e.fg, ColorEditActiveText), RealColor(e.bg, ColorEditActiveBack)
	}

	SetTextColor(fg)
	SetBackColor(bg)
	FillRect(x, y, w, 1, ' ')
	DrawRawText(x, y, textOut)
	if e.Active() {
		SetCursorPos(e.cursorPos+e.x+curOff, e.y)
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

// Clear empties the EditField and emits OnChange event
func (e *EditField) Clear() {
	e.home()
	e.SetTitle("")
}

// SetMaxWidth sets the maximum lenght of the EditField text. If the current text is longer it is truncated
func (e *EditField) SetMaxWidth(w int) {
	e.maxWidth = w
	if w > 0 && xs.Len(e.title) > w {
		e.title = xs.Slice(e.title, 0, w)
		e.end()
	}
}

// MaxWidth returns the current maximum text length. Zero means no limit
func (e *EditField) MaxWidth() int {
	return e.maxWidth
}

// SetSize changes control size. Constant DoNotChange can be
// used as placeholder to indicate that the control attrubute
// should be unchanged.
// Method does nothing if new size is less than minimal size
// EditField height cannot be changed - it equals 1 always
func (e *EditField) SetSize(width, height int) {
	if width != KeepValue && (width > 1000 || width < e.minW) {
		return
	}
	if height != KeepValue && (height > 200 || height < e.minH) {
		return
	}

	if width != KeepValue {
		e.width = width
	}

	e.height = 1
}

// PasswordMode returns whether password mode is enabled for the control
func (e *EditField) PasswordMode() bool {
	return e.showStars
}

// SetPasswordMode changes the way an EditField displays it content.
// If PasswordMode is false then the EditField works as regular text entry
// control. If PasswordMode is true then the EditField shows its content hidden
// with star characters ('*' by default)
func (e *EditField) SetPasswordMode(pass bool) {
	e.showStars = pass
}
