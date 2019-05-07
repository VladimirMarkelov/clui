package clui

import (
	мИнт "./пакИнтерфейсы"
	мКнст "./пакКонстанты"
	мСоб "./пакСобытия"
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
)

// Window is an implemetation of View managed by Composer.
type Window struct {
	BaseControl

	buttons   мКнст.ViewButton
	maximized bool
	// maximization support
	origWidth  int
	origHeight int
	origX      int
	origY      int
	hidden     bool
	immovable  bool
	fixedSize  bool
	border     мКнст.BorderStyle

	onClose        func(мИнт.ИСобытие) bool
	onScreenResize func(мИнт.ИСобытие)

	onKeyDown *keyDownCb
}

type keyDownCb struct {
	data interface{}
	fn   func(evt мИнт.ИСобытие, data interface{}) bool
}

//CreateWindow --
func CreateWindow(x, y, w, h int, title string) *Window {
	wnd := new(Window)
	wnd.BaseControl = NewBaseControl()

	if w == мКнст.AutoSize || w < 1 || w > 1000 {
		w = 10
	}
	if h == мКнст.AutoSize || h < 1 || h > 1000 {
		w = 5
	}

	wnd.SetConstraints(w, h)
	wnd.SetSize(w, h)
	wnd.SetPos(x, y)
	wnd.SetTitle(title)
	wnd.buttons = мКнст.ButtonClose | мКнст.ButtonBottom | мКнст.ButtonMaximize
	wnd.children = make([]Control, 0)
	wnd.SetPaddings(1, 1)
	wnd.SetGaps(1, 0)
	wnd.SetScale(1)
	wnd.SetBorder(мКнст.BorderAuto)

	return wnd
}

func (wnd *Window) buttonCount() int {
	cnt := 0
	if wnd.buttons&мКнст.ButtonClose == мКнст.ButtonClose {
		cnt++
	}
	if wnd.buttons&мКнст.ButtonMaximize == мКнст.ButtonMaximize {
		cnt++
	}
	if wnd.buttons&мКнст.ButtonBottom == мКнст.ButtonBottom {
		cnt++
	}

	return cnt
}

func (wnd *Window) drawFrame() {
	PushAttributes()
	defer PopAttributes()

	var bs мКнст.BorderStyle
	if wnd.border == мКнст.BorderAuto {
		if wnd.inactive {
			bs = мКнст.BorderThin
		} else {
			bs = мКнст.BorderThick
		}
	} else if wnd.border == мКнст.BorderNone {
	} else {
		bs = wnd.border
	}

	DrawFrame(wnd.x, wnd.y, wnd.width, wnd.height, bs)
}

func (wnd *Window) drawTitle() {
	PushAttributes()
	defer PopAttributes()

	btnCount := wnd.buttonCount()
	maxw := wnd.width - 2 - btnCount
	if btnCount > 0 {
		maxw -= 2
	}

	fitTitle := wnd.title
	rawText := UnColorizeText(fitTitle)
	if xs.Len(rawText) > maxw {
		fitTitle = SliceColorized(fitTitle, 0, maxw-3) + "..."
	}

	DrawText(wnd.x+1, wnd.y, fitTitle)
}

func (wnd *Window) drawButtons() {
	btnCount := wnd.buttonCount()
	if btnCount == 0 {
		return
	}

	PushAttributes()
	defer PopAttributes()

	chars := []rune(SysObject(мКнст.ObjViewButtons))
	cMax, cBottom, cClose, cOpenB, cCloseB := chars[0], chars[1], chars[2], chars[3], chars[4]

	pos := wnd.x + wnd.width - btnCount - 2
	putCharUnsafe(pos, wnd.y, cOpenB)
	pos++
	if wnd.buttons&мКнст.ButtonBottom == мКнст.ButtonBottom {
		putCharUnsafe(pos, wnd.y, cBottom)
		pos++
	}
	if wnd.buttons&мКнст.ButtonMaximize == мКнст.ButtonMaximize {
		putCharUnsafe(pos, wnd.y, cMax)
		pos++
	}
	if wnd.buttons&мКнст.ButtonClose == мКнст.ButtonClose {
		putCharUnsafe(pos, wnd.y, cClose)
		pos++
	}
	putCharUnsafe(pos, wnd.y, cCloseB)
}

