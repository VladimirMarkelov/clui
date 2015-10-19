package clui

import (
	"github.com/nsf/termbox-go"
	"log"
	"os"
)

// An service object that manages Views and console, processes events, and provides service methods
// One application must have only one object of this type
type Composer struct {
	// list of visible Views
	views []View
	// console width and height
	width, height int
	// console canvas
	canvas Canvas
	// a channel to communicate with View(e.g, Views send redraw event to this channel)
	channel chan Event

	// current color scheme
	themeManager *ThemeManager

	// multi key sequences support. The flag below are true if the last keyboard combination was Ctrl+S or Ctrl+W respectively
	ctrlKey termbox.Key
	// last pressed key - to make repeatable actions simpler, e.g, at first one presses Ctrl+S and then just repeatedly presses arrow lest to resize View
	lastKey termbox.Key

	//debug
	logger *log.Logger
}

func (c *Composer) initBuffer() {
	c.canvas = NewFrameBuffer(c.width, c.height)
	c.canvas.Clear(ColorBlack)
}

// Initialize library and starts console management
func InitLibrary() *Composer {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	c := new(Composer)
	c.ctrlKey = termbox.KeyEsc

	termbox.HideCursor()

	file, _ := os.OpenFile("debugui.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	c.logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	c.logger.Printf("----------------------------------")

	c.channel = make(chan Event)

	c.views = make([]View, 0)

	c.width, c.height = termbox.Size()
	c.initBuffer()

	c.themeManager = NewThemeManager()

	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)

	c.redrawAll()

	return c
}

func (c *Composer) redrawAll() {
	for j := 0; j < c.height; j++ {
		for i := 0; i < c.width; i++ {
			sym, ok := c.canvas.Symbol(i, j)
			if ok {
				termbox.SetCell(i, j, sym.Ch, termbox.Attribute(sym.Fg), termbox.Attribute(sym.Bg))
			}
		}
	}

	termbox.Flush()
}

// Repaints all View on the screen. Now the method is not efficient: at first clears a console and then draws all Views starting from the bottom
func (c *Composer) RefreshScreen(invalidate bool) {
	c.canvas.Clear(ColorBlack)

	for _, wnd := range c.views {
		if invalidate {
			wnd.Repaint()
		}
		wnd.Draw(c.canvas)
	}

	c.redrawAll()
}

func (c *Composer) CreateView(posX, posY, width, height int, title string) View {
	view := NewWindow(c, posX, posY, width, height, title)

	c.views = append(c.views, view)
	view.Repaint()

	c.activateView(view)

	c.RefreshScreen(false)

	return view
}

func (c *Composer) checkWindowUnderMouse(screenX, screenY int) (View, HitResult) {
	if len(c.views) == 0 {
		return nil, HitOutside
	}

	for i := len(c.views) - 1; i >= 0; i-- {
		window := c.views[i]
		hit := window.HitTest(screenX, screenY)
		if hit != HitOutside {
			return window, hit
		}
	}

	return nil, HitOutside
}

func (c *Composer) activateView(view View) {
	if c.topView() == view {
		for _, v := range c.views {
			v.SetActive(false)
		}
		view.SetActive(true)
		return
	}

	wList := make([]View, 0)
	found := false

	for _, v := range c.views {
		if v != view {
			v.SetActive(false)
			wList = append(wList, v)
		} else {
			found = true
		}
	}

	if !found {
		panic("Invalid view to activate")
	}

	view.SetActive(true)
	c.views = append(wList, view)
}

func (c *Composer) moveActiveWindowToBottom() bool {
	if len(c.views) < 2 {
		return false
	}

	event := Event{Type: EventActivate, X: 0} // send deactivated
	c.sendEventToActiveView(event)

	last := c.topView()

	for i := len(c.views) - 1; i > 0; i-- {
		c.views[i] = c.views[i-1]
	}

	c.views[0] = last
	c.activateView(c.topView())

	event = Event{Type: EventActivate, X: 1} // send 'activated'
	c.sendEventToActiveView(event)
	c.RefreshScreen(true)

	return true
}

func (c *Composer) termboxEventToLocal(ev termbox.Event) Event {
	e := Event{Type: EventType(ev.Type), Ch: ev.Ch, Key: ev.Key, Err: ev.Err, X: ev.MouseX, Y: ev.MouseY, Mod: ev.Mod}
	return e
}

func (c *Composer) sendEventToActiveView(ev Event) bool {
	view := c.topView()
	if view != nil {
		return view.ProcessEvent(ev)
	}

	return false
}

