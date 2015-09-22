package clui

import (
	"github.com/VladimirMarkelov/termbox-go"
	xs "github.com/huandu/xstrings"
	"strings"
)

/*
Text edit field contol. Can be simple edit field(default mode) and edit field
with drop down list. The EditField mode is set during creation and cannot be changed on the fly.
Edit field consumes some keyboard events when it is active: all printable charaters;
Delete, BackSpace, Home, End, left and right arrows; Enter, up and down arrows if EditField
in combobox mode and drop down list is visible; F5 to open drop down list in combobox mode;
Ctrl+R to clear EditField.
Edit text can be limited. By default a user can enter text of any lenght. Use SetMaxWidth to limit the maximum text length. If the text is longer than maximun then the text is automatically truncated.
EditField call funtion onChage in case of its text is changed. Event field Msg contains the new text
*/
type EditField struct {
	posX, posY int
	width      int
	title      string
	anchor     Anchor
	id         WinId
	enabled    bool
	align      Align
	active     bool
	// cursor position in edit text
	cursorPos int
	// the number of the first displayed text character - it is used in case of text is longer than edit width
	offset      int
	readonly    bool
	visible     bool
	tabStop     bool
	editBoxMode EditBoxMode
	parent      Window
	maxWidth    int
	textColor   Color
	backColor   Color
	scale       int

	minW, minH int

	// items below are used only in combobox mode
	list           *ListBox
	listboxVisible bool
	listWidth      int
	listHeight     int
	itemList       []string

	onChange func(Event)
}

func NewEditField(parent Window, id WinId, x, y, width int, text string, props Props) *EditField {
	e := new(EditField)
	e.onChange = nil
	e.SetText(text)
	e.SetEnabled(true)
	e.SetPos(x, y)
	e.SetSize(width, 1)
	e.cursorPos = xs.Len(text)
	e.offset = 0
	e.visible = true
	e.tabStop = true
	e.id = id
	e.editBoxMode = EditBoxSimple
	e.parent = parent
	e.minW, e.minH = 3, 1
	e.anchor = props.Anchors

	e.end()

	return e
}

func NewComboBox(parent Window, id WinId, x, y, width int, text string, props Props) *EditField {
	e := new(EditField)
	e.onChange = nil
	e.SetText(text)
	e.SetEnabled(true)
	e.SetPos(x, y)
	e.SetSize(width, 1)
	e.cursorPos = xs.Len(text)
	e.offset = 0
	e.visible = true
	e.tabStop = true
	e.id = id
	e.editBoxMode = EditBoxCombo
	e.SetEnabled(true)
	e.readonly = props.ReadOnly

	e.listboxVisible = false
	e.listWidth, e.listHeight = -1, -1
	e.itemList = strings.Split(props.Text, "|")
	e.parent = parent
	e.visible = true
	e.minW, e.minH = 4, 1
	e.anchor = props.Anchors

	e.end()

	return e
}

func (e *EditField) OnChange(fn func(Event)) {
	e.onChange = fn
}

func (e *EditField) SetText(title string) {
	if e.title != title {
		e.title = title
		if e.onChange != nil {
			ev := Event{Msg: title}
			go e.onChange(ev)
		}
	}
}

func (e *EditField) GetText() string {
	return e.title
}

func (e *EditField) GetId() WinId {
	return e.id
}

func (e *EditField) GetSize() (int, int) {
	return e.width, 1
}

func (e *EditField) GetConstraints() (int, int) {
	return e.minW, e.minH
}

func (e *EditField) SetConstraints(minW, minH int) {
	if minW <= DoNotChange {
		e.minW = minW
	}
	if minH <= DoNotChange {
		e.minH = minH
	}
}

func (e *EditField) SetSize(width, height int) {
	width, height = ApplyConstraints(e, width, height)
	e.width = width
}

func (e *EditField) GetPos() (int, int) {
	return e.posX, e.posY
}

func (e *EditField) SetPos(x, y int) {
	e.posX = x
	e.posY = y
}

func (e *EditField) drawButton(canvas Canvas) {
	tm := canvas.Theme()
	arrow := tm.GetSysObject(ObjComboboxDropDown)

	fg, bg := ColorDefault, ColorDefault
	if fg == ColorDefault {
		fg = tm.GetSysColor(ColorControlText)
	}
	if bg == ColorDefault {
		bg = tm.GetSysColor(ColorControlBack)
	}

	canvas.DrawRune(e.posX+e.width-1, e.posY, arrow, fg, bg)
}

