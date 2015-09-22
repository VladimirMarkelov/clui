package clui

import (
	"strings"
	"unicode/utf8"
)

/*
Radio is a control for selecting one item from a set of items
exlusively. Items can be displayed horizontally or verically.
The latter case is default.

To set the item list, use Props field Text when creating a Contol.
Text is a string of items separated with '|'

Radio implements a few its own methods:
GetSelectedItem, SetSelectedItem
*/
type Radio struct {
	posX, posY    int
	width, height int
	title         string
	anchor        Anchor
	id            WinId
	enabled       bool
	align         Align
	active        bool
	selected      int
	parent        Window
	direction     Direction
	items         []string
	border        BorderStyle
	itemRects     []Coord
	visible       bool
	tabStop       bool
	textColor     Color
	backColor     Color
	scale         int

	minW, minH int
}

func NewRadio(parent Window, id WinId, x, y, width, height int, title string, props Props) *Radio {
	r := new(Radio)
	r.SetEnabled(true)
	r.SetPos(x, y)

	if height < 3 {
		height = 3
	}

	r.SetSize(width, height)
	r.anchor = props.Anchors
	r.selected = -1
	r.title = title
	r.align = props.Alignment
	r.parent = parent
	r.border = props.Border
	r.direction = props.Dir
	r.items = strings.Split(props.Text, "|")
	r.visible = true
	r.tabStop = true
	r.id = id

	r.minW, r.minH = 5, 3

	return r
}

func (r *Radio) SetText(title string) {
	r.title = title
}

func (r *Radio) GetText() string {
	return r.title
}

func (r *Radio) GetId() WinId {
	return r.id
}

func (r *Radio) GetSize() (int, int) {
	return r.width, r.height
}

func (r *Radio) GetConstraints() (int, int) {
	return r.minW, r.minH
}

func (r *Radio) SetConstraints(minW, minH int) {
	if minW >= 5 {
		r.minW = minW
	}
	if minH >= 3 {
		r.minH = minH
	}
}

func (r *Radio) SetSize(width, height int) {
	width, height = ApplyConstraints(r, width, height)
	r.width = width
	r.height = height
}

func (r *Radio) GetPos() (int, int) {
	return r.posX, r.posY
}

func (r *Radio) SetPos(x, y int) {
	r.posX = x
	r.posY = y
}

func (r *Radio) Redraw(canvas Canvas) {
	x, y := r.GetPos()
	w, h := r.GetSize()

	tm := canvas.Theme()

	fg, bg := r.textColor, r.backColor
	if fg == ColorDefault {
		if r.enabled {
			fg = tm.GetSysColor(ColorControlText)
		} else {
			fg = tm.GetSysColor(ColorGrayText)
		}
	}
	if bg == ColorDefault {
		if r.active {
			bg = tm.GetSysColor(ColorControlActiveBack)
		} else {
			bg = tm.GetSysColor(ColorControlBack)
		}
	}

	canvas.ClearRect(x, y, w, h, bg)
	canvas.DrawFrame(x, y, w, h, r.border, fg, bg)
	canvas.DrawAlignedText(x+1, y, w-2, r.title, fg, bg, r.align)

	r.drawGroup(canvas, fg, bg)
}

func (r *Radio) drawRadio(canvas Canvas, x, y, w int, text string, fg, bg Color, selected bool) {
	tm := canvas.Theme()

	if w < 3 {
		return
	}

	chOpen := tm.GetSysObject(ObjRadioOpen)
	chClose := tm.GetSysObject(ObjRadioClose)
	chOn := tm.GetSysObject(ObjRadioSelected)
	chOff := tm.GetSysObject(ObjRadioUnselected)

	canvas.DrawRune(x, y, chOpen, fg, bg)
	canvas.DrawRune(x+2, y, chClose, fg, bg)
	ch := chOff
	if selected {
		ch = chOn
	}
	canvas.DrawRune(x+1, y, ch, fg, bg)

	if w < 5 {
		return
	}

	canvas.DrawText(x+4, y, w-4, text, fg, bg)
}