//Draw --
func (wnd *Window) Draw() {
	WindowManager().BeginUpdate()
	defer WindowManager().EndUpdate()
	PushAttributes()
	defer PopAttributes()

	fg, bg := RealColor(wnd.fg, wnd.Style(), мКнст.ColorViewText), RealColor(wnd.bg, wnd.Style(), мКнст.ColorViewBack)
	SetBackColor(bg)

	FillRect(wnd.x, wnd.y, wnd.width, wnd.height, ' ')

	wnd.DrawChildren()

	SetBackColor(bg)
	SetTextColor(fg)

	wnd.drawFrame()
	wnd.drawTitle()
	wnd.drawButtons()
}

//HitTest --
func (wnd *Window) HitTest(x, y int) мИнт.HitResult {
	if x > wnd.x && x < wnd.x+wnd.width-1 &&
		y > wnd.y && y < wnd.y+wnd.height-1 {
		return мКнст.HitInside
	}

	hResult := мКнст.HitOutside

	if x == wnd.x && y == wnd.y {
		hResult = мКнст.HitTopLeft
	} else if x == wnd.x+wnd.width-1 && y == wnd.y {
		hResult = мКнст.HitTopRight
	} else if x == wnd.x && y == wnd.y+wnd.height-1 {
		hResult = мКнст.HitBottomLeft
	} else if x == wnd.x+wnd.width-1 && y == wnd.y+wnd.height-1 {
		hResult = мКнст.HitBottomRight
	} else if x == wnd.x && y > wnd.y && y < wnd.y+wnd.height-1 {
		hResult = мКнст.HitLeft
	} else if x == wnd.x+wnd.width-1 && y > wnd.y && y < wnd.y+wnd.height-1 {
		hResult = мКнст.HitRight
	} else if y == wnd.y && x > wnd.x && x < wnd.x+wnd.width-1 {
		btnCount := wnd.buttonCount()
		if x < wnd.x+wnd.width-1-btnCount {
			hResult = мКнст.HitTop
		} else {
			hitRes := []мКнст.HitResult{мКнст.HitTop, мКнст.HitTop, мКнст.HitTop}
			pos := 0

			if wnd.buttons&мКнст.ButtonBottom == мКнст.ButtonBottom {
				hitRes[pos] = мКнст.HitButtonBottom
				pos++
			}
			if wnd.buttons&мКнст.ButtonMaximize == мКнст.ButtonMaximize {
				hitRes[pos] = мКнст.HitButtonMaximize
				pos++
			}
			if wnd.buttons&мКнст.ButtonClose == мКнст.ButtonClose {
				hitRes[pos] = мКнст.HitButtonClose
				pos++
			}

			hResult = hitRes[x-(wnd.x+wnd.width-1-btnCount)]
		}
	} else if y == wnd.y+wnd.height-1 && x > wnd.x && x < wnd.x+wnd.width-1 {
		hResult = мКнст.HitBottom
	}

	if hResult != мКнст.HitOutside {
		if wnd.immovable && hResult == мКнст.HitTop {
			hResult = мКнст.HitInside
		}
		if wnd.fixedSize &&
			(hResult == мКнст.HitBottom || hResult == мКнст.HitLeft ||
				hResult == мКнст.HitRight || hResult == мКнст.HitTopLeft ||
				hResult == мКнст.HitTopRight || hResult == мКнст.HitBottomRight ||
				hResult == мКнст.HitBottomLeft || hResult == мКнст.HitButtonMaximize) {
			hResult = мКнст.HitInside
		}
	}

	return hResult
}

