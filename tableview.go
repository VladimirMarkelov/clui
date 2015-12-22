package clui

import (
	"fmt"
	term "github.com/nsf/termbox-go"
)

type Column struct {
	Title     string
	Width     int
	Alignment Align
	Fg, Bg    term.Attribute
}

type ColumnDrawInfo struct {
	Row          int
	Col          int
	Width        int
	Text         string
	Alignment    Align
	RowSelected  bool
	CellSelected bool
	Fg           term.Attribute
	Bg           term.Attribute
}

/*
TableView is control to display a list of items and allow to user to select any of them.
Content is scrollable with arrow keys or by clicking up and bottom buttons
on the scroll(now content is scrollable with mouse dragging only on Windows).

TableView calls onSelectItem item function after a user changes currently
selected item with mouse or using keyboard (extra case: the event is emitted
when a user presses Enter - the case is used in ComboBox to select an item
from drop down list). Event structure has 2 fields filled: Y - selected
item number in list(-1 if nothing is selected), Msg - text of the selected item.
*/
type TableView struct {
	ControlBase
	// own listbox members
	topRow        int
	topCol        int
	selectedRow   int
	selectedCol   int
	columns       []Column
	rowCount      int
	fullRowSelect bool
	showRowNo     bool
	showVLines    bool

	onDrawCell   func(*ColumnDrawInfo)
	onAction     func()
	onKeyPress   func(term.Key) bool
	onSelectCell func(int, int)
}

