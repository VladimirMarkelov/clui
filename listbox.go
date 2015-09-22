package clui

import (
	"github.com/VladimirMarkelov/termbox-go"
	"strings"
)

/*
Control to display a list of items and allow to user to select any of them.
Content is scrollable with arrow keys or by clicking up and bottom buttons
on the scroll(now content is scrollable with mouse dragging only on Windows).

ListBox calls onSelectItem item function after a user changes currently
selected item with mouse or using keyboard (extra case: the event is emitted
when a user presses Enter - the case is used in ComboBox to select an item
from drop down list). Event structure has 2 fields filled: Y - selected
item number in list(-1 if nothing is selected), Msg - text of the selected item.

ListBox provides a few own methods to manipulate its items:
AddItem, SelectItem, FindItem, GetSelectedItem, RemoveItem, Clear
*/
type ListBox struct {
	posX, posY    int
	width, height int
	title         string
	anchor        Anchor
	id            WinId
	enabled       bool
	align         Align
	active        bool
	parent        Window
	visible       bool
	tabStop       bool
	textColor     Color
	backColor     Color
	scale         int

	// own listbox members
	items         []string
	currSelection int
	topLine       int
	pressY        int
	buttonPos     int

	minW, minH int

	onSelectItem func(Event)
}

func NewListBox(parent Window, id WinId, x, y, width, height int, props Props) *ListBox {
	l := new(ListBox)
	l.SetEnabled(true)
	l.SetPos(x, y)
	l.SetSize(width, height)
	l.currSelection = -1
	l.items = make([]string, 0)
	l.topLine = 0
	l.pressY = -1
	l.buttonPos = -1
	l.parent = parent
	l.visible = true
	l.tabStop = true
	l.id = id
	l.minW, l.minH = 3, 5

	l.onSelectItem = nil

	return l
}

func (l *ListBox) SetText(title string) {
	l.title = title
}

func (l *ListBox) GetText() string {
	return l.title
}

func (l *ListBox) GetId() WinId {
	return l.id
}

func (l *ListBox) GetSize() (int, int) {
	return l.width, l.height
}

func (l *ListBox) GetConstraints() (int, int) {
	return l.minW, l.minH
}

func (l *ListBox) SetConstraints(minW, minH int) {
	if minW >= 3 {
		l.minW = minW
	}
	if minH >= 5 {
		l.minH = minH
	}
}

func (l *ListBox) SetSize(width, height int) {
	width, height = ApplyConstraints(l, width, height)
	l.width = width
	l.height = height
}

func (l *ListBox) GetPos() (int, int) {
	return l.posX, l.posY
}

func (l *ListBox) SetPos(x, y int) {
	l.posX = x
	l.posY = y
}

func (l *ListBox) redrawScroll(canvas Canvas, tm *ThemeManager) {
	chLine := tm.GetSysObject(ObjScrollBar)
	chCursor := tm.GetSysObject(ObjScrollThumb)
	chUp := tm.GetSysObject(ObjScrollUpArrow)
	chDown := tm.GetSysObject(ObjScrollDownArrow)

	fg, bg, fgThumb := ColorDefault, ColorDefault, ColorDefault
	if fg == ColorDefault {
		fg = tm.GetSysColor(ColorScroll)
	}
	if bg == ColorDefault {
		bg = tm.GetSysColor(ColorScrollBack)
	}
	if fgThumb == ColorDefault {
		fgThumb = tm.GetSysColor(ColorScrollThumb)
	}

	canvas.DrawRune(l.posX+l.width-1, l.posY, chUp, fg, bg)
	canvas.DrawRune(l.posX+l.width-1, l.posY+l.height-1, chDown, fg, bg)

	if l.height > 2 {
		for yy := 1; yy < l.height-1; yy++ {
			canvas.DrawRune(l.posX+l.width-1, l.posY+yy, chLine, fg, bg)
		}
	}

	if l.currSelection == -1 {
		return
	}

	if l.height == 3 || l.currSelection <= 0 {
		canvas.DrawRune(l.posX+l.width-1, l.posY+1, chCursor, fgThumb, bg)
		return
	}

	if l.pressY == -1 {
		ydiff := int(float32(l.currSelection) / float32(len(l.items)-1.0) * float32(l.height-3))
		l.buttonPos = ydiff + 1
	}
	canvas.DrawRune(l.posX+l.width-1, l.posY+l.buttonPos, chCursor, fgThumb, bg)
}

