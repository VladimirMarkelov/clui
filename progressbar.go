package clui

import (
	"strings"
)

/*
ProgressBar control visualizes the progression of extended operation.

The control has two sets of colors(almost all other controls have only
one set: foreground and background colors): for filled part and for
empty one. By default colors are the same.

In addition to standard Control methods it has its own ones:
SetLimits, SetValue, Step
*/
type ProgressBar struct {
	posX, posY     int
	width, height  int
	title          string
	anchor         Anchor
	align          Align
	id             WinId
	enabled        bool
	active         bool
	direction      Direction
	min, max       int
	value          int
	visible        bool
	textColor      Color
	backColor      Color
	textEmptyColor Color
	backEmptyColor Color
	scale          int
	parent         Window

	minW, minH int
}

func NewProgressBar(parent Window, id WinId, x, y, width, height int, min, max int, props Props) *ProgressBar {
	b := new(ProgressBar)
	b.SetEnabled(true)
	b.SetPos(x, y)
	b.SetSize(width, height)
	b.anchor = props.Anchors
	b.min = min
	b.max = max
	b.direction = props.Dir
	b.visible = true
	b.minW, b.minH = 1, 1
	b.parent = parent

	return b
}

func (b *ProgressBar) SetText(title string) {
	// nothing to do
}

func (b *ProgressBar) GetText() string {
	return b.title
}

func (b *ProgressBar) GetId() WinId {
	return b.id
}

func (b *ProgressBar) GetSize() (int, int) {
	return b.width, b.height
}

func (b *ProgressBar) GetConstraints() (int, int) {
	return b.minW, b.minH
}

func (b *ProgressBar) SetConstraints(minW, minH int) {
	if minW >= 1 {
		b.minW = minW
	}
	if minH >= 1 {
		b.minH = minH
	}
}

func (b *ProgressBar) SetSize(width, height int) {
	width, height = ApplyConstraints(b, width, height)
	b.width = width
	b.height = height
}

func (b *ProgressBar) GetPos() (int, int) {
	return b.posX, b.posY
}

func (b *ProgressBar) SetPos(x, y int) {
	b.posX = x
	b.posY = y
}

func (b *ProgressBar) Redraw(canvas Canvas) {
	if b.max <= b.min {
		return
	}

	tm := canvas.Theme()

	fgOff, fgOn, bgOff, bgOn := b.textEmptyColor, b.textColor, b.backEmptyColor, b.backColor

	if fgOff == ColorDefault {
		fgOff = tm.GetSysColor(ColorProgressOff)
	}
	if fgOn == ColorDefault {
		fgOn = tm.GetSysColor(ColorProgressOn)
	}
	if bgOff == ColorDefault {
		bgOff = tm.GetSysColor(ColorProgressOffBack)
	}
	if bgOn == ColorDefault {
		bgOn = tm.GetSysColor(ColorProgressOnBack)
	}

	cFilled := tm.GetSysObject(ObjProgressBarFull)
	cEmpty := tm.GetSysObject(ObjProgressBarEmpty)

	prc := 0
	if b.value >= b.max {
		prc = 100
	} else if b.value < b.max && b.value > b.min {
		prc = (100 * (b.value - b.min)) / (b.max - b.min)
	}

	x, y := b.GetPos()
	w, h := b.GetSize()

	if b.direction == DirHorizontal {
		filled := prc * w / 100
		sFilled := strings.Repeat(string(cFilled), filled)
		sEmpty := strings.Repeat(string(cEmpty), w-filled)

		for yy := y; yy < y+h; yy++ {
			canvas.DrawText(x, yy, filled, sFilled, fgOn, bgOn)
			canvas.DrawText(x+filled, yy, w-filled, sEmpty, fgOff, bgOff)
		}
	} else {
		filled := prc * h / 100
		sFilled := strings.Repeat(string(cFilled), w)
		sEmpty := strings.Repeat(string(cEmpty), w)
		for yy := y; yy < y+h-filled; yy++ {
			canvas.DrawText(x, yy, w, sEmpty, fgOff, bgOff)
		}
		for yy := y + h - filled; yy < y+h; yy++ {
			canvas.DrawText(x, yy, w, sFilled, fgOn, bgOn)
		}
	}
}

func (b *ProgressBar) GetEnabled() bool {
	return b.enabled
}

func (b *ProgressBar) SetEnabled(active bool) {
	// nothing to do
}

func (b *ProgressBar) SetAlign(align Align) {
	// nothing
}

func (b *ProgressBar) GetAlign() Align {
	return b.align
}

func (b *ProgressBar) SetAnchors(anchor Anchor) {
	b.anchor = anchor
}

func (b *ProgressBar) GetAnchors() Anchor {
	return b.anchor
}

func (b *ProgressBar) GetActive() bool {
	return b.active
}

func (b *ProgressBar) SetActive(active bool) {
	b.active = active
}

func (b *ProgressBar) GetTabStop() bool {
	return false
}

func (b *ProgressBar) SetTabStop(tab bool) {
	// nothing
}

func (b *ProgressBar) ProcessEvent(event Event) bool {
	return false
}

func (b *ProgressBar) SetVisible(visible bool) {
	b.visible = visible
}

func (b *ProgressBar) GetVisible() bool {
	return b.visible
}

//----------------- own methods -------------------------

// Sets new progress value. If value exeeds ProgressBar
// limits then the limit value is used
func (b *ProgressBar) SetValue(pos int) {
	if pos < b.min {
		b.value = b.min
	} else if pos > b.max {
		b.value = b.max
	} else {
		b.value = pos
	}
}

// Set new ProgressBar limits. The current value is adjusted
// if it exeeds new limits
func (b *ProgressBar) SetLimits(min, max int) {
	b.min = min
	b.max = max

	if b.value < b.min {
		b.value = min
	}
	if b.value > b.max {
		b.value = max
	}
}

// Increase ProgressBar value by 1 if the value is less
// than ProgressBar high limit
func (b *ProgressBar) Step() int {
	b.value++

	if b.value > b.max {
		b.value = b.max
	}

	return b.value
}

func (b *ProgressBar) GetColors() (Color, Color) {
	return b.textColor, b.backColor
}

func (b *ProgressBar) GetSecondColors() (Color, Color) {
	return b.textEmptyColor, b.backEmptyColor
}

func (b *ProgressBar) SetSecondColors(fg, bg Color) {
	b.textEmptyColor, b.backEmptyColor = fg, bg
}

func (b *ProgressBar) SetTextColor(clr Color) {
	b.textColor = clr
}

func (b *ProgressBar) SetBackColor(clr Color) {
	b.backColor = clr
}

func (b *ProgressBar) HideChildren() {
	// nothing to do
}

func (b *ProgressBar) GetScale() int {
	return b.scale
}

func (b *ProgressBar) SetScale(scale int) {
	b.scale = scale
}
