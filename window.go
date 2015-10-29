package clui

import (
	"fmt"
	term "github.com/nsf/termbox-go"
	"log"
)

// Window is an implemetation of View managed by Composer.
type Window struct {
	ControlBase
	buttons  ViewButton
	canvas   Canvas
	parent   Screen
	pack     PackType
	children []Control
	controls []Control
	// dialog support
	modal   bool
	onClose func(Event)
}

/*
NewWindow creates a new View.
parent - is composer that manages all views.
x and y - initial View postion.
w and h - are minimal size of the view.
    The minimal view size cannot be less than 10x5
title - view title.
*/
func NewWindow(parent Screen, x, y, w, h int, title string) *Window {
	d := new(Window)
	d.canvas = NewFrameBuffer(w, h)

	if w == AutoSize {
		w = 10
	}
	if h == AutoSize {
		h = 5
	}

	d.SetSize(w, h)
	d.SetConstraints(w, h)
	d.SetTitle(title)
	d.SetPos(x, y)
	d.SetButtons(ButtonClose | ButtonBottom | ButtonMaximize)

	d.controls = make([]Control, 0)
	d.children = make([]Control, 0)
	d.parent = parent
	d.padSide, d.padTop, d.padX, d.padY = 1, 1, 1, 0

	return d
}

// SetSize changes control size. Constant DoNotChange can be
// used as placeholder to indicate that the control attrubute
// should be unchanged.
// Method panics if new size is less than minimal size.
// View automatically recalculates position and size of its children after changing its size
func (w *Window) SetSize(width, height int) {
	if width == w.width && height == w.height {
		return
	}

	if width != DoNotChange && (width > 1000 || width < w.minW) {
		panic(fmt.Sprintf("Invalid width: %v", width))
	}
	if height != DoNotChange && (height > 200 || height < w.minH) {
		panic(fmt.Sprintf("Invalid height: %v", height))
	}

	if width != DoNotChange {
		w.width = width
	}
	if height != DoNotChange {
		w.height = height
	}

	w.canvas.SetSize(w.width, w.height)
	RepositionControls(0, 0, w)
}

func (w *Window) applyConstraints() {
	width, height := w.Size()
	wM, hM := w.Constraints()

	newW, newH := width, height
	if width < wM {
		newW = wM
	}
	if height < hM {
		newH = hM
	}

	if newW != width || newH != height {
		w.SetSize(newW, newH)
	}
}

// SetConstraints sets new minimal size of control.
// If minimal size of the control is greater than the current
// control size then the control size is changed to fit minimal values
// The minimal constraints for view is width=10, height=5
func (w *Window) SetConstraints(width, height int) {
	if width >= 10 {
		w.minW = width
	}
	if height >= 5 {
		w.minH = height
	}

	w.applyConstraints()
}

// Draw paints the view screen buffer to a canvas. It does not
// repaint all view children
func (w *Window) Draw(canvas Canvas) {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			s, ok := w.canvas.Symbol(x, y)
			if ok {
				canvas.PutSymbol(x+w.x, y+w.y, s)
			} else {
				wx, wy := w.Size()
				panic(fmt.Sprintf("Invalid x, y: %vx%v of %vx%v", x, y, wx, wy))
			}
		}
	}
}

func (w *Window) buttonCount() int {
	count := 0
	if w.buttons&ButtonClose != 0 {
		count++
	}
	if w.buttons&ButtonBottom != 0 {
		count++
	}
	if w.buttons&ButtonMaximize != 0 {
		count++
	}

	return count
}

// Repaint draws the control and its children on the internal canvas
func (w *Window) Repaint() {
	tm := w.parent.Theme()
	bg := RealColor(tm, w.bg, ColorViewBack)

	w.canvas.Clear(bg)
	// paint all controls

	for _, child := range w.children {
		child.Repaint()
	}
	// paint itself - to overpaint any control that draws itself on the window border
	w.drawFrame(tm)
	w.drawTitle(tm)
	w.drawButtons(tm)
}

func (w *Window) drawTitle(tm Theme) {
	if w.title == "" {
		return
	}

	btnWidth := w.buttonCount()
	if btnWidth != 0 {
		btnWidth += 2
	}
	maxWidth := w.width - 2 - btnWidth
	text := Ellipsize(w.title, maxWidth)
	bg := RealColor(tm, w.bg, ColorViewBack)
	fg := RealColor(tm, w.fg, ColorViewText)
	w.canvas.PutText(1, 0, text, fg, bg)
}