func (e *EditField) Redraw(canvas Canvas) {
	x, y := e.GetPos()
	w, _ := e.GetSize()

	dw := 0
	if e.editBoxMode == EditBoxCombo {
		dw = 1
	}

	tm := canvas.Theme()
	chLeft := string(tm.GetSysObject(ObjEditLeftArrow))
	chRight := string(tm.GetSysObject(ObjEditRightArrow))

	var textOut string
	curOff := 0
	if e.offset == 0 && xs.Len(e.title) < e.width-dw {
		textOut = e.title
	} else {
		fromIdx := 0
		toIdx := 0
		if e.offset == 0 {
			toIdx = e.width - dw - 1
			textOut = xs.Slice(e.title, 0, toIdx) + chRight
			curOff = -e.offset
		} else {
			curOff = 1 - e.offset
			fromIdx = e.offset
			if e.width-1-dw <= xs.Len(e.title)-e.offset {
				toIdx = e.offset + e.width - 2
				textOut = chLeft + xs.Slice(e.title, fromIdx, toIdx) + chRight
			} else {
				textOut = chLeft + xs.Slice(e.title, fromIdx, -1)
			}
		}
	}

	fg, bg := e.textColor, e.backColor
	if fg == ColorDefault {
		fg = tm.GetSysColor(ColorEditText)
	}
	if bg == ColorDefault {
		bg = tm.GetSysColor(ColorEditBack)
	}

	canvas.ClearRect(x, y, w-dw, 1, bg)
	canvas.DrawText(x, y, w-dw, textOut, fg, bg)
	if e.active {
		canvas.SetCursorPos(e, e.cursorPos+curOff, 0)
	}

	if e.listboxVisible {
		e.list.Redraw(canvas)
	}

	if e.editBoxMode == EditBoxCombo {
		e.drawButton(canvas)
	}
}

func (e *EditField) GetEnabled() bool {
	return e.enabled
}

func (e *EditField) SetEnabled(active bool) {
	e.enabled = active
}

func (e *EditField) SetAlign(align Align) {
	e.align = align
}

func (e *EditField) GetAlign() Align {
	return e.align
}

func (e *EditField) SetAnchors(anchor Anchor) {
	e.anchor = anchor
}

func (e *EditField) GetAnchors() Anchor {
	return e.anchor
}

func (e *EditField) SetId(id WinId) {
	e.id = id
}

func (e *EditField) GetActive() bool {
	return e.active
}

func (e *EditField) SetActive(active bool) {
	e.active = active
}

func (e *EditField) GetTabStop() bool {
	return e.tabStop
}

func (e *EditField) SetTabStop(tab bool) {
	e.tabStop = tab
}

func (e *EditField) insertRune(ch rune) {
	if e.readonly {
		return
	}

	if e.maxWidth > 0 && xs.Len(e.title) >= e.maxWidth {
		return
	}

	idx := e.cursorPos

	if idx == 0 {
		e.SetText(string(ch) + e.title)
	} else if idx >= xs.Len(e.title) {
		e.SetText(e.title + string(ch))
	} else {
		e.SetText(xs.Slice(e.title, 0, idx) + string(ch) + xs.Slice(e.title, idx, -1))
	}

	e.cursorPos++

	if e.cursorPos >= e.width {
		if e.offset == 0 {
			e.offset = 2
		} else {
			e.offset++
		}
	}
}

func (e *EditField) backspace() {
	if e.title == "" || e.cursorPos == 0 || e.readonly {
		return
	}

	length := xs.Len(e.title)
	if e.cursorPos >= length {
		e.cursorPos--
		e.SetText(xs.Slice(e.title, 0, length-1))
	} else if e.cursorPos == 1 {
		e.cursorPos = 0
		e.SetText(xs.Slice(e.title, 1, -1))
		e.offset = 0
	} else {
		e.cursorPos--
		e.SetText(xs.Slice(e.title, 0, e.cursorPos) + xs.Slice(e.title, e.cursorPos+1, -1))
	}

	if length-1 < e.width {
		e.offset = 0
	}
}

func (e *EditField) del() {
	length := xs.Len(e.title)

	if e.title == "" || e.cursorPos == length || e.readonly {
		return
	}

	if e.cursorPos == length-1 {
		e.SetText(xs.Slice(e.title, 0, length-1))
	} else {
		e.SetText(xs.Slice(e.title, 0, e.cursorPos) + xs.Slice(e.title, e.cursorPos+1, -1))
	}

	if length-1 < e.width {
		e.offset = 0
	}
}

func (e *EditField) charLeft() {
	if e.cursorPos == 0 || e.title == "" {
		return
	}

	if e.cursorPos == e.offset {
		e.offset--
	}

	e.cursorPos--
}

