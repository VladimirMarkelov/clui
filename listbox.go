package clui

import (
	term "github.com/nsf/termbox-go"
	"strings"
	мКнст "./пакКонстанты"
)

/*
ListBox is control to display a list of items and allow to user to select any of them.
Content is scrollable with arrow keys or by clicking up and bottom buttons
on the scroll(now content is scrollable with mouse dragging only on Windows).

ListBox calls onSelectItem item function after a user changes currently
selected item with mouse or using keyboard (extra case: the event is emitted
when a user presses Enter - the case is used in ComboBox to select an item
from drop down list). Event structure has 2 fields filled: Y - selected
item number in list(-1 if nothing is selected), Msg - text of the selected item.
*/
type ListBox struct {
	BaseControl
	// own listbox members
	items         []string
	currSelection int
	topLine       int
	buttonPos     int

	onSelectItem func(мКнст.Event)
	onKeyPress   func(term.Key) bool
}

/*
NewListBox creates a new frame.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width and heigth - are minimal size of the control.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func CreateListBox(parent Control, width, height int, scale int) *ListBox {
	l := new(ListBox)
	l.BaseControl = NewBaseControl()

	if height == мКнст.AutoSize {
		height = 3
	}
	if width == мКнст.AutoSize {
		width = 5
	}

	l.SetSize(width, height)
	l.SetConstraints(width, height)
	l.currSelection = -1
	l.items = make([]string, 0)
	l.topLine = 0
	l.parent = parent
	l.buttonPos = -1

	l.SetTabStop(true)
	l.SetScale(scale)

	l.onSelectItem = nil

	if parent != nil {
		parent.AddChild(l)
	}

	return l
}

func (l *ListBox) drawScroll() {
	PushAttributes()
	defer PopAttributes()

	pos := ThumbPosition(l.currSelection, len(l.items), l.height)
	l.buttonPos = pos

	DrawScrollBar(l.x+l.width-1, l.y, 1, l.height, pos)
}

func (l *ListBox) drawItems() {
	PushAttributes()
	defer PopAttributes()

	maxCurr := len(l.items) - 1
	curr := l.topLine
	dy := 0
	maxDy := l.height - 1
	maxWidth := l.width - 1

	fg, bg := RealColor(l.fg, l.Style(), мКнст.ColorEditText), RealColor(l.bg, l.Style(), мКнст.ColorEditBack)
	if l.Active() {
		fg, bg = RealColor(l.fg, l.Style(), мКнст.ColorEditActiveText), RealColor(l.bg, l.Style(), мКнст.ColorEditActiveBack)
	}
	fgSel, bgSel := RealColor(l.fgActive, l.Style(), мКнст.ColorSelectionText), RealColor(l.bgActive, l.Style(), мКнст.ColorSelectionBack)

	for curr <= maxCurr && dy <= maxDy {
		f, b := fg, bg
		if curr == l.currSelection {
			f, b = fgSel, bgSel
		}

		SetTextColor(f)
		SetBackColor(b)
		FillRect(l.x, l.y+dy, l.width-1, 1, ' ')
		str := SliceColorized(l.items[curr], 0, maxWidth)
		DrawText(l.x, l.y+dy, str)

		curr++
		dy++
	}
}

// Repaint draws the control on its View surface
func (l *ListBox) Draw() {
	if l.hidden {
		return
	}

	PushAttributes()
	defer PopAttributes()

	x, y := l.Pos()
	w, h := l.Size()

	fg, bg := RealColor(l.fg, l.Style(), мКнст.ColorEditText), RealColor(l.bg, l.Style(), мКнст.ColorEditBack)
	if l.Active() {
		fg, bg = RealColor(l.fg, l.Style(), мКнст.ColorEditActiveText), RealColor(l.bg, l.Style(), мКнст.ColorEditActiveBack)
	}
	SetTextColor(fg)
	SetBackColor(bg)
	FillRect(x, y, w, h, ' ')
	l.drawItems()
	l.drawScroll()
}

func (l *ListBox) home() {
	if l.currSelection == 0 {
		return
	}

	if len(l.items) > 0 {
		l.currSelection = 0
	}
	l.topLine = 0

	if l.onSelectItem != nil {
		ev := мКнст.Event{Y: l.currSelection, Msg: l.SelectedItemText()}
		l.onSelectItem(ev)
	}
}

func (l *ListBox) end() {
	length := len(l.items)

	if length == 0 || l.currSelection == length-1 {
		return
	}

	l.currSelection = length - 1
	if length > l.height {
		l.topLine = length - l.height
	}

	if l.onSelectItem != nil {
		ev := мКнст.Event{Y: l.currSelection, Msg: l.SelectedItemText()}
		l.onSelectItem(ev)
	}
}

func (l *ListBox) moveUp(dy int) {
	if l.topLine == 0 && l.currSelection == 0 {
		return
	}

	if l.currSelection == -1 {
		if len(l.items) != 0 {
			l.currSelection = 0
		}
		return
	}

	if l.currSelection < dy {
		l.currSelection = 0
	} else {
		l.currSelection -= dy
	}

	l.EnsureVisible()

	if l.onSelectItem != nil {
		ev :=мКнст.Event{Y: l.currSelection, Msg: l.SelectedItemText()}
		l.onSelectItem(ev)
	}
}

func (l *ListBox) moveDown(dy int) {
	length := len(l.items)

	if length == 0 || l.currSelection == length-1 {
		return
	}

	if l.currSelection+dy >= length {
		l.currSelection = length - 1
	} else {
		l.currSelection += dy
	}

	l.EnsureVisible()

	if l.onSelectItem != nil {
		ev := мКнст.Event{Y: l.currSelection, Msg: l.SelectedItemText()}
		l.onSelectItem(ev)
	}
}

// EnsureVisible makes the currently selected item visible and scrolls the item list if it is required
func (l *ListBox) EnsureVisible() {
	length := len(l.items)

	if length <= l.height || l.currSelection == -1 {
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

// Clear deletes all ListBox items
func (l *ListBox) Clear() {
	l.items = make([]string, 0)
	l.currSelection = -1
	l.topLine = 0
}

func (l *ListBox) processMouseClick(ev мКнст.Event) bool {
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
			l.moveUp(1)
			return true
		}
		if dy == l.height-1 {
			l.moveDown(1)
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
	WindowManager().BeginUpdate()
	onSelFunc := l.onSelectItem
	WindowManager().EndUpdate()
	if onSelFunc != nil {
		ev := мКнст.Event{Y: l.topLine + dy, Msg: l.SelectedItemText()}
		onSelFunc(ev)
	}

	return true
}

func (l *ListBox) recalcPositionByScroll() {
	newPos := ItemByThumbPosition(l.buttonPos, len(l.items), l.height)
	if newPos < 1 {
		return
	}

	l.currSelection = newPos
	l.EnsureVisible()
}

/*
ProcessEvent processes all events come from the control parent. If a control
processes an event it should return true. If the method returns false it means
that the control do not want or cannot process the event and the caller sends
the event to the control parent
*/
func (l *ListBox) ProcessEvent(event мКнст.Event) bool {
	if !l.Active() || !l.Enabled() {
		return false
	}

	switch event.Type {
	case мКнст.EventKey:
		if l.onKeyPress != nil {
			res := l.onKeyPress(event.Key)
			if res {
				return true
			}
		}

		switch event.Key {
		case term.KeyHome:
			l.home()
			return true
		case term.KeyEnd:
			l.end()
			return true
		case term.KeyArrowUp:
			l.moveUp(1)
			return true
		case term.KeyArrowDown:
			l.moveDown(1)
			return true
		case term.KeyPgdn:
			l.moveDown(l.height)
			return true
		case term.KeyPgup:
			l.moveUp(l.height)
			return true
		case term.KeyCtrlM:
			if l.currSelection != -1 && l.onSelectItem != nil {
				ev := мКнст.Event{Y: l.currSelection, Msg: l.SelectedItemText()}
				l.onSelectItem(ev)
			}
		default:
			return false
		}
	case мКнст.EventMouse:
		return l.processMouseClick(event)
	}

	return false
}

