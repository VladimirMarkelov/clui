package clui

import (
	"fmt"
	term "github.com/nsf/termbox-go"
	"log"
)

type Window struct {
	ControlBase
	buttons         ViewButton
	canvas          Canvas
	parent          Screen
	pack            PackType
	padX, padY      int
	padTop, padSide int
	children        []Control
	controls        []Control
}

func NewWindow(parent Screen, x, y, w, h int, title string) *Window {
	d := new(Window)
	d.canvas = NewFrameBuffer(w, h)

	d.SetConstraints(10, 5)
	d.SetTitle(title)
	d.SetSize(w, h)
	d.SetPos(x, y)
	d.SetButtons(ButtonClose | ButtonBottom | ButtonMaximize)

	d.controls = make([]Control, 0)
	d.children = make([]Control, 0)
	d.parent = parent
	d.padSide, d.padTop, d.padX, d.padY = 1, 1, 1, 0

	d.fg = ColorDefault
	d.bg = ColorDefault

	return d
}

func (w *Window) SetSize(width, height int) {
	if width > 1000 || width < w.minW {
		panic(fmt.Sprintf("Invalid width: %v", width))
	}
	if height > 200 || height < w.minH {
		panic(fmt.Sprintf("Invalid height: %v", height))
	}

	w.width = width
	w.height = height

	w.canvas.SetSize(width, height)
	RepositionControls(0, 0, w)
}

func (w *Window) ApplyConstraints() {
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

func (w *Window) SetConstraints(width, height int) {
	if width >= 10 {
		w.minW = width
	}
	if height >= 5 {
		w.minH = height
	}

	w.ApplyConstraints()
}

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

func (w *Window) ButtonCount() int {
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

func (w *Window) Repaint() {
	tm := w.parent.Theme()
	bg := RealColor(tm, w.bg, ColorViewBack)

	w.canvas.Clear(bg)
	// paint all controls

	for _, child := range w.children {
		child.Repaint()
	}
	// paint itself - to overpaint any control that draws itself on the window border
	w.DrawFrame(tm)
	w.DrawTitle(tm)
	w.DrawButtons(tm)
}

func (w *Window) DrawTitle(tm Theme) {
	if w.title == "" {
		return
	}

	btnWidth := w.ButtonCount()
	if btnWidth != 0 {
		btnWidth += 2
	}
	maxWidth := w.width - 2 - btnWidth
	text := Ellipsize(w.title, maxWidth)
	bg := RealColor(tm, w.bg, ColorViewBack)
	fg := RealColor(tm, w.fg, ColorViewText)
	w.canvas.PutText(1, 0, text, fg, bg)
}

func (w *Window) DrawButtons(tm Theme) {
	if w.ButtonCount() == 0 {
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

func (w *Window) DrawFrame(tm Theme) {
	var chars string = ""
	if w.active {
		chars = tm.SysObject(ObjDoubleBorder)
	} else {
		chars = tm.SysObject(ObjSingleBorder)
	}

	bg := RealColor(tm, w.bg, ColorViewBack)
	fg := RealColor(tm, w.fg, ColorViewText)

	w.canvas.DrawFrame(0, 0, w.width, w.height, fg, bg, chars)
}

func (w *Window) Paddings() (int, int, int, int) {
	return w.padSide, w.padTop, w.padX, w.padY
}

func (w *Window) Canvas() Canvas {
	return w.canvas
}

func (w *Window) SetButtons(bi ViewButton) {
	w.buttons = bi
}

func (w *Window) Buttons() ViewButton {
	return w.buttons
}

func (w *Window) SetPack(pk PackType) {
	if len(w.children) > 0 {
		panic("Control already has children")
	}

	w.pack = pk
}

func (w *Window) Pack() PackType {
	return w.pack
}

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

func (w *Window) RegisterControl(c Control) {
	w.controls = append(w.controls, c)
	w.RecalculateConstraints()
	RepositionControls(0, 0, w)
}

func (w *Window) AddChild(c Control, scale int) {
	if w.ChildExists(c) {
		panic("Cannot add the same control twice")
	}

	c.SetScale(scale)
	w.children = append(w.children, c)
	w.RegisterControl(c)
}

func (w *Window) Children() []Control {
	return w.children
}

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

func (w *Window) ProcessEvent(ev Event) bool {
	switch ev.Type {
	case EventKey, EventMouse:
		if ev.Type == EventKey && (ev.Key == term.KeyTab || (ev.Mod&term.ModAlt != 0 && (ev.Key == term.KeyPgup || ev.Key == term.KeyPgdn))) {
			forward := ev.Key != term.KeyPgup
			ctrl := w.ActiveControl()
			if ctrl != nil {
				ctrl.ProcessEvent(Event{Type: EventActivate, X: 0})
			}
			ctrl = w.NextControl(ctrl, forward)
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

			ctrl := w.ActiveControl()
			ctrl.ProcessEvent(Event{Type: EventActivate, X: 0})
			w.ActivateControl(cunder)
		}
		ctrl := w.ActiveControl()
		if ctrl != nil {
			// w.Logger().Printf("Active control %v -- %v", ctrl.Title(), ctrl)
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
		// case EventResize:
		// 	d.hideAllExtraControls()
		// 	d.recalculateControls()
	}

	return false
}

func (w *Window) ChildExists(c Control) bool {
	for _, ctrl := range w.controls {
		if ctrl == c {
			return true
		}
	}

	return false
}

func (w *Window) NextControl(c Control, forward bool) Control {
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
			} else {
				return c
			}
		}
	}
}

func (w *Window) ActiveControl() Control {
	for _, ctrl := range w.controls {
		if ctrl.Active() {
			return ctrl
		}
	}

	return nil
}

func (w *Window) ActivateControl(ctrl Control) {
	active := w.ActiveControl()
	if active == ctrl {
		return
	}
	if active != nil {
		active.SetActive(false)
	}
	ctrl.SetActive(true)
	ctrl.ProcessEvent(Event{Type: EventActivate, X: 1})
}

func (w *Window) Logger() *log.Logger {
	return w.parent.Logger()
}

func (w *Window) Screen() Screen {
	return w.parent
}

func (w *Window) Parent() Control {
	return nil
}

func (w *Window) TabStop() bool {
	return false
}

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