//ProcessEvent --
func (wnd *Window) ProcessEvent(ev мИнт.ИСобытие) bool {
	switch ev.Type {
	case мКнст.EventMove:
		wnd.PlaceChildren()
	case мКнст.EventResize:
		wnd.ResizeChildren()
		wnd.PlaceChildren()
	case мКнст.EventClose:
		if wnd.onClose != nil {
			if !wnd.onClose(ev) {
				return false
			}
		}
		return true
	case мКнст.EventKey:
		if ev.Key == term.KeyTab || ev.Key == term.KeyArrowUp || ev.Key == term.KeyArrowDown {
			if SendEventToChild(wnd, ev) {
				return true
			}

			aC := ActiveControl(wnd)
			nC := NextControl(wnd, aC, ev.Key != term.KeyArrowUp)

			var clipped Control

			if aC != nil && aC.Clipped() {
				clipped = aC
			} else if nC != nil {
				clipped = ClippedParent(nC)
			}

			if clipped != nil {
				dir := 1
				if ev.Key != term.KeyArrowUp {
					dir = -1
				}

				clipped.ProcessEvent(мКнст.Event{Type: мКнст.EventActivateChild, Target: nC, X: dir})
			}

			if nC != aC {
				if aC != nil {
					aC.SetActive(false)
					aC.ProcessEvent(мКнст.Event{Type: мКнст.EventActivate, X: 0})
				}
				if nC != nil {
					aC.SetActive(true)
					aC.ProcessEvent(мКнст.Event{Type: мКнст.EventActivate, X: 1})
				}
			}
			return true
		}
		if SendEventToChild(wnd, ev) {
			return true
		}
		if wnd.onKeyDown != nil {
			return wnd.onKeyDown.fn(ev, wnd.onKeyDown.data)
		}
		return false

	default:
		if ev.Type == мКнст.EventMouse && ev.Key == term.MouseLeft {
			DeactivateControls(wnd)
		}
		return SendEventToChild(wnd, ev)
	}

	return false
}

// OnClose sets the callback that is called when the Window is about to destroy
func (wnd *Window) OnClose(fn func(мИнт.ИСобытие) bool) {
	wnd.onClose = fn
}

// OnKeyDown sets the callback that is called when a user presses a key
// while the Window is active
func (wnd *Window) OnKeyDown(fn func(мИнт.ИСобытие, interface{}) bool, data interface{}) {
	if fn == nil {
		wnd.onKeyDown = nil
	} else {
		wnd.onKeyDown = &keyDownCb{data: data, fn: fn}
	}
}

// OnScreenResize sets the callback that is called when size of terminal changes
func (wnd *Window) OnScreenResize(fn func(мИнт.ИСобытие)) {
	wnd.onScreenResize = fn
}

// Border returns the default window border
func (wnd *Window) Border() мКнст.BorderStyle {
	return wnd.border
}

// SetBorder changes the default window border
func (wnd *Window) SetBorder(border мКнст.BorderStyle) {
	wnd.border = border
}

// SetMaximized opens the view to full screen or restores its
// previous size
func (wnd *Window) SetMaximized(maximize bool) {
	if maximize == wnd.maximized {
		return
	}

	if maximize {
		wnd.origX, wnd.origY = wnd.Pos()
		wnd.origWidth, wnd.origHeight = wnd.Size()
		wnd.maximized = true
		wnd.SetPos(0, 0)
		width, height := ScreenSize()
		wnd.SetSize(width, height)
	} else {
		wnd.maximized = false
		wnd.SetPos(wnd.origX, wnd.origY)
		wnd.SetSize(wnd.origWidth, wnd.origHeight)
	}
	wnd.ResizeChildren()
	wnd.PlaceChildren()
}

// Maximized returns if the view is in full screen mode
func (wnd *Window) Maximized() bool {
	return wnd.maximized
}

// Visible returns if the window must be drawn on the screen
func (wnd *Window) Visible() bool {
	return !wnd.hidden
}

// SetVisible allows to temporarily remove the window from screen
// and show it later without reconstruction
func (wnd *Window) SetVisible(visible bool) {
	if wnd.hidden == visible {
		wnd.hidden = !visible
		if wnd.hidden {
			wnd.SetModal(false)
			if WindowManager().topWindow() == wnd {
				WindowManager().moveActiveWindowToBottom()
			}
		} else {
			WindowManager().activateWindow(wnd)
		}
	}
}

// Movable returns if the Window can be moved with mouse or keyboard
func (wnd *Window) Movable() bool {
	return !wnd.immovable
}

// Sizable returns if size of the Window can be changed with mouse or keyboard
func (wnd *Window) Sizable() bool {
	return !wnd.fixedSize
}

// SetMovable turns on and off ability to change Window position with mouse
// or keyboard
func (wnd *Window) SetMovable(movable bool) {
	wnd.immovable = !movable
}

// SetSizable turns on and off ability to change Window size with mouse
// or keyboard
func (wnd *Window) SetSizable(sizable bool) {
	wnd.fixedSize = !sizable
}

// TitleButtons returns a set of buttons shown in the Window title bar
func (wnd *Window) TitleButtons() мКнст.ViewButton {
	return wnd.buttons
}

// SetTitleButtons sets the title bar buttons available for a user
func (wnd *Window) SetTitleButtons(buttons мКнст.ViewButton) {
	wnd.buttons = buttons
}