func (e *EditField) charRight() {
	length := xs.Len(e.title)
	if e.cursorPos == length || e.title == "" {
		return
	}

	e.cursorPos++
	if e.cursorPos != length && e.cursorPos >= e.offset+e.width-2 {
		e.offset++
	}
}

func (e *EditField) home() {
	e.offset = 0
	e.cursorPos = 0
}

func (e *EditField) end() {
	length := xs.Len(e.title)
	e.cursorPos = length

	if length < e.width {
		return
	}

	e.offset = length - (e.width - 2)
}

func (e *EditField) Clear() {
	e.home()
	e.SetText("")
}

func (e *EditField) ProcessEvent(event Event) bool {
	if !e.active || !e.enabled {
		return false
	}

	if event.Type == EventActivate && event.X == 0 {
		termbox.HideCursor()
	}

	if event.Type == EventKey && event.Key != termbox.KeyTab && event.Key != termbox.KeyEnter {
		switch event.Key {
		case termbox.KeySpace:
			e.insertRune(' ')
			return true
		case termbox.KeyBackspace:
			e.backspace()
			return true
		case termbox.KeyDelete:
			e.del()
			return true
		case termbox.KeyArrowLeft:
			e.charLeft()
			return true
		case termbox.KeyHome:
			e.home()
			return true
		case termbox.KeyEnd:
			e.end()
			return true
		case termbox.KeyCtrlR:
			if !e.readonly {
				e.Clear()
			}
			return true
		case termbox.KeyF5:
			if e.listboxVisible {
				e.closeUpList()
			} else {
				e.dropDownList()
			}
		case termbox.KeyArrowRight:
			e.charRight()
			return true
		default:
			if event.Ch != 0 {
				e.insertRune(event.Ch)
				return true
			}
		}
		return false
	}

	if e.editBoxMode == EditBoxCombo && (event.Type == EventMouse || event.Type == EventMouseClick) {
		if event.Y == e.posY && event.X == e.posX+e.width-1 {
			if e.listboxVisible {
				e.closeUpList()
			} else {
				e.dropDownList()
			}
		}
	}

	return false
}

func (e *EditField) SetVisible(visible bool) {
	e.visible = visible
}

func (e *EditField) GetVisible() bool {
	return e.visible
}

//-----------------------------------------------

func (c *EditField) dropDownList() {
	if c.listboxVisible {
		return
	}

	c.listboxVisible = true

	height := 5
	if c.listHeight != -1 {
		height = c.listHeight
	}
	if height > len(c.itemList) && len(c.itemList) > 0 {
		height = len(c.itemList)
	}
	width := c.width
	if c.listWidth != -1 {
		width = c.listWidth
	}

	var props Props
	idL := c.parent.GetNextControlId()
	c.list = NewListBox(c.parent, idL, c.posX, c.posY+1, width, height, props)
	c.parent.AddControl(c.list)
	c.list.SetTabStop(false)

	for _, str := range c.itemList {
		c.list.AddItem(str)
	}

	if c.title != "" {
		id := c.list.FindItem(c.title, false)
		if id != -1 {
			c.list.SelectItem(id)
		}
	}

	c.list.OnSelectItem(func(ev Event) {
		str := c.list.GetSelectedItem()
		c.closeUpList()
		c.SetText(str)
		if c.parent != nil {
			ev := InternalEvent{act: EventRedraw, sender: c.id}
			c.parent.SendEvent(ev)
		}
	})

	c.parent.ActivateControl(c.list)
}

func (c *EditField) closeUpList() {
	if !c.listboxVisible {
		return
	}
	c.listboxVisible = false
	c.parent.RemoveControl(c.list)
	c.list = nil
	c.SetActive(true)
}

func (e *EditField) SetMaxWidth(w int) {
	e.maxWidth = w
	if w > 0 && xs.Len(e.title) > w {
		e.title = xs.Slice(e.title, 0, w)
		e.end()
	}
}

func (e *EditField) GetMaxWidth() int {
	return e.maxWidth
}

func (e *EditField) GetColors() (Color, Color) {
	return e.textColor, e.backColor
}

func (e *EditField) SetTextColor(clr Color) {
	e.textColor = clr
}

func (e *EditField) SetBackColor(clr Color) {
	e.backColor = clr
}

func (e *EditField) HideChildren() {
	if e.editBoxMode == EditBoxCombo {
		e.closeUpList()
	}
}

func (e *EditField) GetScale() int {
	return e.scale
}

func (e *EditField) SetScale(scale int) {
	e.scale = scale
}