func (l *ListBox) redrawItems(canvas Canvas, tm *ThemeManager) {
	maxCurr := len(l.items) - 1
	curr := l.topLine
	dy := 0
	maxDy := l.height - 1
	maxWidth := l.width - 1

	fg, bg := l.textColor, l.backColor
	fgSel, bgSel := ColorDefault, ColorDefault
	if fg == ColorDefault {
		fg = tm.GetSysColor(ColorEditText)
	}
	if bg == ColorDefault {
		bg = tm.GetSysColor(ColorEditBack)
	}
	if fgSel == ColorDefault {
		fgSel = tm.GetSysColor(ColorSelectionText)
	}
	if bgSel == ColorDefault {
		bgSel = tm.GetSysColor(ColorSelectionBack)
	}

	for curr <= maxCurr && dy <= maxDy {
		f, b := fg, bg
		if curr == l.currSelection {
			f, b = fgSel, bgSel
		}
		canvas.DrawText(l.posX, l.posY+dy, maxWidth, l.items[curr], f, b)

		curr++
		dy++
	}
}

func (l *ListBox) Redraw(canvas Canvas) {
	x, y := l.GetPos()
	w, h := l.GetSize()

	tm := canvas.Theme()
	bg := ColorDefault
	if bg == ColorDefault {
		bg = tm.GetSysColor(ColorEditBack)
	}

	canvas.ClearRect(x, y, w, h, bg)
	l.redrawItems(canvas, tm)
	l.redrawScroll(canvas, tm)
}

func (l *ListBox) GetEnabled() bool {
	return l.enabled
}

func (l *ListBox) SetEnabled(active bool) {
	l.enabled = active
}

func (l *ListBox) SetAlign(align Align) {
	l.align = align
}

func (l *ListBox) GetAlign() Align {
	return l.align
}

func (l *ListBox) SetAnchors(anchor Anchor) {
	l.anchor = anchor
}

func (l *ListBox) GetAnchors() Anchor {
	return l.anchor
}

func (l *ListBox) GetActive() bool {
	return l.active
}

func (l *ListBox) SetActive(active bool) {
	l.active = active
}

func (l *ListBox) GetTabStop() bool {
	return l.tabStop
}

func (l *ListBox) SetTabStop(tab bool) {
	l.tabStop = tab
}

func (l *ListBox) home() {
	if len(l.items) > 0 {
		l.currSelection = 0
	}
	l.topLine = 0
}

func (l *ListBox) end() {
	length := len(l.items)

	if length == 0 {
		return
	}

	l.currSelection = length - 1
	if length > l.height {
		l.topLine = length - l.height
	}
}

func (l *ListBox) moveUp() {
	if l.topLine == 0 && l.currSelection == 0 {
		return
	}

	if l.currSelection == -1 {
		if len(l.items) != 0 {
			l.currSelection = 0
		}
		return
	}

	l.currSelection--
	l.ensureVisible()
}

func (l *ListBox) moveDown() {
	length := len(l.items)

	if length == 0 || l.currSelection == length-1 {
		return
	}

	l.currSelection++
	l.ensureVisible()
}

func (l *ListBox) ensureVisible() {
	length := len(l.items)

	if length <= l.height {
		return
	}

	diff := l.currSelection - l.topLine
	if diff >= 0 && diff < l.height {
		return
	}

	if diff < 0 {
		l.topLine = l.currSelection
	} else {
		top := l.currSelection - l.height + 1
		if length-top > l.height {
			l.topLine = top
		} else {
			l.topLine = length - l.height
		}
	}
}

// Deletes all ListBox items
func (l *ListBox) Clear() {
	l.items = make([]string, 0)
	l.currSelection = -1
	l.topLine = 0
}

func (l *ListBox) processMouseClick(ev Event) bool {
	if ev.Key != termbox.MouseLeft {
		return false
	}

	dx := ev.X - l.posX
	dy := ev.Y - l.posY

	if dx == l.width-1 {
		if dy < 0 || dy >= l.height || len(l.items) < 2 {
			return true
		}

		if dy == 0 {
			l.moveUp()
			return true
		}
		if dy == l.height-1 {
			l.moveDown()
			return true
		}

		l.buttonPos = dy
		l.recalcPositionByScroll()
		return true
	}

	if dx < 0 || dx >= l.width || dy < 0 || dy >= l.height {
		return true
	}

	if dy >= len(l.items) {
		return true
	}

	l.SelectItem(l.topLine + dy)
	if l.onSelectItem != nil {
		ev := Event{Y: l.topLine + dy, Msg: l.GetSelectedItem()}
		go l.onSelectItem(ev)
	}

	return true
}

func (l *ListBox) processMousePress(ev Event) bool {
	if ev.Key != termbox.MouseLeft {
		return false
	}

	dx := ev.X - l.posX
	dy := ev.Y - l.posY

	if dx != l.width-1 || len(l.items) < 2 || dy != l.buttonPos {
		return true
	}

	l.pressY = ev.Y
	return true
}