func (w *Window) drawButtons(tm Theme) {
	if w.buttonCount() == 0 {
		return
	}

	bg, fg := RealColor(tm, w.bg, ColorViewBack), RealColor(tm, w.fg, ColorViewText)
	chars := []rune(tm.SysObject(ObjViewButtons))
	cMax, cBottom, cClose, cOpenB, cCloseB := chars[0], chars[1], chars[2], chars[3], chars[4]

	x := w.width - 2
	w.canvas.PutSymbol(x, 0, term.Cell{Ch: cCloseB, Fg: fg, Bg: bg})
	x--
	if w.buttons&ButtonClose != 0 {
		w.canvas.PutSymbol(x, 0, term.Cell{Ch: cClose, Fg: fg, Bg: bg})
		x--
	}
	if w.buttons&ButtonBottom != 0 {
		w.canvas.PutSymbol(x, 0, term.Cell{Ch: cBottom, Fg: fg, Bg: bg})
		x--
	}
	if w.buttons&ButtonMaximize != 0 {
		w.canvas.PutSymbol(x, 0, term.Cell{Ch: cMax, Fg: fg, Bg: bg})
		x--
	}
	w.canvas.PutSymbol(x, 0, term.Cell{Ch: cOpenB, Fg: fg, Bg: bg})
}

func (w *Window) drawFrame(tm Theme) {
	var chars string
	if w.active {
		chars = tm.SysObject(ObjDoubleBorder)
	} else {
		chars = tm.SysObject(ObjSingleBorder)
	}

	bg := RealColor(tm, w.bg, ColorViewBack)
	fg := RealColor(tm, w.fg, ColorViewText)

	w.canvas.DrawFrame(0, 0, w.width, w.height, fg, bg, chars)
}

// Canvas returns an internal graphic buffer to draw everything.
// Used by children controls - they paint themselves on the canvas
func (w *Window) Canvas() Canvas {
	return w.canvas
}

// SetButtons detemines which button is visible inside view
// title
func (w *Window) SetButtons(bi ViewButton) {
	w.buttons = bi
}

// Buttons returns the bit set of buttons displayed in Windows's title
// A set may contain any combination of: ButtonClose, ButtonBottom, and ButtonMaximize
func (w *Window) Buttons() ViewButton {
	return w.buttons
}

// SetPack changes the direction of children packing. Call the method only before any child is added to view. Otherwise, the method
// panics if a view already contains children
func (w *Window) SetPack(pk PackType) {
	if len(w.children) > 0 {
		panic("Control already has children")
	}

	w.pack = pk
}

// Pack returns direction in which a container packs
// its children: horizontal or vertical
func (w *Window) Pack() PackType {
	return w.pack
}

// RecalculateConstraints used by containers to recalculate new minimal size
// depending on its children constraints after a new child is added
func (w *Window) RecalculateConstraints() {
	width, height := w.Constraints()
	minW, minH := CalculateMinimalSize(w)

	newW, newH := width, height
	if minW > newW {
		newW = minW
	}
	if minH > newH {
		newH = minH
	}

	if newW != width || newH != height {
		w.SetConstraints(newW, newH)
	}
}

// RegisterControl adds a control to the view control list. It
// a list of all controls visible on the view - used to
// calculate the control under mouse when a user clicks, and
// to calculate the next control after a user presses TAB key
func (w *Window) RegisterControl(c Control) {
	w.controls = append(w.controls, c)
	w.RecalculateConstraints()
	RepositionControls(0, 0, w)
}

// AddChild add control to a list of view children. Minimal size
// of the view calculated as a sum of sizes of its children.
// Method panics if the same control is added twice
func (w *Window) AddChild(c Control, scale int) {
	if w.ChildExists(c) {
		panic("Cannot add the same control twice")
	}

	c.SetScale(scale)
	w.children = append(w.children, c)
	w.RegisterControl(c)
}

// Children returns the list of view children
func (w *Window) Children() []Control {
	return w.children
}

// Scale is a stub that always return DoNotScale becaue the
// scaling feature is not applied to views
func (w *Window) Scale() int {
	return DoNotScale
}

func (w *Window) controlAtPos(x, y int) Control {
	x -= w.x
	y -= w.y

	for id := len(w.controls) - 1; id >= 0; id-- {
		ctrl := w.controls[id]
		cw, ch := ctrl.Size()
		cx, cy := ctrl.Pos()

		if x >= cx && x < cx+cw && y >= cy && y < cy+ch {
			return ctrl
		}
	}

	return nil
}