func (r *Radio) drawGroup(canvas Canvas, fg, bg Color) {
	r.itemRects = make([]Coord, 0)
	count := len(r.items)

	dr := r.direction
	if count < 2 {
		dr = DirVertical
	}

	if dr == DirVertical {
		height := r.height - 2
		if height < count {
			// TODO: maybe show some error text?
			return
		}

		dy := 1
		if height >= count+2 {
			dy = 2
			height -= 2
		}

		step := (height - 1) / (count - 1)
		if step == 0 {
			step = 1
		}

		longest, dx := 3, 2
		for i := 0; i < count; i++ {
			l := utf8.RuneCountInString(r.items[i])
			if l > longest {
				longest = l
			}
		}

		if longest > r.width-2-2-4 {
			dx = 1
		} else if longest < r.width-2-2-4-2 {
			dx = ((r.width - 2) - (longest + 4)) / 2
		}

		for i := 0; i < count; i++ {
			l := utf8.RuneCountInString(r.items[i]) + 4
			if l > r.width-2 {
				l = r.width - 2
			}

			c := Coord{x: dx, y: dy, w: l, h: 1}
			r.itemRects = append(r.itemRects, c)
			r.drawRadio(canvas, dx+r.posX, r.posY+dy, c.w, r.items[i], fg, bg, i == r.selected)

			dy += step
		}
	} else {
		// calculate total size
		total := 0
		for i := 0; i < count; i++ {
			total += utf8.RuneCountInString(r.items[i]) + 5
		}

		oneItem := float32(r.width-2-2) / float32(len(r.items))
		dx := float32(2)
		if total >= r.width {
			dx = float32(1)
			oneItem = float32(r.width-2) / float32(len(r.items))
		}

		if int(oneItem) < 4 {
			// TODO: show some error?
			return
		}

		offset := (r.height-2)/2 + 1

		for i := 0; i < count; i++ {
			l := utf8.RuneCountInString(r.items[i]) + 4
			if l > int(oneItem) {
				l = int(oneItem)
			}

			c := Coord{x: int(dx + 0.5), y: offset, w: l, h: 1}
			r.itemRects = append(r.itemRects, c)
			r.drawRadio(canvas, r.posX+int(dx+0.5), r.posY+offset, l, r.items[i], fg, bg, i == r.selected)

			dx += oneItem
		}
	}
}

func (r *Radio) GetEnabled() bool {
	return r.enabled
}

func (r *Radio) SetEnabled(enabled bool) {
	r.enabled = enabled
}

func (r *Radio) SetAlign(align Align) {
	r.align = align
}

func (r *Radio) GetAlign() Align {
	return r.align
}

func (r *Radio) SetAnchors(anchor Anchor) {
	r.anchor = anchor
}

func (r *Radio) GetAnchors() Anchor {
	return r.anchor
}

func (r *Radio) GetActive() bool {
	return r.active
}

func (r *Radio) SetActive(active bool) {
	r.active = active
}

func (r *Radio) GetTabStop() bool {
	return r.tabStop
}

func (r *Radio) SetTabStop(tab bool) {
	r.tabStop = tab
}

func (r *Radio) itemUnderCursor(x, y int) int {
	if len(r.itemRects) == 0 {
		return -1
	}

	for id, obj := range r.itemRects {
		px, py := obj.x+r.posX, obj.y+r.posY
		if x >= px && y >= py && x < px+obj.w && y < py+obj.h {
			return id
		}
	}

	return -1
}

func (r *Radio) ProcessEvent(event Event) bool {
	if (!r.active && event.Type == EventKey) || !r.enabled {
		return false
	}

	// TODO: should process ArrowKeys if it is active
	if event.Type == EventMouseClick || event.Type == EventMouse {
		id := r.itemUnderCursor(event.X, event.Y)

		if id != -1 && id != r.selected {
			r.selected = id
			// TODO send selectItem event
		}

		return true
	}

	return false
}

func (r *Radio) SetVisible(visible bool) {
	r.visible = visible
}

func (r *Radio) GetVisible() bool {
	return r.visible
}

func (r *Radio) GetColors() (Color, Color) {
	return r.textColor, r.backColor
}

func (r *Radio) SetTextColor(clr Color) {
	r.textColor = clr
}

func (r *Radio) SetBackColor(clr Color) {
	r.backColor = clr
}

func (r *Radio) HideChildren() {
	// nothing to do
}

func (r *Radio) GetScale() int {
	return r.scale
}

func (r *Radio) SetScale(scale int) {
	r.scale = scale
}

// Returns id of selected item, and -1 in case of
// no item is selected
func (r *Radio) GetSelectedItem() int {
	return r.selected
}

// Selects an item inside Radio.
// Returns true if item is successfully selected
// or false if id is greater than item count
func (r *Radio) SetSelectedItem(id int) bool {
	if id >= len(r.items) {
		return false
	}

	r.selected = id
	return true
}
