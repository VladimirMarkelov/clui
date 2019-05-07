package clui

import (
	"fmt"
	term "github.com/nsf/termbox-go"
	мКнст "./пакКонстанты"
	мИнт "./пакИнтерфейсы"
)

/*
TableView is control to display a list of items in a table(grid).
Content is scrollable with arrow keys and mouse.
TableView always works in virtual mode - it does not keep table
data and always asks for the cell value using callback OnDrawCell.

Predefined hotkeys:
  Arrows - move cursor
  Home, End - move cursor to first and last column, respectively
  Alt+Home, Alt+End - move cursor to first and last row, respectively
  PgDn, PgUp - move cursor to a screen down and up
  Enter, F2 - emits event TableActionEdit
  Insert - emits event TableActionNew
  Delete - emits event TableActionDelete
  F4 - Change sort mode

Events:
  OnDrawCell - called every time the table is going to draw a cell.
        The argument is ColumnDrawInfo prefilled with the current
        cell attributes. Callback should fill at least the field
        Title. Filling Bg, Fg, and Alignment are optional. Changing
        other fields in callback does not make any difference - they
        are only for caller convenience
  OnAction - called when a user pressed some hotkey(please, see
        above) or clicks any column header(in this case, the control
        sends TableActionSort event and fills column number and
        sorting type - no sort, ascending, descending)
  OnKeyPress - called every time a user presses a key. Callback should
        return true if TableView must skip internal key processing.
        E.g, a user can disable emitting TableActionDelete event by
        adding callback OnKeyPress and return true in case of Delete
        key is pressed
  OnSelectCell - called in case of the currently selected row or
        column is changed
  OnBeforeDraw - called right before the TableView is going to repaint
        itself. It can be used to prepare all the data beforehand and
        then quickly use cached data inside OnDrawCell. Callback
        receives 4 arguments: first visible column, first visible row,
        number of visible columns, number of visible rows.
*/
type TableView struct {
	*BaseControl
	// own TableView members
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
	onAction     func(TableEvent)
	onKeyPress   func(term.Key) bool
	onSelectCell func(int, int)
	onBeforeDraw func(int, int, int, int)

	// internal variable to avoid sending onSelectCell twice or more
	// in case of current cell is unchanged
	lastEventCol int
	lastEventRow int
}

// Column is a information about a table column.
// When one sets a column list, it the fields Title
// and Width should be set. All other fields can be
// undefined.
type Column struct {
	Title     string
	Width     int
	Alignment мИнт.Align
	Fg, Bg    term.Attribute
	Sort      мКнст.SortOrder
}

// ColumnDrawInfo is a structure used in OnDrawCell event.
// A callback should assign Text field otherwise the cell
// will be empty. In addition to it, the callback can
// change Bg, Fg, and Alignment to display customizes
// info. All other non-mentioned fields are for a user
// convenience and used to describe the cell more detailed,
// changing that fields affects nothing
type ColumnDrawInfo struct {
	// row number
	Row int
	// column number
	Col int
	// width of the cell
	Width int
	// cell displayed text
	Text string
	// text alignment
	Alignment мИнт.Align
	// is the row that contains the cell selected(active)
	RowSelected bool
	// is the column that contains the cell selected(active)
	CellSelected bool
	// current text color
	Fg term.Attribute
	// current background color
	Bg term.Attribute
}

// TableEvent is structure to describe the common action that a
// TableView ask for while a user is interacting with the table
type TableEvent struct {
	// requested action: Add, Edit, Delete, Sort data
	Action мКнст.TableAction
	// Currently selected column
	Col int
	// Currently selected row (it is not used for TableActionSort)
	Row int
	// Sort order (it is used only in TableActionSort event)
	Sort мКнст.SortOrder
}