func (c *Composer) topView() View {
	if len(c.views) == 0 {
		return nil
	}

	return c.views[len(c.views)-1]
}

func (c *Composer) resizeTopView(ev termbox.Event) bool {
	view := c.topView()
	if view == nil {
		return false
	}

	w, h := view.Size()
	w1, h1 := w, h
	minW, minH := view.Constraints()
	if ev.Key == termbox.KeyArrowUp && minH < h {
		h--
	} else if ev.Key == termbox.KeyArrowLeft && minW < w {
		w--
	} else if ev.Key == termbox.KeyArrowDown {
		h++
	} else if ev.Key == termbox.KeyArrowRight {
		w++
	}

	if w1 != w || h1 != h {
		view.SetSize(w, h)
		// event := Event{Type: EventResize, X: w, Y: h}
		// c.sendEventToActiveView(event)
		c.RefreshScreen(true)
	}

	return true
}

func (c *Composer) moveTopView(ev termbox.Event) bool {
	if len(c.views) > 0 {
		view := c.topView()
		if view != nil {
			x, y := view.Pos()
			w, h := view.Size()
			x1, y1 := x, y
			cx, cy := termbox.Size()
			if ev.Key == termbox.KeyArrowUp && y > 0 {
				y--
			} else if ev.Key == termbox.KeyArrowDown && y+h < cy {
				y++
			} else if ev.Key == termbox.KeyArrowLeft && x > 0 {
				x--
			} else if ev.Key == termbox.KeyArrowRight && x+w < cx {
				x++
			}

			if x1 != x || y1 != y {
				view.SetPos(x, y)
				// event := Event{Type: EventMove, X: x, Y: y}
				// c.sendEventToActiveView(event)
				c.RefreshScreen(true)
			}
		}
		return true
	} else {
		return false
	}
}

func (c *Composer) isDeadKey(ev termbox.Event) bool {
	if ev.Key == termbox.KeyCtrlS || ev.Key == termbox.KeyCtrlW {
		c.ctrlKey = ev.Key
		c.lastKey = termbox.KeyEsc
		return true
	}

	c.ctrlKey = termbox.KeyEsc
	return false
}

func (c *Composer) processKeySeq(ev termbox.Event) bool {
	if c.ctrlKey == termbox.KeyEsc {
		return false
	}

	if c.ctrlKey == termbox.KeyCtrlS {
		if c.lastKey == termbox.KeyEsc {
			c.lastKey = ev.Key
		} else if c.lastKey != ev.Key {
			c.ctrlKey = termbox.KeyEsc
			return false
		}

		switch ev.Key {
		case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight:
			evCopy := ev
			c.resizeTopView(evCopy)
			return true
		}

		return false
	}

	if c.ctrlKey == termbox.KeyCtrlW {
		if c.lastKey == termbox.KeyEsc {
			c.lastKey = ev.Key
		} else if c.lastKey != ev.Key {
			c.ctrlKey = termbox.KeyEsc
			return false
		}

		switch ev.Key {
		case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight:
			evCopy := ev
			c.moveTopView(evCopy)
			return true
		case termbox.KeyCtrlH:
			// if len(c.windowOrder) > 1 && ev.Mod&termbox.ModControl != 0 {
			//  c.moveActiveWindowToBottom()
			// }
			return true
		}

		return false
	}

	c.ctrlKey = termbox.KeyEsc
	return true
}

func (c *Composer) processKey(ev termbox.Event) bool {
	if c.processKeySeq(ev) {
		return false
	}
	if c.isDeadKey(ev) {
		return false
	}

	switch ev.Key {
	case termbox.KeyCtrlQ:
		return true
	// case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight:
	//  if c.sendEventToActiveView(c.termboxEventToLocal(ev)) {
	//      c.RefreshScreen()
	//  }
	// case termbox.KeyEnd:
	//  if c.sendEventToActiveView(c.termboxEventToLocal(ev)) {
	//      c.RefreshScreen()
	//  }
	default:
		if c.sendEventToActiveView(c.termboxEventToLocal(ev)) {
			c.topView().Repaint()
			c.RefreshScreen(false)
		}
	}

	return false
}

