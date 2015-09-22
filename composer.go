package clui

import (
	"github.com/VladimirMarkelov/termbox-go"
	"log"
	"os"
)

// An service object that manages Windows and console, processes events, and provides service methods
// One application must have only one object of this type
type Composer struct {
	// list of visible Windows
	windows map[WinId]Window
	// console width and height
	width, height int
	// console canvas
	screen *FrameBuffer
	// Window draw order. The last Window is active Window
	windowOrder []WinId
	// ID assigned to the last created Window
	lastWindowId WinId
	// a channel to communicate with Windows(e.g, Windows send redraw event to this channel)
	channel chan InternalEvent
	// What kind of dragging is active: nothing is dragged, something is being moved, something is being resized
	dragging DragAction
	// Window ID that is currently dragged. NullWindow if no Window is dragged
	dragObject WinId
	// Console position where mouse button was pressed and drag action started
	dragStartX, dragStartY int
	// which part of dragged Window is active, it determines whether Window is resizing or dragging and what size is changing
	dragHit HitResult

	// current color scheme
	themeManager *ThemeManager

	// multi key sequences support. The flag below are true if the last keyboard combination was Ctrl+S or Ctrl+W respectively
	ctrlSpressed bool
	ctrlWpressed bool
	// last pressed key - to make repeatable actions simpler, e.g, at first one presses Ctrl+S and then just repeatedly presses arrow lest to resize Window
	lastKey termbox.Key

	//debug
	logger *log.Logger
}

func (c *Composer) initBuffer(w, h int) *FrameBuffer {
	bf := NewFrameBuffer(w, h)
	bf.Clear(ColorBlack)
	return bf
}

// Initialize library and starts console management
func (c *Composer) Init() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	termbox.HideCursor()

	file, _ := os.OpenFile("debug.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	c.logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	c.logger.Printf("----------------------------------")

	c.lastWindowId = 0
	c.channel = make(chan InternalEvent)

	c.windows = make(map[WinId]Window)
	c.windowOrder = make([]WinId, 0)

	c.width, c.height = termbox.Size()
	c.screen = c.initBuffer(c.width, c.height)

	c.themeManager = NewThemeManager()

	c.ctrlSpressed = false
	c.ctrlWpressed = false

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	c.redrawAll()
}

func (c *Composer) redrawAll() {
	for j := 0; j < c.height; j++ {
		for i := 0; i < c.width; i++ {
			sym, ok := c.screen.GetSymbol(i, j)
			if ok {
				termbox.SetCell(i, j, sym.ch, termbox.Attribute(sym.fg), termbox.Attribute(sym.bg))
			}
		}
	}

	termbox.Flush()
}

// Repaints all Windows on the screen. Now the method is not efficient: at first clears a console and then draws all Windows starting from the bottom
func (c *Composer) RefreshScreen() {
	// TODO: make more intelligent (for each Cell ask starting from the top most if it is its Cell and write it)
	for yy := 0; yy < c.height; yy++ {
		for xx := 0; xx < c.width; xx++ {
			c.screen.rawRuneOut(xx, yy, ' ', ColorWhite, ColorBlack)
		}
	}

	for _, winId := range c.windowOrder {
		wnd, ok := c.windows[winId]
		if !ok {
			continue
		}

		wnd.Redraw()

		posx, posy := wnd.GetPos()
		w, h := wnd.GetSize()
		for y := posy; y < posy+h; y++ {
			for x := posx; x < posx+w; x++ {
				if x < c.width && y < c.height {
					sym, ok := wnd.GetScreenSymbol(x, y)
					if ok {
						c.screen.rawRuneOut(x, y, sym.ch, sym.fg, sym.bg)
					}
				}
			}
		}
	}
	c.redrawAll()
}