func (l *ListBox) processMouseRelease(ev Event) bool {
	if ev.Key == termbox.MouseLeft {
		l.pressY = -1
		return true
	}

	return false
}

func (l *ListBox) recalcPositionByScroll() {
	if len(l.items) < 2 {
		return
	}

	newPos := int(float32(len(l.items)-1)*float32(l.buttonPos-1)/float32(l.height-3) + 0.9)

	if newPos < 0 {
		newPos = 0
	} else if newPos >= len(l.items) {
		newPos = len(l.items) - 1
	}

	l.currSelection = newPos
	l.ensureVisible()
}

func (l *ListBox) processMouseMove(ev Event) bool {
	if l.pressY == -1 {
		return false
	}

	length := len(l.items)
	if length < 2 {
		return true
	}

	dy := ev.Y - l.pressY
	l.pressY = ev.Y

	if dy > 0 && l.buttonPos == l.height-2 {
		return true
	}
	if dy < 0 && l.buttonPos == 1 {
		return true
	}

	l.buttonPos += dy
	if l.buttonPos < 1 {
		l.buttonPos = 1
	} else if l.buttonPos >= l.height-1 {
		l.buttonPos = l.height - 2
	}
	l.recalcPositionByScroll()

	return true
}

func (l *ListBox) ProcessEvent(event Event) bool {
	if !l.active || !l.enabled {
		return false
	}

	switch event.Type {
	case EventKey:
		switch event.Key {
		case termbox.KeyHome:
			l.home()
			return true
		case termbox.KeyEnd:
			l.end()
			return true
		case termbox.KeyArrowUp:
			l.moveUp()
			return true
		case termbox.KeyArrowDown:
			l.moveDown()
			return true
		case termbox.KeyCtrlM:
			if l.currSelection != -1 && l.onSelectItem != nil {
				ev := Event{Y: l.currSelection, Msg: l.GetSelectedItem()}
				go l.onSelectItem(ev)
			}
		default:
			return false
		}
	case EventMouseScroll:
		if event.Y > 0 {
			l.moveDown()
			return true
		} else if event.Y < 0 {
			l.moveUp()
			return true
		}
	case EventMouseClick, EventMouse:
		return l.processMouseClick(event)
	case EventMousePress:
		return l.processMousePress(event)
	case EventMouseRelease:
		return l.processMouseRelease(event)
	case EventMouseMove:
		return l.processMouseMove(event)
	}

	return false
}

func (l *ListBox) SetVisible(visible bool) {
	l.visible = visible
}

func (l *ListBox) GetVisible() bool {
	return l.visible
}

// own methods

// Adds a new item to item list
// Returns true if the operation is successful
func (l *ListBox) AddItem(item string) bool {
	l.items = append(l.items, item)
	return true
}

// Selects item which number in the list equals id. If the item exists the
// ListBox scrolls the list to make the item visible.
// Returns true if the item is selected successfully
func (l *ListBox) SelectItem(id int) bool {
	if len(l.items) <= id || id < 0 {
		return false
	}

	l.currSelection = id
	l.ensureVisible()
	return true
}

// Finds an item in list which text equals to text, by default the search
// is casesensitive.
// Returns item number in item list or -1 if nothing is found.
func (l *ListBox) FindItem(text string, caseSensitive bool) int {
	for idx, itm := range l.items {
		if itm == text || (caseSensitive && strings.EqualFold(itm, text)) {
			return idx
		}
	}

	return -1
}

// Returns text of currently selected item or empty sting if nothing is
// selected or ListBox is empty
func (l *ListBox) GetSelectedItem() string {
	if l.currSelection == -1 {
		return ""
	}

	return l.items[l.currSelection]
}

// Deletes an item which number is id in item list
// Returns true if item is deleted
func (l *ListBox) RemoveItem(id int) bool {
	if id < 0 || id >= len(l.items) {
		return false
	}

	l.items = append(l.items[:id], l.items[id+1:]...)
	return true
}

func (l *ListBox) OnSelectItem(fn func(Event)) {
	l.onSelectItem = fn
}

func (l *ListBox) GetColors() (Color, Color) {
	return l.textColor, l.backColor
}

func (l *ListBox) SetTextColor(clr Color) {
	l.textColor = clr
}

func (l *ListBox) SetBackColor(clr Color) {
	l.backColor = clr
}

func (l *ListBox) HideChildren() {
	// nothing to do
}

func (l *ListBox) GetScale() int {
	return l.scale
}

func (l *ListBox) SetScale(scale int) {
	l.scale = scale
}
