package clui

import (
	term "github.com/nsf/termbox-go"
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
	ControlBase
	// own listbox members
	items         []string
	currSelection int
	topLine       int
	maxItems      int
	bgSel, fgSel  term.Attribute
	buttonPos     int

	onSelectItem func(Event)
}

func NewListBox(view View, parent Control, width, height int, scale int) *ListBox {
	l := new(ListBox)
	l.SetSize(width, height)
	l.SetConstraints(width, height)
	l.currSelection = -1
	l.items = make([]string, 0)
	l.topLine = 0
	l.parent = parent
	l.view = view
	l.maxItems = 0
	l.buttonPos = -1

	l.fg = ColorBlack
	l.bg = ColorWhite
	l.fgSel = ColorYellow
	l.bgSel = ColorBlue
	l.SetTabStop(true)

	l.onSelectItem = nil

	if parent != nil {
		parent.AddChild(l, scale)
	}

	return l
}

func (l *ListBox) redrawScroll(canvas Canvas, tm Theme) {
	parts := []rune(tm.SysObject(ObjScrollBar))

	chLine, chCursor, chUp, chDown := parts[0], parts[1], parts[2], parts[3]

	fg, bg := RealColor(tm, l.fg, ColorScrollText), RealColor(tm, l.bg, ColorScrollBack)
	fgThumb, bgThumb := RealColor(tm, l.fg, ColorThumbText), RealColor(tm, l.bg, ColorThumbBack)

	canvas.PutSymbol(l.x+l.width-1, l.y, term.Cell{Ch: chUp, Fg: fg, Bg: bg})
	canvas.PutSymbol(l.x+l.width-1, l.y+l.height-1, term.Cell{Ch: chDown, Fg: fg, Bg: bg})

	if l.height > 2 {
		for yy := 1; yy < l.height-1; yy++ {
			canvas.PutSymbol(l.x+l.width-1, l.y+yy, term.Cell{Ch: chLine, Fg: fg, Bg: bg})
		}
	}

	if l.currSelection == -1 {
		return
	}

	if l.height == 3 || l.currSelection <= 0 {
		canvas.PutSymbol(l.x+l.width-1, l.y+1, term.Cell{Ch: chCursor, Fg: fgThumb, Bg: bgThumb})
		return
	}

	// if l.pressY == -1 {
	ydiff := int(float32(l.currSelection) / float32(len(l.items)-1.0) * float32(l.height-3))
	l.buttonPos = ydiff + 1
	// }
	canvas.PutSymbol(l.x+l.width-1, l.y+l.buttonPos, term.Cell{Ch: chCursor, Fg: fgThumb, Bg: bgThumb})
}

func (l *ListBox) redrawItems(canvas Canvas, tm Theme) {
	maxCurr := len(l.items) - 1
	curr := l.topLine
	dy := 0
	maxDy := l.height - 1
	maxWidth := l.width - 1

	fg, bg := RealColor(tm, l.fg, ColorEditText), RealColor(tm, l.bg, ColorEditBack)
	fgSel, bgSel := RealColor(tm, l.fgSel, ColorSelectionText), RealColor(tm, l.bgSel, ColorSelectionBack)

	for curr <= maxCurr && dy <= maxDy {
		f, b := fg, bg
		if curr == l.currSelection {
			f, b = fgSel, bgSel
		}

		_, text := AlignText(l.items[curr], maxWidth, AlignLeft)
		canvas.PutText(l.x, l.y+dy, text, f, b)

		curr++
		dy++
	}
}

func (l *ListBox) Repaint() {
	canvas := l.view.Canvas()
	tm := l.view.Screen().Theme()

	x, y := l.Pos()
	w, h := l.Size()

	bg := RealColor(tm, l.bg, ColorEditText)
	canvas.FillRect(x, y, w, h, term.Cell{Bg: bg, Ch: ' '})
	l.redrawItems(canvas, tm)
	l.redrawScroll(canvas, tm)
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
	if ev.Key != term.MouseLeft {
		return false
	}

	dx := ev.X - l.x
	dy := ev.Y - l.y

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

func (l *ListBox) ProcessEvent(event Event) bool {
	if !l.Active() || !l.Enabled() {
		return false
	}

	switch event.Type {
	case EventKey:
		switch event.Key {
		case term.KeyHome:
			l.home()
			return true
		case term.KeyEnd:
			l.end()
			return true
		case term.KeyArrowUp:
			l.moveUp()
			return true
		case term.KeyArrowDown:
			l.moveDown()
			return true
		case term.KeyCtrlM:
			if l.currSelection != -1 && l.onSelectItem != nil {
				ev := Event{Y: l.currSelection, Msg: l.GetSelectedItem()}
				go l.onSelectItem(ev)
			}
		default:
			return false
		}
	case EventMouse:
		return l.processMouseClick(event)
	}

	return false
}

// own methods

// Adds a new item to item list
// Returns true if the operation is successful
func (l *ListBox) AddItem(item string) bool {
	if l.maxItems > 0 && len(l.items) > l.maxItems {
		l.RemoveItem(0)
	}

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

func (l *ListBox) MaxItems() int {
	return l.maxItems
}

func (l *ListBox) SetMaxItems(max int) {
	l.maxItems = max
}