/*
NewTableView creates a new frame.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width and heigth - are minimal size of the control.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func NewTableView(view View, parent Control, width, height int, scale int) *TableView {
	l := new(TableView)

	if height == AutoSize {
		height = 3
	}
	if width == AutoSize {
		width = 10
	}

	l.SetSize(width, height)
	l.SetConstraints(width, height)
	l.selectedCol = 0
	l.selectedRow = 0
	l.parent = parent
	l.view = view
	l.columns = make([]Column, 0)

	l.SetTabStop(true)

	l.onDrawCell = nil
	l.onAction = nil
	l.onKeyPress = nil
	l.onSelectCell = nil

	if parent != nil {
		parent.AddChild(l, scale)
	}

	return l
}

func (l *TableView) redrawHeader(canvas Canvas, tm Theme) {
	fg, bg := RealColor(tm, l.fg, ColorTableHeaderText), RealColor(tm, l.bg, ColorTableHeaderBack)
	fgLine := RealColor(tm, l.fg, ColorTableLineText)
	x, y := l.Pos()
	w, _ := l.Size()
	canvas.FillRect(x, y, w, 1, term.Cell{Fg: fg, Bg: bg, Ch: ' '})
	parts := []rune(tm.SysObject(ObjTableView))

	for i := 0; i < w; i++ {
		canvas.PutSymbol(x+i, y+1, term.Cell{Ch: parts[0], Fg: fg, Bg: bg})
	}
	w-- // scrollbar

	dx := 0
	if l.showVLines {
		dx = 1
	}

	pos := 0
	if l.showRowNo {
		cW := l.counterWidth()
		shift, str := AlignText("#", cW, AlignRight)
		canvas.PutText(x+pos+shift, y, str, fg, bg)
		if l.showVLines {
			canvas.PutSymbol(x+pos+cW, y, term.Cell{Ch: parts[1], Fg: fgLine, Bg: bg})
			canvas.PutSymbol(x+pos+cW, y+1, term.Cell{Ch: parts[2], Fg: fgLine, Bg: bg})
			pos++
		}
		pos = cW + dx
	}

	idx := l.topCol
	for pos < w && idx < len(l.columns) {
		w := l.columns[idx].Width
		if l.width-pos < w {
			w = l.width - pos
		}
		if w <= 0 {
			break
		}

		shift, str := AlignText(l.columns[idx].Title, w, l.columns[idx].Alignment)
		canvas.PutText(x+pos+shift, y, str, fg, bg)
		pos += w

		if l.showVLines && idx < len(l.columns)-1 {
			canvas.PutSymbol(x+pos, y, term.Cell{Ch: parts[1], Fg: fgLine, Bg: bg})
			canvas.PutSymbol(x+pos, y+1, term.Cell{Ch: parts[2], Fg: fgLine, Bg: bg})
			pos++
		}

		idx++
	}
}

func (l *TableView) counterWidth() int {
	width := 0

	if l.showRowNo {
		s := fmt.Sprintf("%v", l.rowCount)
		if s == "" {
			s = " "
		}
		width = len(s)
	}

	return width
}

func (l *TableView) redrawScroll(canvas Canvas, tm Theme) {
	fg, bg := RealColor(tm, l.fg, ColorScrollText), RealColor(tm, l.bg, ColorScrollBack)
	fgThumb, bgThumb := RealColor(tm, l.fg, ColorThumbText), RealColor(tm, l.bg, ColorThumbBack)

	pos := ThumbPosition(l.selectedRow, l.rowCount, l.height-1)
	canvas.DrawScroll(l.x+l.width-1, l.y, 1, l.height-1, pos, fg, bg, fgThumb, bgThumb, tm.SysObject(ObjScrollBar))

	pos = ThumbPosition(l.selectedCol, len(l.columns), l.width-1)
	canvas.DrawScroll(l.x, l.y+l.height-1, l.width-1, 1, pos, fg, bg, fgThumb, bgThumb, tm.SysObject(ObjScrollBar))
	canvas.PutSymbol(l.x+l.width-1, l.y+l.height-1, term.Cell{Ch: ' ', Fg: fg, Bg: bg})
}

func (l *TableView) redrawCells(canvas Canvas, tm Theme) {
	maxRow := l.rowCount - 1
	rowNo := l.topRow
	dy := 2
	maxDy := l.height - 2

	fg, bg := RealColor(tm, l.fg, ColorTableText), RealColor(tm, l.bg, ColorTableBack)
	fgRow, bgRow := RealColor(tm, l.fg, ColorTableSelectedText), RealColor(tm, l.bg, ColorTableSelectedBack)
	fgCell, bgCell := RealColor(tm, l.fg, ColorTableActiveCellText), RealColor(tm, l.bg, ColorTableActiveCellBack)
	fgLine := RealColor(tm, l.fg, ColorTableLineText)
	parts := []rune(tm.SysObject(ObjTableView))

	start := 0
	if l.showRowNo {
		start = l.counterWidth()
		for idx := 1; idx < l.height-2; idx++ {
			if l.topRow+idx > l.rowCount {
				break
			}
			s := fmt.Sprintf("%v", idx+l.topRow)
			shift, str := AlignText(s, start, AlignRight)
			canvas.PutText(l.x+shift, l.y+dy+idx-1, str, fg, bg)
			if l.showVLines {
				canvas.PutSymbol(l.x+start, l.y+dy+idx-1, term.Cell{Ch: parts[1], Fg: fgLine, Bg: bg})
			}
		}
		if l.showVLines {
			start++
		}
	}

	for rowNo <= maxRow && dy <= maxDy {
		colNo := l.topCol
		dx := start
		for colNo < len(l.columns) && dx < l.width-1 {
			c := l.columns[colNo]
			info := ColumnDrawInfo{Row: rowNo, Col: colNo, Width: c.Width, Alignment: c.Alignment}
			if l.selectedRow == rowNo && l.selectedCol == colNo {
				info.RowSelected = true
				info.CellSelected = true
				info.Bg = bgCell
				info.Fg = fgCell
			} else if l.selectedRow == rowNo && l.fullRowSelect {
				info.RowSelected = true
				info.Bg = bgRow
				info.Fg = fgRow
			} else {
				info.Fg = fg
				info.Bg = bg
			}

			if l.onDrawCell != nil {
				l.onDrawCell(&info)
			}

			length := c.Width
			if length+dx >= l.width-1 {
				length = l.width - 1 - dx
			}
			canvas.FillRect(l.x+dx, l.y+dy, length, 1, term.Cell{Bg: info.Bg, Ch: ' ', Fg: info.Fg})
			shift, text := AlignText(info.Text, length, info.Alignment)
			canvas.PutText(l.x+dx+shift, l.y+dy, text, info.Fg, info.Bg)

			dx += c.Width
			if l.showVLines && dx < l.width-1 && colNo < len(l.columns)-1 {
				canvas.PutSymbol(l.x+dx, l.y+dy, term.Cell{Ch: parts[1], Fg: fg, Bg: bg})
				dx++
			}

			colNo++
		}

		rowNo++
		dy++
	}
}

// Repaint draws the control on its View surface
func (l *TableView) Repaint() {
	canvas := l.view.Canvas()
	tm := l.view.Screen().Theme()

	x, y := l.Pos()
	w, h := l.Size()

	bg := RealColor(tm, l.bg, ColorTableBack)
	canvas.FillRect(x, y+2, w, h-2, term.Cell{Bg: bg, Ch: ' '})
	l.redrawHeader(canvas, tm)
	l.redrawScroll(canvas, tm)
	l.redrawCells(canvas, tm)
}

func (l *TableView) home() {
	if len(l.columns) > 0 {
		l.selectedCol = 0
	}
	l.topCol = 0
	l.EnsureColVisible()
}

func (l *TableView) end() {
	length := len(l.columns)

	if length == 0 {
		return
	}

	l.selectedCol = length - 1
	l.EnsureColVisible()
}

func (l *TableView) firstRow() {
	if l.rowCount > 0 {
		l.selectedRow = 0
	}
	l.topRow = 0
	l.EnsureRowVisible()
}

func (l *TableView) lastRow() {
	if l.rowCount == 0 {
		return
	}

	l.selectedRow = l.rowCount - 1
	l.EnsureColVisible()
}

func (l *TableView) moveUp(dy int) {
	if l.topRow == 0 && l.selectedRow == 0 {
		return
	}

	if l.selectedRow == -1 {
		if l.rowCount != 0 {
			l.selectedRow = 0
		}
		return
	}

	if l.selectedRow < dy {
		l.selectedRow = 0
	} else {
		l.selectedRow -= dy
	}

	l.EnsureRowVisible()
}

func (l *TableView) moveDown(dy int) {
	length := l.rowCount

	if length == 0 || l.selectedRow == length-1 {
		return
	}

	if l.selectedRow+dy >= length {
		l.selectedRow = length - 1
	} else {
		l.selectedRow += dy
	}

	l.EnsureRowVisible()
}

func (l *TableView) moveRight(dx int) {
	colCnt := len(l.columns)
	if l.selectedCol == colCnt-1 || colCnt == 0 {
		return
	}

	if l.selectedCol == -1 {
		l.selectedCol = 0
	} else {
		if l.selectedCol+dx >= colCnt {
			l.selectedCol = colCnt - 1
		} else {
			l.selectedCol += dx
		}
	}

	l.EnsureColVisible()
}

func (l *TableView) moveLeft(dx int) {
	colCnt := len(l.columns)
	if l.selectedCol == 0 || colCnt == 0 {
		return
	}

	if l.selectedCol == -1 {
		l.selectedCol = 0
	} else {
		if l.selectedCol-dx < 0 {
			l.selectedCol = 0
		} else {
			l.selectedCol -= dx
		}
	}

	l.EnsureColVisible()
}

func (l *TableView) isColVisible(idx int) bool {
	if idx < l.topCol {
		return false
	}

	width := l.width - 1
	width -= l.counterWidth()
	if l.showVLines && l.showRowNo {
		width--
	}

	for i := l.topCol; i < len(l.columns); i++ {
		if i == idx && l.columns[i].Width <= width {
			return true
		}

		width -= l.columns[i].Width
	}

	return false
}

func (l *TableView) EnsureColVisible() {
	if l.isColVisible(l.selectedCol) {
		return
	}

	if l.selectedCol < l.topCol {
		l.topCol = l.selectedCol
		return
	}

	width := l.width - 1 - l.counterWidth()
	if l.showRowNo && l.showVLines {
		width--
	}

	toShow := l.selectedCol
	for width > 0 {
		if l.columns[toShow].Width > width {
			if toShow == l.selectedCol {
				break
			} else {
				toShow++
				break
			}
		} else if l.columns[toShow].Width == width {
			break
		} else {
			width -= l.columns[toShow].Width
			if width < 0 {
				break
			}
			toShow--
			if toShow == 0 {
				break
			}
		}
	}

	l.topCol = toShow
}

// EnsureRowVisible makes the currently selected row is visible and scrolls the item list if it is required
func (l *TableView) EnsureRowVisible() {
	length := l.rowCount

	hgt := l.height - 3

	if length <= hgt || l.selectedRow == -1 {
		return
	}

	diff := l.selectedRow - l.topRow
	if diff >= 0 && diff < hgt {
		return
	}

	if diff < 0 {
		l.topRow = l.selectedRow
	} else {
		top := l.selectedRow - hgt + 1
		if length-top > hgt {
			l.topRow = top
		} else {
			l.topRow = length - hgt
		}
	}
}

// Clear deletes all TableView items
func (l *TableView) Clear() {
	l.selectedRow = 0
	l.selectedCol = 0
	l.topRow = 0
	l.topCol = 0
}

func (l *TableView) mouseToCol(dx int) int {
	shift := l.counterWidth()
	if l.showVLines {
		shift++
	}

	if dx < shift {
		return l.selectedCol
	}

	idx := l.topCol
	selectedCol := l.selectedCol
	for {
		if shift+l.columns[idx].Width > dx {
			selectedCol = idx
			break
		}

		if idx == len(l.columns)-1 {
			selectedCol = idx
			break
		}

		shift += l.columns[idx].Width
		if l.showVLines {
			shift++
		}
		idx++
	}

	return selectedCol
}

func (l *TableView) horizontalScrollClick(dx int) {
	if dx == 0 {
		l.moveLeft(1)
		return
	} else if dx == l.width-2 {
		l.moveRight(1)
	} else if dx > 0 && dx < l.width-2 {
		pos := ThumbPosition(l.selectedCol, len(l.columns), l.width-1)
		if pos < dx {
			l.moveRight(1)
		} else if pos > dx {
			l.moveLeft(1)
		}
	}
}

func (l *TableView) verticalScrollClick(dy int) {
	if dy == 0 {
		l.moveUp(1)
		return
	} else if dy == l.height-2 {
		l.moveDown(1)
	} else if dy > 0 && dy < l.height-2 {
		pos := ThumbPosition(l.selectedRow, l.rowCount, l.height-1)
		l.Logger().Printf("POS: %v, DY: %v", pos, dy)
		if pos > dy {
			l.moveUp(l.height - 3)
		} else if pos < dy {
			l.moveDown(l.height - 3)
		}
	}
}

func (l *TableView) processMouseClick(ev Event) bool {
	if ev.Key != term.MouseLeft {
		return false
	}

	dx := ev.X - l.x
	dy := ev.Y - l.y

	if dy == l.height-1 && dx == l.width-1 {
		l.selectedRow = l.rowCount - 1
		l.selectedCol = len(l.columns) - 1
		return true
	}

	if dy == l.height-1 {
		l.horizontalScrollClick(dx)
		return true
	}

	if dx == l.width-1 {
		l.verticalScrollClick(dy)
		return true
	}

	if dy < 2 {
		// Header - no action now
		return true
	}

	dy -= 2
	l.selectedRow = l.topRow + dy

	oldCol := l.selectedCol
	l.selectedCol = l.mouseToCol(dx)

	if oldCol != l.selectedCol {
		l.EnsureColVisible()
	}

	return true
}

/*
ProcessEvent processes all events come from the control parent. If a control
processes an event it should return true. If the method returns false it means
that the control do not want or cannot process the event and the caller sends
the event to the control parent
*/
func (l *TableView) ProcessEvent(event Event) bool {
	if !l.Active() || !l.Enabled() {
		return false
	}

	switch event.Type {
	case EventKey:
		if l.onKeyPress != nil {
			res := l.onKeyPress(event.Key)
			if res {
				return true
			}
		}

		switch event.Key {
		case term.KeyHome:
			if event.Mod == term.ModAlt {
				l.selectedRow = 0
				l.EnsureRowVisible()
			} else {
				l.home()
			}
			return true
		case term.KeyEnd:
			if event.Mod == term.ModAlt {
				l.selectedRow = l.rowCount - 1
				l.EnsureRowVisible()
			} else {
				l.end()
			}
			return true
		case term.KeyArrowUp:
			l.moveUp(1)
			return true
		case term.KeyArrowDown:
			l.moveDown(1)
			return true
		case term.KeyArrowLeft:
			l.moveLeft(1)
			return true
		case term.KeyArrowRight:
			l.moveRight(1)
			return true
		case term.KeyPgdn:
			l.moveDown(l.height - 3)
			return true
		case term.KeyPgup:
			l.moveUp(l.height - 3)
			return true
			// 	case term.KeyCtrlM:
			// 		if l.currSelection != -1 && l.onSelectItem != nil {
			// 			ev := Event{Y: l.currSelection, Msg: l.SelectedItemText()}
			// 			go l.onSelectItem(ev)
			// 		}
			// 	default:
			// 		return false
		}
	case EventMouse:
		return l.processMouseClick(event)
	}

	return false
}