func (c *Composer) makeTopViewActive() {
	for i := 0; i < len(c.windowOrder)-1; i++ {
		window, ok := c.windows[c.windowOrder[i]]
		if ok {
			window.SetActive(false)
		}
	}

	window, ok := c.windows[c.windowOrder[len(c.windowOrder)-1]]
	if ok {
		window.SetActive(true)
	}
}

func (c *Composer) CreateWindow(posX, posY, width, height int, title string, props Props) Window {
	c.lastWindowId++
	id := c.lastWindowId
	view := NewView(c, id, posX, posY, width, height, title)
	view.SetBorderIcons(IconClose | IconBottom)

	c.windows[id] = view
	view.Redraw()
	c.windowOrder = append(c.windowOrder, id)

	c.makeTopViewActive()

	c.RefreshScreen()

	return view
}

func (c *Composer) DestroyControl(parent Window, control Control) {
	parent.RemoveControl(control)
}

func (c *Composer) checkWindowUnderMouse(screenX, screenY int) (WinId, HitResult) {
	if len(c.windowOrder) == 0 {
		return NullWindow, HitOutside
	}

	for i := len(c.windowOrder) - 1; i >= 0; i-- {
		window, ok := c.windows[c.windowOrder[i]]
		if ok {
			hit := window.HitTest(screenX, screenY)
			if hit != HitOutside {
				return c.windowOrder[i], hit
			}
		}
	}

	return NullWindow, HitOutside
}

func (c *Composer) getWindow(id WinId) Window {
	wnd, ok := c.windows[id]

	if !ok {
		return nil
	}

	return wnd
}

func (c *Composer) getTopWindow() WinId {
	if len(c.windowOrder) == 0 {
		return NullWindow
	}

	return c.windowOrder[len(c.windowOrder)-1]
}

func (c *Composer) makeWindowTop(wid WinId) {
	for i := 0; i < len(c.windowOrder); i++ {
		if c.windowOrder[i] == wid {
			c.windowOrder = append(c.windowOrder[:i], c.windowOrder[i+1:]...)
			break
		}
	}

	c.windowOrder = append(c.windowOrder, wid)
	c.makeTopViewActive()
}

func (c *Composer) moveActiveWindowToBottom() bool {
	if len(c.windowOrder) < 2 {
		return false
	}

	event := Event{Type: EventActivate, X: 0} // send deactivated
	c.sendEventToActiveView(event)

	last := c.windowOrder[len(c.windowOrder)-1]

	for i := len(c.windowOrder) - 1; i > 0; i-- {
		c.windowOrder[i] = c.windowOrder[i-1]
	}

	c.windowOrder[0] = last
	c.makeTopViewActive()

	event = Event{Type: EventActivate, X: 1} // send 'activated'
	c.sendEventToActiveView(event)
	c.RefreshScreen()

	return true
}

func (c *Composer) termboxEventToLocal(ev termbox.Event) Event {
	e := Event{Type: EventType(ev.Type), Ch: ev.Ch, Key: ev.Key, Err: ev.Err, X: ev.MouseX, Y: ev.MouseY, Mod: ev.Mod}
	return e
}

func (c *Composer) sendEventToActiveView(ev Event) bool {
	winid := c.getTopWindow()
	if winid != NullWindow {
		window := c.getWindow(winid)
		if window != nil {
			return window.ProcessEvent(ev)
		}
	}

	return false
}

func (c *Composer) resizeTopWindow(ev termbox.Event) bool {
	if len(c.windowOrder) > 0 {
		if ev.Mod&termbox.ModControl != 0 {
			window := c.getWindow(c.getTopWindow())
			if window != nil {
				w, h := window.GetSize()
				w1, h1 := w, h
				minW, minH := window.GetConstraints()
				if ev.Key == termbox.KeyArrowUp && (minH == -1 || minH < h) {
					h--
				} else if ev.Key == termbox.KeyArrowLeft && (minW < w || minW == -1) {
					w--
				}

				if w1 != w || h1 != h {
					window.SetSize(w, h)
					event := Event{Type: EventResize, X: w, Y: h}
					c.sendEventToActiveView(event)
					c.RefreshScreen()
				}
			}
		}
		return true
	} else {
		return false
	}
}

