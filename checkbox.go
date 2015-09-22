package clui

import (
	"github.com/VladimirMarkelov/termbox-go"
)

// CheckBox control. It can be two-state one(on and off) - it is default mode - or tree-state.
// State values are 0=off, 1=on, 2=third state
// Minimal width of a checkbox cannot be less than 3
type CheckBox struct {
	posX, posY    int
	width, height int
	title         string
	anchor        Anchor
	id            WinId
	enabled       bool
	align         Align
	active        bool
	state         int // 0 - off, 1 - on, 2 - third state(?)
	allow3state   bool
	visible       bool
	tabStop       bool
	textColor     Color
	backColor     Color
	scale         int

	minW, minH int

	parent Window
}

func NewCheckBox(parent Window, id WinId, x, y, width, height int, title string, props Props) *CheckBox {
	c := new(CheckBox)
	c.SetEnabled(true)
	c.SetPos(x, y)
	c.SetSize(width, 1) // TODO: only one line checkboxes are supported at that moment
	c.anchor = props.Anchors
	c.state = 0
	c.title = title
	c.align = props.Alignment
	c.parent = parent
	c.allow3state = false
	c.visible = true
	c.tabStop = true
	c.id = id

	c.minW, c.minH = 3, 1

	return c
}

func (c *CheckBox) SetText(title string) {
	c.title = title
}

func (c *CheckBox) GetText() string {
	return c.title
}

func (c *CheckBox) GetId() WinId {
	return c.id
}

func (c *CheckBox) GetSize() (int, int) {
	return c.width, c.height
}

func (c *CheckBox) GetConstraints() (int, int) {
	return c.minW, c.minH
}

func (c *CheckBox) SetConstraints(minW, minH int) {
	if minW > DoNotChange && minW >= 3 {
		c.minW = minW
	}
	if minH > DoNotChange {
		c.minH = minH
	}
}

func (c *CheckBox) SetSize(width, height int) {
	width, height = ApplyConstraints(c, width, height)

	// TODO: support multiline
	if height > 1 {
		height = 1
	}

	c.width = width
	c.height = height
}

func (c *CheckBox) GetPos() (int, int) {
	return c.posX, c.posY
}

func (c *CheckBox) SetPos(x, y int) {
	c.posX = x
	c.posY = y
}

func (c *CheckBox) Redraw(canvas Canvas) {
	x, y := c.GetPos()
	w, h := c.GetSize()

	tm := canvas.Theme()

	fg, bg := c.textColor, c.backColor
	if fg == ColorDefault {
		if c.enabled {
			fg = tm.GetSysColor(ColorControlText)
		} else {
			fg = tm.GetSysColor(ColorGrayText)
		}
	}
	if bg == ColorDefault {
		if c.active {
			bg = tm.GetSysColor(ColorControlActiveBack)
		} else {
			bg = tm.GetSysColor(ColorControlBack)
		}
	}

	cOpen := tm.GetSysObject(ObjCheckboxOpen)
	cClose := tm.GetSysObject(ObjCheckboxClose)
	cCheck := tm.GetSysObject(ObjCheckboxChecked)
	cEmpty := tm.GetSysObject(ObjCheckboxUnchecked)
	cUnknown := tm.GetSysObject(ObjCheckboxUnknown)
	cState := []rune{cEmpty, cCheck, cUnknown}

	canvas.ClearRect(x, y, w, h, bg)
	if w < 3 {
		return
	}

	canvas.DrawRune(x, y, cOpen, fg, bg)
	canvas.DrawRune(x+2, y, cClose, fg, bg)
	canvas.DrawRune(x+1, y, cState[c.state], fg, bg)

	if w < 5 {
		return
	}

	canvas.DrawAlignedText(x+4, y, w-4, c.title, fg, bg, c.align)
}

func (c *CheckBox) GetEnabled() bool {
	return c.enabled
}

func (c *CheckBox) SetEnabled(enabled bool) {
	c.enabled = enabled
}

func (c *CheckBox) SetAlign(align Align) {
	c.align = align
}

func (c *CheckBox) GetAlign() Align {
	return c.align
}

func (c *CheckBox) SetAnchors(anchor Anchor) {
	c.anchor = anchor
}

func (c *CheckBox) GetAnchors() Anchor {
	return c.anchor
}

func (c *CheckBox) GetActive() bool {
	return c.active
}

func (c *CheckBox) SetActive(active bool) {
	c.active = active
}

func (c *CheckBox) GetTabStop() bool {
	return c.tabStop
}

func (c *CheckBox) SetTabStop(tab bool) {
	c.tabStop = tab
}

func (c *CheckBox) ProcessEvent(event Event) bool {
	if (!c.active && event.Type == EventKey) || !c.enabled {
		return false
	}

	if (event.Type == EventKey && event.Key == termbox.KeySpace) || event.Type == EventMouseClick || event.Type == EventMouse {
		if c.state == 0 {
			c.state = 1
		} else if c.state == 2 {
			c.state = 0
		} else {
			if c.allow3state {
				c.state = 2
			} else {
				c.state = 0
			}
		}
		return true
	}

	return false
}

func (c *CheckBox) SetVisible(visible bool) {
	c.visible = visible
}

func (c *CheckBox) GetVisible() bool {
	return c.visible
}

func (c *CheckBox) GetColors() (Color, Color) {
	return c.textColor, c.backColor
}

func (c *CheckBox) SetTextColor(clr Color) {
	c.textColor = clr
}

func (c *CheckBox) SetBackColor(clr Color) {
	c.backColor = clr
}

func (c *CheckBox) HideChildren() {
	// nothing to do
}

func (c *CheckBox) GetScale() int {
	return c.scale
}

func (c *CheckBox) SetScale(scale int) {
	c.scale = scale
}
