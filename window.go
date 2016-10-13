package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
)

// Window is an implemetation of View managed by Composer.
type Window struct {
	BaseControl

	buttons   ViewButton
	maximized bool
	// maximization support
	origWidth  int
	origHeight int
	origX      int
	origY      int

	onClose func(Event)
}

func CreateWindow(x, y, w, h int, title string) *Window {
	wnd := new(Window)
	if w == AutoSize || w < 1 || w > 1000 {
		w = 10
	}
	if h == AutoSize || h < 1 || h > 1000 {
		w = 5
	}

	wnd.SetConstraints(w, h)
	wnd.SetSize(w, h)
	wnd.SetPos(x, y)
	wnd.SetTitle(title)
	wnd.buttons = ButtonClose | ButtonBottom | ButtonMaximize
	wnd.children = make([]Control, 0)
	wnd.SetPaddings(1, 1)
	wnd.SetGaps(1, 0)
	wnd.SetScale(1)

	return wnd
}

func (wnd *Window) buttonCount() int {
	cnt := 0
	if wnd.buttons&ButtonClose == ButtonClose {
		cnt += 1
	}
	if wnd.buttons&ButtonMaximize == ButtonMaximize {
		cnt += 1
	}
	if wnd.buttons&ButtonBottom == ButtonBottom {
		cnt += 1
	}

	return cnt
}

func (wnd *Window) drawFrame() {
	PushAttributes()
	defer PopAttributes()

	var bs BorderStyle
	if wnd.inactive {
		bs = BorderThin
	} else {
		bs = BorderThick
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

	chars := []rune(SysObject(ObjViewButtons))
	cMax, cBottom, cClose, cOpenB, cCloseB := chars[0], chars[1], chars[2], chars[3], chars[4]

	pos := wnd.x + wnd.width - btnCount - 2
	putCharUnsafe(pos, wnd.y, cOpenB)
	pos += 1
	if wnd.buttons&ButtonBottom == ButtonBottom {
		putCharUnsafe(pos, wnd.y, cBottom)
		pos += 1
	}
	if wnd.buttons&ButtonMaximize == ButtonMaximize {
		putCharUnsafe(pos, wnd.y, cMax)
		pos += 1
	}
	if wnd.buttons&ButtonClose == ButtonClose {
		putCharUnsafe(pos, wnd.y, cClose)
		pos += 1
	}
	putCharUnsafe(pos, wnd.y, cCloseB)
}

func (wnd *Window) Draw() {
	PushAttributes()
	defer PopAttributes()

	fg, bg := RealColor(wnd.fg, ColorViewText), RealColor(wnd.bg, ColorViewBack)
	SetBackColor(bg)

	FillRect(wnd.x, wnd.y, wnd.width, wnd.height, ' ')

	wnd.DrawChildren()

	SetBackColor(bg)
	SetTextColor(fg)

	wnd.drawFrame()
	wnd.drawTitle()
	wnd.drawButtons()
}

func (c *Window) HitTest(x, y int) HitResult {
	if x > c.x && x < c.x+c.width-1 &&
		y > c.y && y < c.y+c.height-1 {
		return HitInside
	}

	if x == c.x && y == c.y {
		return HitTopLeft
	}

	if x == c.x+c.width-1 && y == c.y {
		return HitTopRight
	}

	if x == c.x && y == c.y+c.height-1 {
		return HitBottomLeft
	}

	if x == c.x+c.width-1 && y == c.y+c.height-1 {
		return HitBottomRight
	}

	if x == c.x && y > c.y && y < c.y+c.height-1 {
		return HitLeft
	}

	if x == c.x+c.width-1 && y > c.y && y < c.y+c.height-1 {
		return HitRight
	}

	if y == c.y && x > c.x && x < c.x+c.width-1 {
		btnCount := c.buttonCount()
		if x < c.x+c.width-1-btnCount {
			return HitTop
		}

		hitRes := []HitResult{HitTop, HitTop, HitTop}
		pos := 0

		if c.buttons&ButtonBottom == ButtonBottom {
			hitRes[pos] = HitButtonBottom
			pos += 1
		}
		if c.buttons&ButtonMaximize == ButtonMaximize {
			hitRes[pos] = HitButtonMaximize
			pos += 1
		}
		if c.buttons&ButtonClose == ButtonClose {
			hitRes[pos] = HitButtonClose
			pos += 1
		}

		return hitRes[x-(c.x+c.width-1-btnCount)]
	}

	if y == c.y+c.height-1 && x > c.x && x < c.x+c.width-1 {
		return HitBottom
	}

	return HitOutside
}

func (c *Window) ProcessEvent(ev Event) bool {
	switch ev.Type {
	case EventMove:
		c.PlaceChildren()
	case EventResize:
		c.ResizeChildren()
		c.PlaceChildren()
	case EventClose:
		if c.onClose != nil {
			c.onClose(ev)
		}
	case EventKey:
		if ev.Key == term.KeyTab {
			aC := ActiveControl(c)
			nC := NextControl(c, aC, true)
			if nC != aC {
				if aC != nil {
					aC.SetActive(false)
					aC.ProcessEvent(Event{Type: EventActivate, X: 0})
				}
				if nC != nil {
					nC.SetActive(true)
					nC.ProcessEvent(Event{Type: EventActivate, X: 1})
				}
			}
			return true
		} else {
			return SendEventToChild(c, ev)
		}
	default:
		if ev.Type == EventMouse && ev.Key == term.MouseLeft {
			DeactivateControls(c)
		}
		return SendEventToChild(c, ev)
	}

	return false
}

func (w *Window) OnClose(fn func(Event)) {
	w.onClose = fn
}

// SetMaximized opens the view to full screen or restores its
// previous size
func (w *Window) SetMaximized(maximize bool) {
	if maximize == w.maximized {
		return
	}

	if maximize {
		w.origX, w.origY = w.Pos()
		w.origWidth, w.origHeight = w.Size()
		w.maximized = true
		w.SetPos(0, 0)
		width, height := ScreenSize()
		w.SetSize(width, height)
	} else {
		w.maximized = false
		w.SetPos(w.origX, w.origY)
		w.SetSize(w.origWidth, w.origHeight)
	}
	w.ResizeChildren()
	w.PlaceChildren()
}

// Maximized returns if the view is in full screen mode
func (w *Window) Maximized() bool {
	return w.maximized
}