func (c *Composer) moveTopWindow(ev termbox.Event) bool {
	if len(c.windowOrder) > 0 {
		window := c.getWindow(c.getTopWindow())
		if window != nil {
			x, y := window.GetPos()
			w, h := window.GetSize()
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
				window.SetPos(x, y)
				event := Event{Type: EventMove, X: x, Y: y}
				c.sendEventToActiveView(event)
				c.RefreshScreen()
			}
		}
		return true
	} else {
		return false
	}
}

func (c *Composer) isDeadKey(ev termbox.Event) bool {
	if ev.Key == termbox.KeyCtrlS {
		c.ctrlSpressed = true
		c.lastKey = termbox.KeyEsc
		return true
	}
	if ev.Key == termbox.KeyCtrlW {
		c.ctrlWpressed = true
		c.lastKey = termbox.KeyEsc
		return true
	}

	c.ctrlSpressed = false
	c.ctrlWpressed = false

	return false
}

func (c *Composer) processKeySeq(ev termbox.Event) bool {
	if !c.ctrlSpressed && !c.ctrlWpressed {
		return false
	}

	if c.ctrlSpressed {
		c.ctrlWpressed = false

		if c.lastKey == termbox.KeyEsc {
			c.lastKey = ev.Key
		} else if c.lastKey != ev.Key {
			c.ctrlSpressed = false
			return false
		}

		switch ev.Key {
		case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight:
			evCopy := ev
			evCopy.Mod = termbox.ModControl
			c.resizeTopWindow(evCopy)
			return true
		}

		return false
	}
	if c.ctrlWpressed {
		c.ctrlSpressed = false

		if c.lastKey == termbox.KeyEsc {
			c.lastKey = ev.Key
		} else if c.lastKey != ev.Key {
			c.ctrlWpressed = false
			return false
		}

		switch ev.Key {
		case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight:
			evCopy := ev
			evCopy.Mod = termbox.ModAlt
			c.moveTopWindow(evCopy)
			return true
		case termbox.KeyCtrlH:
			if len(c.windowOrder) > 1 && ev.Mod&termbox.ModControl != 0 {
				c.moveActiveWindowToBottom()
			}
			return true
		}

		return false
	}

	c.ctrlSpressed = false
	c.ctrlWpressed = false
	return true
}

func (c *Composer) processKey(ev termbox.Event) bool {
	processed := false

	if c.processKeySeq(ev) {
		return false
	}
	if c.isDeadKey(ev) {
		return false
	}

	switch ev.Key {
	case termbox.KeyCtrlQ:
		return true
	case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight:
		if ev.Mod&termbox.ModControl != 0 {
			processed = processed || c.resizeTopWindow(ev)
		}

		if ev.Mod&termbox.ModAlt != 0 {
			processed = processed || c.moveTopWindow(ev)
		}

		if !processed {
			if c.sendEventToActiveView(c.termboxEventToLocal(ev)) {
				c.RefreshScreen()
			}
		}
	case termbox.KeyEnd:
		if len(c.windowOrder) > 1 && ev.Mod&termbox.ModControl != 0 {
			processed = c.moveActiveWindowToBottom()
		}
		if !processed {
			if c.sendEventToActiveView(c.termboxEventToLocal(ev)) {
				c.RefreshScreen()
			}
		}
	default:
		if c.sendEventToActiveView(c.termboxEventToLocal(ev)) {
			c.RefreshScreen()
		}
	}

	return false
}

