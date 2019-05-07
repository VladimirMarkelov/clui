package clui

import (
	мКнст "./пакКонстанты"
	term "github.com/nsf/termbox-go"
	"sync"
	"sync/atomic"
	мИнт "./пакИнтерфейсы"
	мСоб "./пакСобытия"
)

// BaseControl is a base for all visible controls.
// Every new control must inherit it or implement
// the same set of methods
type BaseControl struct {
	refID         int64
	title         string
	x, y          int
	width, height int
	minW, minH    int
	scale         int
	fg, bg        term.Attribute
	fgActive      term.Attribute
	bgActive      term.Attribute
	tabSkip       bool
	disabled      bool
	hidden        bool
	align         мИнт.Align
	parent        мИнт.ИВиджет
	inactive      bool
	modal         bool
	padX, padY    int
	gapX, gapY    int
	pack          мИнт.PackType
	children      []мИнт.ИВиджет
	mtx           sync.RWMutex
	onActive      func(active bool)
	style         string
	clipped       bool
	clipper       *rect
}

var (
	globalRefId int64
)

func nextRefId() int64 {
	return atomic.AddInt64(&globalRefId, 1)
}

//NewBaseControl --
func NewBaseControl() *BaseControl {
	return &BaseControl{refID: nextRefId()}
}

//SetClipped --
func (c *BaseControl) SetClipped(clipped bool) {
	c.clipped = clipped
}

//Clipped --
func (c *BaseControl) Clipped() bool {
	return c.clipped
}

//SetStyle --
func (c *BaseControl) SetStyle(style string) {
	c.style = style
}

//Style --
func (c *BaseControl) Style() string {
	return c.style
}

//RefID --
func (c *BaseControl) RefID() int64 {
	return c.refID
}

//Title --
func (c *BaseControl) Title() string {
	return c.title
}

//SetTitle --
func (c *BaseControl) SetTitle(title string) {
	c.title = title
}

//Size --
func (c *BaseControl) Size() (widht int, height int) {
	return c.width, c.height
}

//SetSize --
func (c *BaseControl) SetSize(width, height int) {
	if width < c.minW {
		width = c.minW
	}
	if height < c.minH {
		height = c.minH
	}

	if height != c.height || width != c.width {
		c.height = height
		c.width = width
	}
}

//Pos --
func (c *BaseControl) Pos() (x int, y int) {
	return c.x, c.y
}

//SetPos --
func (c *BaseControl) SetPos(x, y int) {
	if c.clipped && c.clipper != nil {
		cx, cy, _, _ := c.Clipper()
		px, py := c.Paddings()

		distX := cx - c.x
		distY := cy - c.y

		c.clipper.x = x + px
		c.clipper.y = y + py

		c.x = (x - distX) + px
		c.y = (y - distY) + py
	} else {
		c.x = x
		c.y = y
	}
}

//applyConstraints --
func (c *BaseControl) applyConstraints() {
	ww, hh := c.width, c.height
	if ww < c.minW {
		ww = c.minW
	}
	if hh < c.minH {
		hh = c.minH
	}
	if hh != c.height || ww != c.width {
		c.SetSize(ww, hh)
	}
}

//Constraints --
func (c *BaseControl) Constraints() (minw int, minh int) {
	return c.minW, c.minH
}

//SetConstraints --
func (c *BaseControl) SetConstraints(minw, minh int) {
	c.minW = minw
	c.minH = minh
	c.applyConstraints()
}

//Active --
func (c *BaseControl) Active() bool {
	return !c.inactive
}

//SetActive --
func (c *BaseControl) SetActive(active bool) {
	c.inactive = !active

	if c.onActive != nil {
		c.onActive(active)
	}
}

//OnActive --
func (c *BaseControl) OnActive(fn func(active bool)) {
	c.onActive = fn
}

//TabStop --
func (c *BaseControl) TabStop() bool {
	return !c.tabSkip
}

//SetTabStop --
func (c *BaseControl) SetTabStop(tabstop bool) {
	c.tabSkip = !tabstop
}

//Enabled --
func (c *BaseControl) Enabled() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return !c.disabled
}

//SetEnabled --
func (c *BaseControl) SetEnabled(enabled bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.disabled = !enabled
}

//Visible --
func (c *BaseControl) Visible() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return !c.hidden
}

//SetVisible --
func (c *BaseControl) SetVisible(visible bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if visible == !c.hidden {
		return
	}

	c.hidden = !visible
	if c.parent == nil {
		return
	}

	p := c.Parent()
	for p.Parent() != nil {
		p = p.Parent()
	}

	go func() {
		if FindFirstActiveControl(c) != nil && !c.inactive {
			PutEvent(мСоб.Event{Type: мКнст.EventKey, Key: term.KeyTab})
		}
		PutEvent(мСоб.Event{Type: мКнст.EventLayout, Target: p})
	}()
}

//Parent --
func (c *BaseControl) Parent() мИнт.ИВиджет {
	return c.parent
}

//SetParent --
func (c *BaseControl) SetParent(parent мИнт.ИВиджет) {
	if c.parent == nil {
		c.parent = parent
	}
}