// own methods

// AddItem adds a new item to item list.
// Returns true if the operation is successful
func (l *ListBox) AddItem(item string) bool {
	l.items = append(l.items, item)
	return true
}

// SelectItem selects item which number in the list equals
// id. If the item exists the ListBox scrolls the list to
// make the item visible.
// Returns true if the item is selected successfully
func (l *ListBox) SelectItem(id int) bool {
	if len(l.items) <= id || id < 0 {
		return false
	}

	l.currSelection = id
	l.EnsureVisible()
	return true
}

// Item returns item text by its index.
// If index is out of range an empty string and false are returned
func (l *ListBox) Item(id int) (string, bool) {
	if len(l.items) <= id || id < 0 {
		return "", false
	}

	return l.items[id], true
}

// FindItem looks for an item in list which text equals
// to text, by default the search is casesensitive.
// Returns item number in item list or -1 if nothing is found.
func (l *ListBox) FindItem(text string, caseSensitive bool) int {
	for idx, itm := range l.items {
		if itm == text || (caseSensitive && strings.EqualFold(itm, text)) {
			return idx
		}
	}

	return -1
}

// PartialFindItem looks for an item in list which text starts from
// the given substring, by default the search is casesensitive.
// Returns item number in item list or -1 if nothing is found.
func (l *ListBox) PartialFindItem(text string, caseSensitive bool) int {
	if !caseSensitive {
		text = strings.ToLower(text)
	}

	for idx, itm := range l.items {
		if caseSensitive {
			if strings.HasPrefix(itm, text) {
				return idx
			}
		} else {
			low := strings.ToLower(itm)
			if strings.HasPrefix(low, text) {
				return idx
			}
		}
	}

	return -1
}

// SelectedItem returns currently selected item id
func (l *ListBox) SelectedItem() int {
	return l.currSelection
}

// SelectedItemText returns text of currently selected item or empty sting if nothing is
// selected or ListBox is empty.
func (l *ListBox) SelectedItemText() string {
	if l.currSelection == -1 {
		return ""
	}

	return l.items[l.currSelection]
}

// RemoveItem deletes an item which number is id in item list
// Returns true if item is deleted
func (l *ListBox) RemoveItem(id int) bool {
	if id < 0 || id >= len(l.items) {
		return false
	}

	l.items = append(l.items[:id], l.items[id+1:]...)
	return true
}

// OnSelectItem sets a callback that is called every time
// the selected item is changed
func (l *ListBox) OnSelectItem(fn func(мКнст.Event)) {
	l.onSelectItem = fn
}

// OnKeyPress sets the callback that is called when a user presses a Key while
// the controls is active. If a handler processes the key it should return
// true. If handler returns false it means that the default handler will
// process the key
func (l *ListBox) OnKeyPress(fn func(term.Key) bool) {
	l.onKeyPress = fn
}

// ItemCount returns the number of items in the ListBox
func (l *ListBox) ItemCount() int {
	return len(l.items)
}