func (c *Composer) processMouseClick(ev termbox.Event) {
	winId, hit := c.checkWindowUnderMouse(ev.MouseX, ev.MouseY)
	if c.getTopWindow() != winId {
		event := Event{Type: EventActivate, X: 0} // send 'deactivated'
		c.sendEventToActiveView(event)
		c.makeWindowTop(winId)
		event = Event{Type: EventActivate, X: 1} // send 'activated'
		c.sendEventToActiveView(event)
		c.RefreshScreen()
	} else if hit == HitInside {
		// c.logger.Printf("Clicked inside window %v", int(winId))
		c.sendEventToActiveView(c.termboxEventToLocal(ev))
		c.RefreshScreen()
	} else if hit == HitButtonClose {
		if len(c.windowOrder) > 1 {
			event := Event{Type: EventClose}
			c.sendEventToActiveView(event)

			c.DestroyWindow(winId)
			activate := c.windowOrder[len(c.windowOrder)-1]
			c.makeWindowTop(activate)
			event = Event{Type: EventActivate, X: 1} // send 'activated'
			c.sendEventToActiveView(event)

			c.RefreshScreen()
		}
	} else if hit == HitButtonBottom {
		c.moveActiveWindowToBottom()
	}
}

func (c *Composer) stopDragging() {
	c.dragging = DragNone
	c.dragStartX, c.dragStartY = -1, -1
	c.dragObject = NullWindow
}

func (c *Composer) processMouseRelease(ev termbox.Event) {
	if c.dragging == DragNone {
		c.sendEventToActiveView(c.termboxEventToLocal(ev))
	}
	c.stopDragging()
}

func (c *Composer) processMousePress(ev termbox.Event) {
	winId, hit := c.checkWindowUnderMouse(ev.MouseX, ev.MouseY)
	if c.getTopWindow() != winId {
		event := Event{Type: EventActivate, X: 0} // send 'deactivated'
		c.sendEventToActiveView(event)
		c.makeWindowTop(winId)
		event = Event{Type: EventActivate, X: 1} // send 'activated'
		c.sendEventToActiveView(event)
		c.RefreshScreen()
	} else {
		if hit == HitHeader {
			c.dragging = DragMove
			c.dragObject = winId
			c.dragStartX, c.dragStartY = ev.MouseX, ev.MouseY
		} else if hit == HitLeftBorder || hit == HitRightBorder || hit == HitBottomBorder || hit == HitBottomLeft || hit == HitBottomRight || hit == HitTopLeft || hit == HitTopRight {
			c.dragging = DragResize
			c.dragObject = winId
			c.dragStartX, c.dragStartY = ev.MouseX, ev.MouseY
			c.dragHit = hit
		} else if hit == HitInside {
			c.sendEventToActiveView(c.termboxEventToLocal(ev))
		}
	}
}

func (c *Composer) moveAndResize(dx, dy int, wnd Window) bool {
	if dx == 0 && dy == 0 {
		return false
	}

	posX, posY := wnd.GetPos()
	w, h := wnd.GetSize()
	mnX, mnY := wnd.GetConstraints()
	newW, newH := w, h
	newX, newY := posX, posY

	if c.dragHit == HitLeftBorder {
		newX = posX + dx
		newW = w - dx
	} else if c.dragHit == HitRightBorder {
		newW = w + dx
	} else if c.dragHit == HitBottomBorder {
		newH = h + dy
	} else if c.dragHit == HitTopLeft {
		newX, newY = posX+dx, posY+dy
		newW, newH = w-dx, h-dy
	} else if c.dragHit == HitBottomLeft {
		newX = posX + dx
		newW, newH = w-dx, h+dy
	} else if c.dragHit == HitTopRight {
		newY = posY + dy
		newW, newH = w+dx, h-dy
	} else if c.dragHit == HitBottomRight {
		newW, newH = w+dx, h+dy
	}

	if (mnX != -1 && newW < mnX) || (mnY != -1 && newH < mnY) {
		return false
	}

	wnd.SetPos(newX, newY)
	wnd.SetSize(newW, newH)

	if posX != newX || posY != newY {
		event := Event{Type: EventMove, X: newX, Y: newY}
		c.sendEventToActiveView(event)
	}
	if newW != w || newH != h {
		event := Event{Type: EventResize, X: w, Y: h}
		c.sendEventToActiveView(event)
	}

	wnd.Redraw()

	return true
}