//Modal --
func (c *BaseControl) Modal() bool {
	return c.modal
}

//SetModal --
func (c *BaseControl) SetModal(modal bool) {
	c.modal = modal
}

//Paddings --
func (c *BaseControl) Paddings() (px int, py int) {
	return c.padX, c.padY
}

//SetPaddings __
func (c *BaseControl) SetPaddings(px, py int) {
	if px >= 0 {
		c.padX = px
	}
	if py >= 0 {
		c.padY = py
	}
}

//Gaps --
func (c *BaseControl) Gaps() (dx int, dy int) {
	return c.gapX, c.gapY
}

//SetGaps --
func (c *BaseControl) SetGaps(dx, dy int) {
	if dx >= 0 {
		c.gapX = dx
	}
	if dy >= 0 {
		c.gapY = dy
	}
}

//Pack --
func (c *BaseControl) Pack() мИнт.PackType {
	return c.pack
}

//SetPack --
func (c *BaseControl) SetPack(pack мИнт.PackType) {
	c.pack = pack
}

//Scale --
func (c *BaseControl) Scale() int {
	return c.scale
}

//SetScale --
func (c *BaseControl) SetScale(scale int) {
	if scale >= 0 {
		c.scale = scale
	}
}

//Align --
func (c *BaseControl) Align() мИнт.Align {
	return c.align
}

//SetAlign --
func (c *BaseControl) SetAlign(align мИнт.Align) {
	c.align = align
}

//TextColor --
func (c *BaseControl) TextColor() term.Attribute {
	return c.fg
}

//SetTextColor --
func (c *BaseControl) SetTextColor(clr term.Attribute) {
	c.fg = clr
}

//BackColor --
func (c *BaseControl) BackColor() term.Attribute {
	return c.bg
}

//SetBackColor --
func (c *BaseControl) SetBackColor(clr term.Attribute) {
	c.bg = clr
}

//childCount --
func (c *BaseControl) childCount() int {
	cnt := 0
	for _, child := range c.children {
		if child.Visible() {
			cnt++
		}
	}

	return cnt
}

//ResizeChildren --
func (c *BaseControl) ResizeChildren() {
	children := c.childCount()
	if children == 0 {
		return
	}

	fullWidth := c.width - 2*c.padX
	fullHeight := c.height - 2*c.padY
	if c.pack == мКнст.Horizontal {
		fullWidth -= (children - 1) * c.gapX
	} else {
		fullHeight -= (children - 1) * c.gapY
	}

	totalSc := c.ChildrenScale()
	minWidth := 0
	minHeight := 0
	for _, child := range c.children {
		if !child.Visible() {
			continue
		}

		cw, ch := child.MinimalSize()
		if c.pack == мКнст.Horizontal {
			minWidth += cw
		} else {
			minHeight += ch
		}
	}

	aStep := 0
	diff := fullWidth - minWidth
	if c.pack == мКнст.Vertical {
		diff = fullHeight - minHeight
	}
	if totalSc > 0 {
		aStep = int(float32(diff) / float32(totalSc))
	}

	for _, ctrl := range c.children {
		if !ctrl.Visible() {
			continue
		}

		tw, th := ctrl.MinimalSize()
		sc := ctrl.Scale()
		d := int(ctrl.Scale() * aStep)
		if c.pack == мКнст.Horizontal {
			if sc != 0 {
				if sc == totalSc {
					tw += diff
					d = diff
				} else {
					tw += d
				}
			}
			th = fullHeight
		} else {
			if sc != 0 {
				if sc == totalSc {
					th += diff
					d = diff
				} else {
					th += d
				}
			}
			tw = fullWidth
		}
		diff -= d
		totalSc -= sc

		ctrl.SetSize(tw, th)
		ctrl.ResizeChildren()
	}
}

//AddChild --
func (c *BaseControl) AddChild(control мИнт.ИВиджет) {
	if c.children == nil {
		c.children = make([]мИнт.ИВиджет, 1)
		c.children[0] = control
	} else {
		if c.ChildExists(control) {
			panic("Double adding a child")
		}

		c.children = append(c.children, control)
	}

	var ctrl мИнт.ИВиджет
	var mainCtrl мИнт.ИВиджет
	ctrl = c
	for ctrl != nil {
		ww, hh := ctrl.MinimalSize()
		cw, ch := ctrl.Size()
		if ww > cw || hh > ch {
			if ww > cw {
				cw = ww
			}
			if hh > ch {
				ch = hh
			}
			ctrl.SetConstraints(cw, ch)
		}

		if ctrl.Parent() == nil {
			mainCtrl = ctrl
		}
		ctrl = ctrl.Parent()
	}

	if mainCtrl != nil {
		mainCtrl.ResizeChildren()
		mainCtrl.PlaceChildren()
	}

	if c.clipped && c.clipper == nil {
		c.setClipper()
	}
}

//Children --
func (c *BaseControl) Children() []мИнт.ИВиджет {
	child := make([]мИнт.ИВиджет, len(c.children))
	copy(child, c.children)
	return child
}

