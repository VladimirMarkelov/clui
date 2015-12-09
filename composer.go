package clui

import (
	term "github.com/nsf/termbox-go"
	"log"
	"os"
)

// Composer is a service object that manages Views and console, processes
// events, and provides service methods. One application must have only
// one object of this type
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
	ctrlKey term.Key
	// last pressed key - to make repeatable actions simpler, e.g, at first one presses Ctrl+S and then just repeatedly presses arrow lest to resize View
	lastKey term.Key

	//debug
	logger *log.Logger
}

func (c *Composer) initBuffer() {
	c.canvas = NewFrameBuffer(c.width, c.height)
	c.canvas.Clear(ColorBlack)
}

// InitLibrary initializes library and starts console management.
// Retuns nil in case of error
func InitLibrary() *Composer {
	err := term.Init()
	if err != nil {
		return nil
	}

	c := new(Composer)
	c.ctrlKey = term.KeyEsc

	term.HideCursor()

	file, _ := os.OpenFile("debugui.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	c.logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	c.logger.Printf("----------------------------------")

	c.channel = make(chan Event)

	c.views = make([]View, 0)

	c.width, c.height = term.Size()
	c.initBuffer()

	c.themeManager = NewThemeManager()

	term.SetInputMode(term.InputAlt | term.InputMouse)

	c.redrawAll()

	return c
}

func (c *Composer) redrawAll() {
	for j := 0; j < c.height; j++ {
		for i := 0; i < c.width; i++ {
			sym, ok := c.canvas.Symbol(i, j)
			if ok {
				term.SetCell(i, j, sym.Ch, term.Attribute(sym.Fg), term.Attribute(sym.Bg))
			}
		}
	}

	term.Flush()
}

// Repaints all View on the screen. Now the method is not efficient: at first clears a console and then draws all Views starting from the bottom
func (c *Composer) refreshScreen(invalidate bool) {
	c.canvas.Clear(ColorBlack)

	for _, wnd := range c.views {
		if invalidate {
			wnd.Repaint()
		}
		wnd.Draw(c.canvas)
	}

	c.redrawAll()
}

// CreateView constucts a new View
// posX and posY are top left coordinates of the View
// width and height are View size
// title is a View title shown inside the top View line
func (c *Composer) CreateView(posX, posY, width, height int, title string) View {
	view := NewWindow(c, posX, posY, width, height, title)

	c.views = append(c.views, view)
	view.Repaint()

	c.activateView(view)

	c.refreshScreen(false)

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

func (c *Composer) activateView(view View) bool {
	if c.topView() == view {
		for _, v := range c.views {
			v.SetActive(false)
		}
		view.SetActive(true)
		return true
	}

	var wList []View
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
		return false
	}

	view.SetActive(true)
	c.views = append(wList, view)
	return true
}

func (c *Composer) moveActiveWindowToBottom() bool {
	if len(c.views) < 2 {
		return false
	}

	if c.topView().Modal() {
		return false
	}

	event := Event{Type: EventActivate, X: 0} // send deactivated
	c.sendEventToActiveView(event)

	last := c.topView()

	for i := len(c.views) - 1; i > 0; i-- {
		c.views[i] = c.views[i-1]
	}

	c.views[0] = last
	if !c.activateView(c.topView()) {
		return false
	}

	event = Event{Type: EventActivate, X: 1} // send 'activated'
	c.sendEventToActiveView(event)
	c.refreshScreen(true)

	return true
}

func (c *Composer) termboxEventToLocal(ev term.Event) Event {
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

func (c *Composer) resizeTopView(ev term.Event) bool {
	view := c.topView()
	if view == nil {
		return false
	}

	w, h := view.Size()
	w1, h1 := w, h
	minW, minH := view.Constraints()
	if ev.Key == term.KeyArrowUp && minH < h {
		h--
	} else if ev.Key == term.KeyArrowLeft && minW < w {
		w--
	} else if ev.Key == term.KeyArrowDown {
		h++
	} else if ev.Key == term.KeyArrowRight {
		w++
	}

	if w1 != w || h1 != h {
		view.SetSize(w, h)
		// event := Event{Type: EventResize, X: w, Y: h}
		// c.sendEventToActiveView(event)
		c.refreshScreen(true)
	}

	return true
}

func (c *Composer) moveTopView(ev term.Event) bool {
	if len(c.views) > 0 {
		view := c.topView()
		if view != nil {
			x, y := view.Pos()
			w, h := view.Size()
			x1, y1 := x, y
			cx, cy := term.Size()
			if ev.Key == term.KeyArrowUp && y > 0 {
				y--
			} else if ev.Key == term.KeyArrowDown && y+h < cy {
				y++
			} else if ev.Key == term.KeyArrowLeft && x > 0 {
				x--
			} else if ev.Key == term.KeyArrowRight && x+w < cx {
				x++
			}

			if x1 != x || y1 != y {
				view.SetPos(x, y)
				// event := Event{Type: EventMove, X: x, Y: y}
				// c.sendEventToActiveView(event)
				c.refreshScreen(true)
			}
		}
		return true
	}

	return false
}

func (c *Composer) isDeadKey(ev term.Event) bool {
	if ev.Key == term.KeyCtrlS || ev.Key == term.KeyCtrlP || ev.Key == term.KeyCtrlW || ev.Key == term.KeyCtrlQ {
		c.ctrlKey = ev.Key
		c.lastKey = term.KeyEsc
		return true
	}

	c.ctrlKey = term.KeyEsc
	return false
}

func (c *Composer) processKeySeq(ev term.Event) bool {
	if c.ctrlKey == term.KeyEsc {
		return false
	}

	if c.ctrlKey == term.KeyCtrlQ {
		if c.ctrlKey == ev.Key {
			go c.Stop()
		} else if c.lastKey == term.KeyEsc {
			c.ctrlKey = ev.Key
		} else {
			c.ctrlKey = term.KeyEsc
			return false
		}

		return true
	}

	if c.ctrlKey == term.KeyCtrlS {
		if c.lastKey == term.KeyEsc {
			c.lastKey = ev.Key
		} else if c.lastKey != ev.Key {
			c.ctrlKey = term.KeyEsc
			return false
		}

		switch ev.Key {
		case term.KeyArrowUp, term.KeyArrowDown, term.KeyArrowLeft, term.KeyArrowRight:
			evCopy := ev
			c.resizeTopView(evCopy)
			return true
		}

		return false
	}

	if c.ctrlKey == term.KeyCtrlP {
		if c.lastKey == term.KeyEsc {
			c.lastKey = ev.Key
		} else if c.lastKey != ev.Key {
			c.ctrlKey = term.KeyEsc
			return false
		}

		switch ev.Key {
		case term.KeyArrowUp, term.KeyArrowDown, term.KeyArrowLeft, term.KeyArrowRight:
			evCopy := ev
			c.moveTopView(evCopy)
			return true
		default:
			return false
		}
	}

	if c.ctrlKey == term.KeyCtrlW {
		switch ev.Key {
		case term.KeyCtrlH:
			return c.moveActiveWindowToBottom()
		case term.KeyCtrlM:
			v := c.topView()
			maximized := v.Maximized()
			v.SetMaximized(!maximized)
			c.refreshScreen(true)
			return true
		case term.KeyCtrlC:
			c.closeTopView()
			return true
		default:
			return false
		}
	}

	c.ctrlKey = term.KeyEsc
	return true
}

// processKey returns false in case of the application should be terminated
func (c *Composer) processKey(ev term.Event) bool {
	if c.processKeySeq(ev) {
		return false
	}
	if c.isDeadKey(ev) {
		return false
	}

	if c.sendEventToActiveView(c.termboxEventToLocal(ev)) {
		c.topView().Repaint()
		c.refreshScreen(false)
	}

	return false
}

func (c *Composer) closeTopView() {
	if len(c.views) > 1 {
		view := c.topView()
		event := Event{Type: EventClose, X: 1}
		c.sendEventToActiveView(event)

		c.DestroyView(view)
		activate := c.topView()
		c.activateView(activate)
		event = Event{Type: EventActivate, X: 1} // send 'activated'
		c.sendEventToActiveView(event)

		c.refreshScreen(true)
	} else {
		go c.Stop()
	}
}

func (c *Composer) processMouseClick(ev term.Event) {
	view, hit := c.checkWindowUnderMouse(ev.MouseX, ev.MouseY)

	if view == nil {
		return
	}

	if c.topView() != view {
		if c.topView().Modal() {
			return
		}
		event := Event{Type: EventActivate, X: 0} // send 'deactivated'
		c.sendEventToActiveView(event)
		c.activateView(view)
		event = Event{Type: EventActivate, X: 1} // send 'activated'
		c.sendEventToActiveView(event)
		c.refreshScreen(true)
	} else if hit == HitInside {
		c.sendEventToActiveView(c.termboxEventToLocal(ev))
		c.refreshScreen(true)
	} else if hit == HitButtonClose {
		c.closeTopView()
	} else if hit == HitButtonBottom {
		c.moveActiveWindowToBottom()
	} else if hit == HitButtonMaximize {
		v := c.topView()
		maximized := v.Maximized()
		v.SetMaximized(!maximized)
		c.refreshScreen(true)
	}
}

// Stop sends termination event to Composer. Composer should stop
// console management and quit application
func (c *Composer) Stop() {
	ev := Event{Type: EventQuit}
	go c.PutEvent(ev)
}

// MainLoop starts the main application event loop
func (c *Composer) MainLoop() {
	c.refreshScreen(true)

	eventQueue := make(chan term.Event)
	go func() {
		for {
			eventQueue <- term.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventQueue:
			switch ev.Type {
			case term.EventKey:
				if c.processKey(ev) {
					return
				}
			case term.EventMouse:
				c.processMouseClick(ev)
			case term.EventError:
				panic(ev.Err)
			case term.EventResize:
				term.Flush()
				c.width, c.height = term.Size()
				c.initBuffer()
				for _, view := range c.views {
					if view.Maximized() {
						view.SetSize(c.width, c.height)
					}
				}
				c.refreshScreen(true)
			}
		case cmd := <-c.channel:
			if cmd.Type == EventRedraw {
				c.refreshScreen(true)
			} else if cmd.Type == EventQuit {
				return
			}
		}
	}
}

// PutEvent send event to a Composer directly.
// Used by Views to ask for repainting or for quitting the application
func (c *Composer) PutEvent(ev Event) {
	c.channel <- ev
}

// Close closes console management and makes a console cursor visible
func (c *Composer) Close() {
	term.SetCursor(3, 3)
	term.Close()
}

// SetCursorPos shows consolse cursor at given position.
// Setting cursor to -1,-1 hides cursor
func (c *Composer) SetCursorPos(x, y int) {
	term.SetCursor(x, y)
}

// DestroyView removes the View from the list of managed views
func (c *Composer) DestroyView(view View) {
	ev := Event{Type: EventClose}
	c.sendEventToActiveView(ev)

	var newOrder []View
	for i := 0; i < len(c.views); i++ {
		if c.views[i] != view {
			newOrder = append(newOrder, c.views[i])
		}
	}
	c.views = newOrder
	c.activateView(c.topView())
}

// Theme returns the theme manager. Theme manager implements
// Theme interface, so a caller can read the current colors
func (c *Composer) Theme() Theme {
	return c.themeManager
}

// Size returns size of the console(visible) buffer
func (c *Composer) Size() (int, int) {
	return term.Size()
}

func (c *Composer) Logger() *log.Logger {
	return c.logger
}