func (c *Composer) processMouseMove(ev termbox.Event) {
	if c.dragging == DragNone {
		c.sendEventToActiveView(c.termboxEventToLocal(ev))
		c.RefreshScreen()
		return
	}

	wnd, ok := c.windows[c.dragObject]
	if !ok || c.dragObject == NullWindow {
		c.stopDragging()
		return
	}

	if c.dragging == DragMove {
		dx, dy := ev.MouseX-c.dragStartX, ev.MouseY-c.dragStartY
		if dx != 0 || dy != 0 {
			c.dragStartY = ev.MouseY
			c.dragStartX = ev.MouseX
			posX, posY := wnd.GetPos()
			wnd.SetPos(posX+dx, posY+dy)
			event := Event{Type: EventMove, X: posX + dx, Y: posY + dy}
			c.sendEventToActiveView(event)
			c.RefreshScreen()
		}
	} else if c.dragging == DragResize && c.dragHit != HitOutside {

		dx, dy := 0, 0
		if c.dragHit == HitLeftBorder || c.dragHit == HitRightBorder {
			dx = ev.MouseX - c.dragStartX
		} else if c.dragHit == HitBottomBorder {
			dy = ev.MouseY - c.dragStartY
		} else if c.dragHit == HitTopLeft || c.dragHit == HitTopRight || c.dragHit == HitBottomLeft || c.dragHit == HitBottomRight {
			dx = ev.MouseX - c.dragStartX
			dy = ev.MouseY - c.dragStartY
		}

		if c.moveAndResize(dx, dy, wnd) {
			c.dragStartX = ev.MouseX
			c.dragStartY = ev.MouseY

			c.RefreshScreen()
		}
	}
}

// Asks a Composer to stops console management and quit application
func (c *Composer) Stop() {
	ev := InternalEvent{act: EventQuit}
	go c.SendEvent(ev)
}

// Main event loop
func (c *Composer) MainLoop() {
	c.redrawAll()

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
			case termbox.EventMouseScroll:
				if c.sendEventToActiveView(c.termboxEventToLocal(ev)) {
					c.RefreshScreen()
				}
			case termbox.EventKey:
				if c.processKey(ev) {
					return
				}
			case termbox.EventMouse, termbox.EventMouseClick:
				c.processMouseClick(ev)
			case termbox.EventMouseRelease:
				c.processMouseRelease(ev)
			case termbox.EventMousePress:
				c.processMousePress(ev)
			case termbox.EventMouseMove:
				c.processMouseMove(ev)
			case termbox.EventError:
				panic(ev.Err)
			case termbox.EventResize:
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				c.width, c.height = termbox.Size()
				c.screen = c.initBuffer(c.width, c.height)
				c.RefreshScreen()
			}
		case cmd := <-c.channel:
			if cmd.act == EventRedraw {
				c.RefreshScreen()
			} else if cmd.act == EventQuit {
				return
			}
		}
	}
}

// Send event to a Composer. Used by Windows to ask for repainting or for quitting the application
func (c *Composer) SendEvent(ev InternalEvent) {
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

// Returns thememanager to get current theme colors etc
func (c *Composer) GetThemeManager() *ThemeManager {
	return c.themeManager
}

func (c *Composer) DestroyWindow(winId WinId) {
	ev := Event{Type: EventClose}
	c.sendEventToActiveView(ev)

	newOrder := make([]WinId, 0)
	for i := 0; i < len(c.windowOrder); i++ {
		if c.windowOrder[i] != winId {
			newOrder = append(newOrder, c.windowOrder[i])
		}
	}
	c.windowOrder = newOrder
	delete(c.windows, winId)
}