// ProcessEvent processes all events come from the composer.
// If a view processes an event it should return true. If
// the method returns false it means that the view do
// not want or cannot process the event and the caller sends
// the event to the next target
func (w *Window) ProcessEvent(ev Event) bool {
	switch ev.Type {
	case EventKey, EventMouse:
		if ev.Type == EventKey && (ev.Key == term.KeyTab || (ev.Mod&term.ModAlt != 0 && (ev.Key == term.KeyPgup || ev.Key == term.KeyPgdn))) {
			forward := ev.Key != term.KeyPgup
			ctrl := w.ActiveControl()
			if ctrl != nil {
				ctrl.ProcessEvent(Event{Type: EventActivate, X: 0})
			}
			ctrl = w.nextControl(ctrl, forward)
			if ctrl != nil {
				// w.Logger().Printf("Activate control: %v", ctrl)
				w.ActivateControl(ctrl)
			}
			return true
		}
		if ev.Type == EventMouse {
			cunder := w.controlAtPos(ev.X, ev.Y)
			if cunder == nil {
				return true
			}

			w.ActivateControl(cunder)
		}
		ctrl := w.ActiveControl()
		if ctrl != nil {
			cx, cy := ctrl.Pos()
			cw, ch := ctrl.Size()
			ctrlX, ctrlY := ev.X-w.x, ev.Y-w.y
			if ev.Type == EventMouse && (ctrlX < cx || ctrlY < cy || ctrlX >= cx+cw || ctrlY >= cy+ch) {
				return false
			}
			copyEv := ev
			copyEv.X, copyEv.Y = ctrlX, ctrlY
			ctrl.ProcessEvent(copyEv)
			return true
		}
	case EventActivate:
		if ev.X == 0 {
			w.canvas.SetCursorPos(-1, -1)
		}
	case EventClose:
		if w.onClose != nil {
			w.onClose(Event{Type: EventClose, X: ev.X})
		}
		// case EventResize:
		// 	d.hideAllExtraControls()
		// 	d.recalculateControls()
	}

	return false
}

// ChildExists returns true if the container already has
// the control in its children list
func (w *Window) ChildExists(c Control) bool {
	for _, ctrl := range w.controls {
		if ctrl == c {
			return true
		}
	}

	return false
}

func (w *Window) nextControl(c Control, forward bool) Control {
	length := len(w.controls)

	if length == 0 {
		return nil
	}

	if length == 1 {
		return w.controls[0]
	}

	id := 0
	if c != nil {
		id = -1
		for idx, ct := range w.controls {
			if ct == c {
				id = idx
				break
			}
		}

		if id == -1 {
			return nil
		}
	}

	orig := id
	for {
		if forward {
			id++
		} else {
			id--
		}

		if id >= length {
			id = 0
		} else if id < 0 {
			id = length - 1
		}

		if w.controls[id].TabStop() {
			return w.controls[id]
		}

		if orig == id {
			if !w.controls[id].TabStop() {
				return nil
			}
			return c
		}
	}
}

// ActiveControl returns control that currently has focus or nil
// if there is no active control
func (w *Window) ActiveControl() Control {
	for _, ctrl := range w.controls {
		if ctrl.Active() {
			return ctrl
		}
	}

	return nil
}

// ActivateControl make the control active and previously
// focused control loses the focus. As a side effect the method
// emits two events: deactivate for previously focused and
// activate for new one if it is possible (EventActivate with
// different X values)
func (w *Window) ActivateControl(ctrl Control) {
	active := w.ActiveControl()
	if active == ctrl {
		return
	}
	if active != nil {
		active.ProcessEvent(Event{Type: EventActivate, X: 0})
		active.SetActive(false)
	}
	ctrl.SetActive(true)
	ctrl.ProcessEvent(Event{Type: EventActivate, X: 1})
}

func (w *Window) Logger() *log.Logger {
	return w.parent.Logger()
}

// Screen returns the composer that manages the view
func (w *Window) Screen() Screen {
	return w.parent
}

// Parent is a stub that always returns nil because the view
// cannot be added to any container
func (w *Window) Parent() Control {
	return nil
}

// TabStop is a stub that always returns false because the view
// cannot be selected by pressing TAB key
func (w *Window) TabStop() bool {
	return false
}

// HitTest returns the area that corresponds to the clicked
// position X, Y (absolute position in console window): title,
// internal view area, title button, border or outside the view
func (w *Window) HitTest(x, y int) HitResult {
	if x < w.x || y < w.y || x >= w.x+w.width || y >= w.y+w.height {
		return HitOutside
	}

	if x == w.x || x == w.x+w.width-1 || y == w.y+w.height-1 {
		return HitBorder
	}

	if y == w.y {
		dx := -3
		if w.buttons&ButtonClose != 0 {
			if x == w.x+w.width+dx {
				return HitButtonClose
			}
			dx--
		}
		if w.buttons&ButtonBottom != 0 {
			if x == w.x+w.width+dx {
				return HitButtonBottom
			}
			dx--
		}
		if w.buttons&ButtonMaximize != 0 {
			if x == w.x+w.width+dx {
				return HitButtonMaximize
			}
		}
	}

	return HitInside
}

// SetModal enables or disables modal mode
func (w *Window) SetModal(modal bool) {
	w.modal = modal
}

// Modal returns if the view is in modal mode.In modal mode a
// user cannot switch to any other view until the user closes
// the modal view. Used by confirmation and select dialog to be
// sure that the user has made a choice before continuing work
func (w *Window) Modal() bool {
	return w.modal
}

// OnClose sets a callback that is called when view is closed.
// For dialogs after windows is closed a user can check the
// close result
func (w *Window) OnClose(fn func(Event)) {
	w.onClose = fn
}
