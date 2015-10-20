package clui

import (
	"fmt"
	term "github.com/nsf/termbox-go"
	"log"
)

type ControlBase struct {
	title           string
	x, y            int
	width, height   int
	minW, minH      int
	scale           int
	fg, bg          term.Attribute
	fgActive        term.Attribute
	bgActive        term.Attribute
	tabSkip         bool
	disabled        bool
	align           Align
	parent          Control
	view            View
	active          bool
	padX, padY      int
	padTop, padSide int
}

func (c *ControlBase) Title() string {
	return c.title
}

func (c *ControlBase) SetTitle(title string) {
	c.title = title
}

func (c *ControlBase) Size() (int, int) {
	return c.width, c.height
}

func (c *ControlBase) SetSize(width, height int) {
	if width != DoNotChange && (width > 1000 || width < c.minW) {
		panic(fmt.Sprintf("Invalid width: %v", width))
	}
	if height != DoNotChange && (height > 200 || height < c.minH) {
		panic(fmt.Sprintf("Invalid height: %v", height))
	}

	if width != DoNotChange {
		c.width = width
	}
	if height != DoNotChange {
		c.height = height
	}
}

func (c *ControlBase) Pos() (int, int) {
	return c.x, c.y
}

func (c *ControlBase) SetPos(x, y int) {
	c.x = x
	c.y = y
}

func (c *ControlBase) ApplyConstraints() {
	w, h := c.Size()
	wM, hM := c.Constraints()

	newW, newH := w, h
	if w < wM {
		newW = wM
	}
	if h < hM {
		newH = hM
	}

	if newW != w || newH != h {
		c.SetSize(newW, newH)
	}
}

func (c *ControlBase) SetConstraints(width, height int) {
	if width >= 1 {
		c.minW = width
	}
	if height >= 1 {
		c.minH = height
	}

	c.ApplyConstraints()
}

func (c *ControlBase) Constraints() (int, int) {
	return c.minW, c.minH
}

func (c *ControlBase) Scale() int {
	return c.scale
}

func (c *ControlBase) SetScale(scale int) {
	c.scale = scale
}

func (c *ControlBase) Pack() PackType {
	return Vertical
}

func (c *ControlBase) SetPack(pk PackType) {
}

func (c *ControlBase) AddChild(ctrl Control, scale int) {
	panic("This control cannot have children")
}

func (c *ControlBase) Children() []Control {
	return make([]Control, 0)
}

func (c *ControlBase) Colors() (term.Attribute, term.Attribute) {
	return c.fg, c.bg
}

func (c *ControlBase) SetTextColor(clr term.Attribute) {
	c.fg = clr
}

func (c *ControlBase) SetBackColor(clr term.Attribute) {
	c.bg = clr
}

func (c *ControlBase) ActiveColors() (term.Attribute, term.Attribute) {
	return c.fg, c.bg
}

func (c *ControlBase) SetActiveTextColor(clr term.Attribute) {
	c.fg = clr
}

func (c *ControlBase) SetActiveBackColor(clr term.Attribute) {
	c.bg = clr
}

func (c *ControlBase) TabStop() bool {
	return !c.tabSkip
}

func (c *ControlBase) SetTabStop(skip bool) {
	c.tabSkip = !skip
}

func (c *ControlBase) Enabled() bool {
	return !c.disabled
}

func (c *ControlBase) SetEnabled(enable bool) {
	c.disabled = !enable
}

func (c *ControlBase) SetAlign(align Align) {
	c.align = align
}

func (c *ControlBase) GetAlign() Align {
	return c.align
}

func (c *ControlBase) Active() bool {
	return c.active
}

func (c *ControlBase) SetActive(active bool) {
	c.active = active
}

func (c *ControlBase) ProcessEvent(ev Event) bool {
	return false
}

func (c *ControlBase) Parent() Control {
	return c.parent
}

func (c *ControlBase) RecalculateConstraints() {
}

func (c *ControlBase) Paddings() (int, int, int, int) {
	return c.padSide, c.padTop, c.padX, c.padY
}

func (c *ControlBase) SetPaddings(side, top, dx, dy int) {
	if side >= 0 {
		c.padSide = side
	}
	if top >= 0 {
		c.padTop = top
	}
	if dx >= 0 {
		c.padX = dx
	}
	if dy >= 0 {
		c.padY = dy
	}
}

//---------- debug ----------------
func (c *ControlBase) Logger() *log.Logger {
	if c.parent == nil {
		return nil
	} else {
		return c.parent.Logger()
	}
}