/*
CreateTableView creates a new frame.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width and height - are minimal size of the control.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func CreateTableView(parent мИнт.ИВиджет, width, height int, scale int) *TableView {
	l := new(TableView)
	l.BaseControl = NewBaseControl()

	if height == мКнст.AutoSize {
		height = 3
	}
	if width == мКнст.AutoSize {
		width = 10
	}

	l.SetSize(width, height)
	l.SetConstraints(width, height)
	l.selectedCol = 0
	l.selectedRow = 0
	l.parent = parent
	l.columns = make([]Column, 0)
	l.SetScale(scale)

	l.SetTabStop(true)

	l.onDrawCell = nil
	l.onAction = nil
	l.onKeyPress = nil
	l.onSelectCell = nil
	l.lastEventCol = -1
	l.lastEventRow = -1

	if parent != nil {
		parent.AddChild(l)
	}

	return l
}

func (l *TableView) drawHeader() {
	PushAttributes()
	defer PopAttributes()

	fg, bg := RealColor(l.fg, l.Style(), мКнст.ColorTableHeaderText), RealColor(l.bg, l.Style(), мКнст.ColorTableHeaderBack)
	fgLine := RealColor(l.fg, l.Style(), мКнст.ColorTableLineText)
	x, y := l.Pos()
	w, _ := l.Size()
	SetTextColor(fg)
	SetBackColor(bg)
	FillRect(x, y, w, 1, ' ')
	parts := []rune(SysObject(мКнст.ObjTableView))

	for i := 0; i < w; i++ {
		PutChar(x+i, y+1, parts[0])
	}
	w-- // scrollbar

	dx := 0
	if l.showVLines {
		dx = 1
	}

	pos := 0
	SetBackColor(bg)
	if l.showRowNo {
		cW := l.counterWidth()
		shift, str := AlignText("#", cW, мКнст.AlignRight)
		SetTextColor(fg)
		DrawRawText(x+pos+shift, y, str)
		if l.showVLines {
			SetTextColor(fgLine)
			PutChar(x+pos+cW, y, parts[1])
			PutChar(x+pos+cW, y+1, parts[2])
			pos++
		}
		pos = cW + dx
	}

	idx := l.topCol
	for pos < w && idx < len(l.columns) {
		w := l.columns[idx].Width
		if l.width-1-pos < w {
			w = l.width - 1 - pos
		}
		if w <= 0 {
			break
		}

		dw := 0
		if l.columns[idx].Sort != мКнст.SortNone {
			dw = -1
			ch := parts[3]
			if l.columns[idx].Sort == мКнст.SortDesc {
				ch = parts[4]
			}
			SetTextColor(fg)
			PutChar(x+pos+w-1, y, ch)
		}

		shift, str := AlignColorizedText(l.columns[idx].Title, w+dw, l.columns[idx].Alignment)
		SetTextColor(fg)
		DrawText(x+pos+shift, y, str)
		pos += w

		if l.showVLines && idx < len(l.columns)-1 {
			SetTextColor(fgLine)
			PutChar(x+pos, y, parts[1])
			PutChar(x+pos, y+1, parts[2])
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

func (l *TableView) drawScroll() {

	pos := ThumbPosition(l.selectedRow, l.rowCount, l.height-1)
	DrawScrollBar(l.x+l.width-1, l.y, 1, l.height-1, pos)

	pos = ThumbPosition(l.selectedCol, len(l.columns), l.width-1)
	DrawScrollBar(l.x, l.y+l.height-1, l.width-1, 1, pos)
	PutChar(l.x+l.width-1, l.y+l.height-1, ' ')
}

func (l *TableView) drawCells() {
	PushAttributes()
	defer PopAttributes()

	maxRow := l.rowCount - 1
	rowNo := l.topRow
	dy := 2
	maxDy := l.height - 2

	fg, bg := RealColor(l.fg, l.Style(), мКнст.ColorTableText), RealColor(l.bg, l.Style(), мКнст.ColorTableBack)
	fgRow, bgRow := RealColor(l.fg, l.Style(), мКнст.ColorTableSelectedText), RealColor(l.bg, l.Style(), мКнст.ColorTableSelectedBack)
	fgCell, bgCell := RealColor(l.fg, l.Style(), мКнст.ColorTableActiveCellText), RealColor(l.bg, l.Style(), мКнст.ColorTableActiveCellBack)
	fgLine := RealColor(l.fg, l.Style(), мКнст.ColorTableLineText)
	parts := []rune(SysObject(мКнст.ObjTableView))

	start := 0
	if l.showRowNo {
		start = l.counterWidth()
		for idx := 1; idx < l.height-2; idx++ {
			if l.topRow+idx > l.rowCount {
				break
			}
			s := fmt.Sprintf("%v", idx+l.topRow)
			shift, str := AlignText(s, start, мКнст.AlignRight)
			SetTextColor(fg)
			SetBackColor(bg)
			DrawText(l.x+shift, l.y+dy+idx-1, str)
			if l.showVLines {
				SetTextColor(fgLine)
				PutChar(l.x+start, l.y+dy+idx-1, parts[1])
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
			SetTextColor(info.Fg)
			SetBackColor(info.Bg)
			FillRect(l.x+dx, l.y+dy, length, 1, ' ')
			shift, text := AlignColorizedText(info.Text, length, info.Alignment)
			DrawText(l.x+dx+shift, l.y+dy, text)

			dx += c.Width
			if l.showVLines && dx < l.width-1 && colNo < len(l.columns)-1 {
				SetTextColor(fg)
				SetBackColor(bg)
				PutChar(l.x+dx, l.y+dy, parts[1])
				dx++
			}

			colNo++
		}

		rowNo++
		dy++
	}
}

//Draw Repaint draws the control on its View surface
func (l *TableView) Draw() {
	if l.hidden {
		return
	}

	l.mtx.RLock()
	defer l.mtx.RUnlock()
	PushAttributes()
	defer PopAttributes()

	x, y := l.Pos()
	w, h := l.Size()

	if l.onBeforeDraw != nil {
		firstCol, firstRow, colCount, rowCount := l.VisibleArea()
		l.onBeforeDraw(firstCol, firstRow, colCount, rowCount)
	}

	bg := RealColor(l.bg, l.Style(), мКнст.ColorTableBack)
	SetBackColor(bg)
	FillRect(x, y+2, w, h-2, ' ')
	l.drawHeader()
	l.drawScroll()
	l.drawCells()
}

func (l *TableView) emitSelectionChange() {
	if l.lastEventRow == l.selectedRow && l.lastEventCol == l.selectedCol {
		return
	}

	if l.selectedCol != -1 && l.selectedRow != -1 && l.onSelectCell != nil {
		l.onSelectCell(l.selectedCol, l.selectedRow)
		l.lastEventRow = l.selectedRow
		l.lastEventCol = l.selectedCol
	}
}

func (l *TableView) home() {
	if len(l.columns) > 0 {
		l.selectedCol = 0
	}
	l.topCol = 0
	l.EnsureColVisible()
	l.emitSelectionChange()
}

func (l *TableView) end() {
	length := len(l.columns)

	if length == 0 {
		return
	}

	l.selectedCol = length - 1
	l.EnsureColVisible()
	l.emitSelectionChange()
}

func (l *TableView) firstRow() {
	if l.rowCount > 0 {
		l.selectedRow = 0
	}
	l.topRow = 0
	l.EnsureRowVisible()
	l.emitSelectionChange()
}

func (l *TableView) lastRow() {
	if l.rowCount == 0 {
		return
	}

	l.selectedRow = l.rowCount - 1
	l.EnsureColVisible()
	l.emitSelectionChange()
}

func (l *TableView) moveUp(dy int) {
	if l.topRow == 0 && l.selectedRow == 0 {
		return
	}

	if l.selectedRow == -1 {
		if l.rowCount != 0 {
			l.selectedRow = 0
			l.emitSelectionChange()
		}
		return
	}

	if l.selectedRow < dy {
		l.selectedRow = 0
	} else {
		l.selectedRow -= dy
	}

	l.EnsureRowVisible()
	l.emitSelectionChange()
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
	l.emitSelectionChange()
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
	l.emitSelectionChange()
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
	l.emitSelectionChange()
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

// EnsureColVisible scrolls the table horizontally
// to make the currently selected column fully visible
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

// EnsureRowVisible scrolls the table vertically
// to make the currently selected row visible
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

func (l *TableView) mouseToCol(dx int) int {
	shift := l.counterWidth()
	if l.showVLines {
		shift++
	}

	if dx < shift {
		return -1
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
		if pos > dy {
			l.moveUp(l.height - 3)
		} else if pos < dy {
			l.moveDown(l.height - 3)
		}
	}
}

func (l *TableView) processMouseClick(ev мИнт.ИСобытие) bool {
	if ev.Key() != term.MouseLeft {
		return false
	}

	dx := ev.X() - l.x
	dy := ev.Y() - l.y

	if l.topRow+dy-2 >= l.rowCount && dy != l.height-1 && dx != l.width-1 {
		return false
	}

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
		l.headerClicked(dx)
		return true
	}

	dy -= 2
	newRow := l.topRow + dy

	newCol := l.mouseToCol(dx)
	if newCol == -1 && newRow != l.selectedRow {
		l.selectedRow = newRow
		l.EnsureColVisible()
		l.emitSelectionChange()
	} else if newCol != -1 && (newCol != l.selectedCol || newRow != l.selectedRow) {
		l.selectedCol = newCol
		l.selectedRow = newRow
		l.EnsureColVisible()
		l.emitSelectionChange()
	}

	return true
}

func (l *TableView) headerClicked(dx int) {
	colID := l.mouseToCol(dx)
	if colID == -1 {
		if l.onAction != nil {
			ev := TableEvent{Action: мКнст.TableActionSort, Col: -1, Row: -1}
			l.onAction(ev)
		}
	} else {
		sort := l.columns[colID].Sort

		for idx := range l.columns {
			l.columns[idx].Sort = мКнст.SortNone
		}

		if sort == мКнст.SortAsc {
			sort = мКнст.SortDesc
		} else if sort == мКнст.SortNone {
			sort = мКнст.SortAsc
		} else {
			sort = мКнст.SortNone
		}
		l.columns[colID].Sort = sort

		if l.onAction != nil {
			ev := TableEvent{Action: мКнст.TableActionSort, Col: colID, Row: -1, Sort: sort}
			l.onAction(ev)
		}
	}
}

/*
ProcessEvent processes all events come from the control parent. If a control
processes an event it should return true. If the method returns false it means
that the control do not want or cannot process the event and the caller sends
the event to the control parent
*/
func (l *TableView) ProcessEvent(event мИнт.ИСобытие) bool {
	if !l.Active() || !l.Enabled() {
		return false
	}

	switch event.Type() {
	case мИнт.EventKey:
		if l.onKeyPress != nil {
			res := l.onKeyPress(event.Key())
			if res {
				return true
			}
		}

		switch event.Key() {
		case term.KeyHome:
			if event.Mod() == term.ModAlt {
				l.selectedRow = 0
				l.EnsureRowVisible()
				l.emitSelectionChange()
			} else {
				l.home()
			}
			return true
		case term.KeyEnd:
			if event.Mod() == term.ModAlt {
				l.selectedRow = l.rowCount - 1
				l.EnsureRowVisible()
				l.emitSelectionChange()
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
		case term.KeyCtrlM, term.KeyF2:
			if l.selectedRow != -1 && l.selectedCol != -1 && l.onAction != nil {
				ev := TableEvent{Action: мКнст.TableActionEdit, Col: l.selectedCol, Row: l.selectedRow}
				l.onAction(ev)
			}
		case term.KeyDelete:
			if l.selectedRow != -1 && l.onAction != nil {
				ev := TableEvent{Action: мКнст.TableActionDelete, Col: l.selectedCol, Row: l.selectedRow}
				l.onAction(ev)
			}
		case term.KeyInsert:
			if l.onAction != nil {
				ev := TableEvent{Action: мКнст.TableActionNew, Col: l.selectedCol, Row: l.selectedRow}
				l.onAction(ev)
			}
		case term.KeyF4:
			if l.onAction != nil {
				colID := l.selectedCol
				sort := l.columns[colID].Sort

				for idx := range l.columns {
					l.columns[idx].Sort = мКнст.SortNone
				}

				if sort == мКнст.SortAsc {
					sort = мКнст.SortDesc
				} else if sort == мКнст.SortNone {
					sort = мКнст.SortAsc
				} else {
					sort = мКнст.SortNone
				}
				l.columns[colID].Sort = sort

				ev := TableEvent{Action: мКнст.TableActionSort, Col: colID, Row: -1, Sort: sort}
				l.onAction(ev)
			}
		default:
			return false
		}
	case мИнт.EventMouse:
		return l.processMouseClick(event)
	}

	return false
}

// own methods

// ShowLines returns true if table displays vertical
// lines to separate columns
func (l *TableView) ShowLines() bool {
	return l.showVLines
}

// SetShowLines disables and enables displaying vertical
// lines inside TableView
func (l *TableView) SetShowLines(show bool) {
	l.showVLines = show
}

// ShowRowNumber returns true if the table shows the
// row number as the first table column. This virtual
// column is always fixed and a user cannot change
// displayed text
func (l *TableView) ShowRowNumber() bool {
	return l.showRowNo
}

// SetShowRowNumber turns on and off the first fixed
// column of the table that displays the row number
func (l *TableView) SetShowRowNumber(show bool) {
	l.showRowNo = show
}

// Columns returns the current list of table columns
func (l *TableView) Columns() []Column {
	c := make([]Column, len(l.columns))
	copy(c, l.columns)
	return c
}

// SetColumns replaces existing table column list with
// a new one. Be sure that every item has correct
// Title and Width, all other column properties may
// be undefined
func (l *TableView) SetColumns(cols []Column) {
	l.columns = cols
}

// SetColumnInfo replaces the existing column info
func (l *TableView) SetColumnInfo(id int, col Column) {
	if id < len(l.columns) {
		l.columns[id] = col
	}
}

// RowCount returns current row count
func (l *TableView) RowCount() int {
	return l.rowCount
}

// SetRowCount sets the new row count
func (l *TableView) SetRowCount(count int) {
	l.rowCount = count
}

// FullRowSelect returns if TableView hilites the selected
// cell only or the whole row that contains the selected
// cell. By default the colors for selected row and cell
// are different
func (l *TableView) FullRowSelect() bool {
	return l.fullRowSelect
}

// SetFullRowSelect enables or disables hiliting of the
// full row that contains the selected cell
func (l *TableView) SetFullRowSelect(fullRow bool) {
	l.fullRowSelect = fullRow
}

// OnSelectCell sets a callback that is called every time
// the selected cell is changed
func (l *TableView) OnSelectCell(fn func(int, int)) {
	l.onSelectCell = fn
}

// OnKeyPress sets the callback that is called when a user presses a Key while
// the controls is active. If a handler processes the key it should return
// true. If handler returns false it means that the default handler has to
// process the key
func (l *TableView) OnKeyPress(fn func(term.Key) bool) {
	l.onKeyPress = fn
}

// OnDrawCell is called every time the table is going to display
// a cell
func (l *TableView) OnDrawCell(fn func(*ColumnDrawInfo)) {
	l.mtx.Lock()
	l.onDrawCell = fn
	l.mtx.Unlock()
}

// OnAction is called when the table wants a user application to
// do some job like add, delete, edit or sort data
func (l *TableView) OnAction(fn func(TableEvent)) {
	l.onAction = fn
}

// SelectedRow returns currently selected row number or
// -1 if no row is selected
func (l *TableView) SelectedRow() int {
	return l.selectedRow
}

// SelectedCol returns currently selected column number or
// -1 if no column is selected
func (l *TableView) SelectedCol() int {
	return l.selectedCol
}

// SetSelectedRow changes the currently selected row.
// If row is greater than number of row the last row
// is selected. Set row to -1 to turn off selection.
// The table scrolls automatically to display the column
func (l *TableView) SetSelectedRow(row int) {
	oldSelection := l.selectedRow
	if row >= l.rowCount {
		l.selectedRow = l.rowCount - 1
	} else if row < -1 {
		l.selectedRow = -1
	}

	if l.selectedRow != oldSelection {
		l.EnsureRowVisible()
		l.emitSelectionChange()
	}
}

// SetSelectedCol changes the currently selected column.
// If column is greater than number of columns the last
// column is selected. Set row to -1 to turn off selection.
// The table scrolls automatically to display the column
func (l *TableView) SetSelectedCol(col int) {
	oldSelection := l.selectedCol
	if col >= len(l.columns) {
		l.selectedCol = len(l.columns) - 1
	} else if col < -1 {
		l.selectedCol = -1
	}

	if l.selectedCol != oldSelection {
		l.EnsureColVisible()
		l.emitSelectionChange()
	}
}

// OnBeforeDraw is called when TableView is going to draw its cells.
// Can be used to precache the data, and make OnDrawCell faster.
// Callback receives 4 arguments: first visible column, first visible row,
// the number of visible columns, the number of visible rows
func (l *TableView) OnBeforeDraw(fn func(int, int, int, int)) {
	l.mtx.Lock()
	l.onBeforeDraw = fn
	l.mtx.Unlock()
}

// VisibleArea returns which rows and columns are currently visible. It can be
// used instead of OnBeforeDraw event to prepare the data for drawing without
// waiting until TableView starts drawing itself.
// It can be useful in case of you update your database, so at the same moment
// you can request the visible area and update database cache - it can improve
// performance.
// Returns:
// * firstCol - first visible column
// * firstRow - first visible row
// * colCount - the number of visible columns
// * rowCount - the number of visible rows
func (l *TableView) VisibleArea() (firstCol, firstRow, colCount, rowCount int) {
	firstRow = l.topRow
	maxDy := l.height - 3
	if firstRow+maxDy < l.rowCount {
		rowCount = maxDy
	} else {
		rowCount = l.rowCount - l.topRow
	}

	total := l.width - 1
	if l.showRowNo {
		total -= l.counterWidth()
		if l.showVLines {
			total--
		}
	}

	colNo := l.topCol
	colCount = 0
	for colNo < len(l.columns) && total > 0 {
		w := l.columns[colNo].Width
		total -= w
		colNo++
		colCount++
	}

	return l.topCol, l.topRow, colCount, rowCount
}