// own methods

func (l *TableView) ShowLines() bool {
	return l.showVLines
}

func (l *TableView) SetShowLines(show bool) {
	l.showVLines = show
}

func (l *TableView) ShowRowNumber() bool {
	return l.showRowNo
}

func (l *TableView) SetShowRowNumber(show bool) {
	l.showRowNo = show
}

func (l *TableView) Columns() []Column {
	c := make([]Column, len(l.columns))
	copy(c, l.columns)
	return c
}

func (l *TableView) SetColumns(cols []Column) {
	l.columns = cols
}

func (l *TableView) SetColumnInfo(id int, col Column) {
	if id < len(l.columns) {
		l.columns[id] = col
	}
}

func (l *TableView) RowCount() int {
	return l.rowCount
}

func (l *TableView) SetRowCount(count int) {
	l.rowCount = count
}

func (l *TableView) FullRowSelect() bool {
	return l.fullRowSelect
}

func (l *TableView) SetFullRowSelect(fullRow bool) {
	l.fullRowSelect = fullRow
}

// OnSelectItem sets a callback that is called every time
// the selected item is changed
func (l *TableView) OnSelectCell(fn func(int, int)) {
	l.onSelectCell = fn
}

// OnKeyPress sets the callback that is called when a user presses a Key while
// the controls is active. If a handler processes the key it should return
// true. If handler returns false it means that the default handler will
// process the key
func (l *TableView) OnKeyPress(fn func(term.Key) bool) {
	l.onKeyPress = fn
}

func (l *TableView) OnDrawCell(fn func(*ColumnDrawInfo)) {
	l.onDrawCell = fn
}