func (c *Composer) processMouseClick(ev termbox.Event) {
	view, hit := c.checkWindowUnderMouse(ev.MouseX, ev.MouseY)

	if view == nil {
		return
	}

	if c.topView() != view {
		event := Event{Type: EventActivate, X: 0} // send 'deactivated'
		c.sendEventToActiveView(event)
		c.activateView(view)
		event = Event{Type: EventActivate, X: 1} // send 'activated'
		c.sendEventToActiveView(event)
		c.RefreshScreen(true)
	} else if hit == HitInside {
		c.sendEventToActiveView(c.termboxEventToLocal(ev))
		c.RefreshScreen(true)
	} else if hit == HitButtonClose {
		if len(c.views) > 1 {
			event := Event{Type: EventClose}
			c.sendEventToActiveView(event)

			c.DestroyWindow(view)
			activate := c.topView()
			c.activateView(activate)
			event = Event{Type: EventActivate, X: 1} // send 'activated'
			c.sendEventToActiveView(event)

			c.RefreshScreen(true)
		}
	} else if hit == HitButtonBottom {
		c.moveActiveWindowToBottom()
	}
}

// func (c *Composer) moveAndResize(dx, dy int, wnd Window) bool {
//  if dx == 0 && dy == 0 {
//      return false
//  }
//
//  posX, posY := wnd.GetPos()
//  w, h := wnd.GetSize()
//  mnX, mnY := wnd.GetConstraints()
//  newW, newH := w, h
//  newX, newY := posX, posY
//
//  if c.dragHit == HitLeftBorder {
//      newX = posX + dx
//      newW = w - dx
//  } else if c.dragHit == HitRightBorder {
//      newW = w + dx
//  } else if c.dragHit == HitBottomBorder {
//      newH = h + dy
//  } else if c.dragHit == HitTopLeft {
//      newX, newY = posX+dx, posY+dy
//      newW, newH = w-dx, h-dy
//  } else if c.dragHit == HitBottomLeft {
//      newX = posX + dx
//      newW, newH = w-dx, h+dy
//  } else if c.dragHit == HitTopRight {
//      newY = posY + dy
//      newW, newH = w+dx, h-dy
//  } else if c.dragHit == HitBottomRight {
//      newW, newH = w+dx, h+dy
//  }
//
//  if (mnX != -1 && newW < mnX) || (mnY != -1 && newH < mnY) {
//      return false
//  }
//
//  wnd.SetPos(newX, newY)
//  wnd.SetSize(newW, newH)
//
//  if posX != newX || posY != newY {
//      event := Event{Type: EventMove, X: newX, Y: newY}
//      c.sendEventToActiveView(event)
//  }
//  if newW != w || newH != h {
//      event := Event{Type: EventResize, X: w, Y: h}
//      c.sendEventToActiveView(event)
//  }
//
//  wnd.Redraw()
//
//  return true
// }

// Asks a Composer to stops console management and quit application
func (c *Composer) Stop() {
	ev := Event{Type: EventQuit}
	go c.PutEvent(ev)
}

// Main event loop
func (c *Composer) MainLoop() {
	// c.redrawAll()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventQueue:
			switch ev.Type {
			case termbox.EventKey:
				if c.processKey(ev) {
					return
				}
			case termbox.EventMouse:
				c.processMouseClick(ev)
			case termbox.EventError:
				panic(ev.Err)
				// case termbox.EventResize:
				//  termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				//  c.width, c.height = termbox.Size()
				//  c.screen = c.initBuffer(c.width, c.height)
				//  c.RefreshScreen()
			}
		case cmd := <-c.channel:
			if cmd.Type == EventRedraw {
				c.RefreshScreen(true)
			} else if cmd.Type == EventQuit {
				return
			}
		}
	}
}

// Send event to a Composer. Used by Windows to ask for repainting or for quitting the application
func (c *Composer) PutEvent(ev Event) {
	c.channel <- ev
}

// Closes console management and makes a console cursor visible
func (c *Composer) Close() {
	termbox.SetCursor(3, 3)
	termbox.Close()
}

// Shows consolse cursor at given position. Setting cursor to -1,-1 hides cursor
func (c *Composer) SetCursorPos(x, y int) {
	termbox.SetCursor(x, y)
}

func (c *Composer) DestroyWindow(view View) {
	ev := Event{Type: EventClose}
	c.sendEventToActiveView(ev)

	newOrder := make([]View, 0)
	for i := 0; i < len(c.views); i++ {
		if c.views[i] != view {
			newOrder = append(newOrder, c.views[i])
		}
	}
	c.views = newOrder
}

func (c *Composer) Theme() Theme {
	return c.themeManager
}

func (c *Composer) Logger() *log.Logger {
	return c.logger
}