//ChildExists --
func (c *BaseControl) ChildExists(control мИнт.ИВиджет) bool {
	if len(c.children) == 0 {
		return false
	}

	for _, ctrl := range c.children {
		if ctrl == control {
			return true
		}
	}

	return false
}

//ChildrenScale --
func (c *BaseControl) ChildrenScale() int {
	if c.childCount() == 0 {
		return c.scale
	}

	total := 0
	for _, ctrl := range c.children {
		if ctrl.Visible() {
			total += ctrl.Scale()
		}
	}

	return total
}

//MinimalSize --
func (c *BaseControl) MinimalSize() (w int, h int) {
	children := c.childCount()
	if children == 0 {
		return c.minW, c.minH
	}

	totalX := 2 * c.padX
	totalY := 2 * c.padY

	if c.pack == мКнст.Vertical {
		totalY += (children - 1) * c.gapY
	} else {
		totalX += (children - 1) * c.gapX
	}

	for _, ctrl := range c.children {
		if ctrl.Clipped() {
			continue
		}

		if !ctrl.Visible() {
			continue
		}
		ww, hh := ctrl.MinimalSize()
		if c.pack == мКнст.Vertical {
			totalY += hh
			if ww+2*c.padX > totalX {
				totalX = ww + 2*c.padX
			}
		} else {
			totalX += ww
			if hh+2*c.padY > totalY {
				totalY = hh + 2*c.padY
			}
		}
	}

	if totalX < c.minW {
		totalX = c.minW
	}
	if totalY < c.minH {
		totalY = c.minH
	}

	return totalX, totalY
}

//Draw --
func (c *BaseControl) Draw() {
	panic("BaseControl Draw Called")
}

//DrawChildren --
func (c *BaseControl) DrawChildren() {
	if c.hidden {
		return
	}

	PushClip()
	defer PopClip()

	cp := ClippedParent(c)
	var cTarget мВид.ИВиджет

	cTarget = c
	if cp != nil {
		cTarget = cp
	}

	x, y, w, h := cTarget.Clipper()
	SetClipRect(x, y, w, h)

	for _, child := range c.children {
		child.Draw()
	}
}

//Clipper --
func (c *BaseControl) Clipper() (int, int, int, int) {
	clipped := ClippedParent(c)

	if clipped == nil || (c.clipped && c.clipper != nil) {
		return c.clipper.x, c.clipper.y, c.clipper.w, c.clipper.h
	}

	return CalcClipper(c)
}

func (c *BaseControl) setClipper() {
	x, y, w, h := CalcClipper(c)
	c.clipper = &rect{x: x, y: y, w: w, h: h}
}
//HitTest --
func (c *BaseControl) HitTest(x, y int) мИнт.HitResult {
	if x > c.x && x < c.x+c.width-1 &&
		y > c.y && y < c.y+c.height-1 {
		return мКнст.HitInside
	}

	if (x == c.x || x == c.x+c.width-1) &&
		y >= c.y && y < c.y+c.height {
		return мКнст.HitBorder
	}

	if (y == c.y || y == c.y+c.height-1) &&
		x >= c.x && x < c.x+c.width {
		return мКнст.HitBorder
	}

	return мКнст.HitOutside
}

//ProcessEvent --
func (c *BaseControl) ProcessEvent(ev мИнт.ИСобытие) bool {
	return SendEventToChild(c, ev)
}

//PlaceChildren --
func (c *BaseControl) PlaceChildren() {
	children := c.childCount()
	if c.children == nil || children == 0 {
		return
	}

	xx, yy := c.x+c.padX, c.y+c.padY
	for _, ctrl := range c.children {
		if !ctrl.Visible() {
			continue
		}

		ctrl.SetPos(xx, yy)
		ww, hh := ctrl.Size()
		if c.pack == мКнст.Vertical {
			yy += c.gapY + hh
		} else {
			xx += c.gapX + ww
		}

		ctrl.PlaceChildren()
	}
}

// ActiveColors return the attrubutes for the controls when it
// is active: text and background colors
func (c *BaseControl) ActiveColors() (term.Attribute, term.Attribute) {
	return c.fgActive, c.bgActive
}

// SetActiveTextColor changes text color of the active control
func (c *BaseControl) SetActiveTextColor(clr term.Attribute) {
	c.fgActive = clr
}

// SetActiveBackColor changes background color of the active control
func (c *BaseControl) SetActiveBackColor(clr term.Attribute) {
	c.bgActive = clr
}
//RemoveChild --
func (c *BaseControl) RemoveChild(control мИнт.ИВиджет) {
	children := []мИнт.ИВиджет{}

	for _, child := range c.children {
		if child.RefID() == control.RefID() {
			continue
		}

		children = append(children, child)
	}
	c.children = nil

	for _, child := range children {
		c.AddChild(child)
	}
}

// Destroy removes an object from its parental chain
func (c *BaseControl) Destroy() {
	c.parent.RemoveChild(c)
	c.parent.SetConstraints(0, 0)
}
