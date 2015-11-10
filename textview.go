package clui

import (
	"bufio"
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	"os"
	"strings"
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
type TextView struct {
	ControlBase
	// own listbox members
	lines   []string
	lengths []int
	// for up/down scroll
	topLine int
	// for side scroll
	leftShift     int
	wordWrap      bool
	colorized     bool
	virtualHeight int
	virtualWidth  int
	multicolor    bool
}

/*
NewListBox creates a new frame.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width and heigth - are minimal size of the control.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func NewTextView(view View, parent Control, width, height int, scale int) *TextView {
	l := new(TextView)

	if height == AutoSize {
		height = 3
	}
	if width == AutoSize {
		width = 5
	}

	l.SetSize(width, height)
	l.SetConstraints(width, height)
	l.topLine = 0
	l.lines = make([]string, 0)
	l.parent = parent
	l.view = view

	l.SetTabStop(true)

	if parent != nil {
		parent.AddChild(l, scale)
	}

	return l
}

func (l *TextView) outputHeight() int {
	h := l.height
	if !l.wordWrap {
		h--
	}
	return h
}

func (l *TextView) redrawScrolls(canvas Canvas, tm Theme) {
	fg, bg := RealColor(tm, l.fg, ColorScrollText), RealColor(tm, l.bg, ColorScrollBack)
	fgThumb, bgThumb := RealColor(tm, l.fg, ColorThumbText), RealColor(tm, l.bg, ColorThumbBack)

	height := l.outputHeight()
	pos := ThumbPosition(l.topLine, l.virtualHeight-l.outputHeight(), height)
	canvas.DrawScroll(l.x+l.width-1, l.y, 1, height, pos, fg, bg, fgThumb, bgThumb, tm.SysObject(ObjScrollBar))

	if !l.wordWrap {
		pos = ThumbPosition(l.leftShift, l.virtualWidth-l.width+1, l.width-1)
		canvas.DrawScroll(l.x, l.y+l.height-1, l.width-1, 1, pos, fg, bg, fgThumb, bgThumb, tm.SysObject(ObjScrollBar))
	}
}

func (l *TextView) redrawText(canvas Canvas, tm Theme) {
	maxWidth := l.width - 1
	maxHeight := l.outputHeight()

	fg, bg := RealColor(tm, l.fg, ColorEditText), RealColor(tm, l.bg, ColorEditBack)
	if l.Active() {
		fg, bg = RealColor(tm, l.fg, ColorEditActiveText), RealColor(tm, l.bg, ColorEditActiveBack)
	}

	if l.wordWrap {
		lineId := l.posToItemNo(l.topLine)
		linePos := l.itemNoToPos(lineId)

		y := 0
		for {
			if y >= maxHeight || lineId >= len(l.lines) {
				break
			}

			remained := l.lengths[lineId]
			start := 0
			for remained > 0 {
				var s string
				if l.multicolor {
					s = SliceColorized(l.lines[lineId], start, start+maxWidth)
				} else {
					if remained <= maxWidth {
						s = xs.Slice(l.lines[lineId], start, -1)
					} else {
						s = xs.Slice(l.lines[lineId], start, start+maxWidth)
					}
				}

				if linePos >= l.topLine {
					if l.multicolor {
						canvas.PutColorizedText(l.x, l.y+y, maxWidth, s, fg, bg, Horizontal)
					} else {
						canvas.PutText(l.x, l.y+y, s, fg, bg)
					}
				}

				remained -= maxWidth
				y++
				linePos++
				start += maxWidth

				if y >= maxHeight {
					break
				}
			}

			lineId++
		}
	} else {
		y := 0
		total := len(l.lines)
		for {
			if y+l.topLine >= total {
				break
			}
			if y >= maxHeight {
				break
			}

			str := l.lines[l.topLine+y]
			lineLength := l.lengths[l.topLine+y]
			if l.multicolor {
				if l.leftShift == 0 {
					if lineLength > maxWidth {
						str = SliceColorized(str, 0, maxWidth)
					}
				} else {
					if l.leftShift+maxWidth >= lineLength {
						str = SliceColorized(str, l.leftShift, -1)
					} else {
						str = SliceColorized(str, l.leftShift, maxWidth+l.leftShift)
					}
				}
				canvas.PutColorizedText(l.x, l.y+y, maxWidth, str, fg, bg, Horizontal)
			} else {
				if l.leftShift == 0 {
					if lineLength > maxWidth {
						str = CutText(str, maxWidth)
					}
				} else {
					if l.leftShift+maxWidth >= lineLength {
						str = xs.Slice(str, l.leftShift, -1)
					} else {
						str = xs.Slice(str, l.leftShift, maxWidth+l.leftShift)
					}
				}
				canvas.PutText(l.x, l.y+y, str, fg, bg)
			}

			y++
		}
	}
}

// Repaint draws the control on its View surface
func (l *TextView) Repaint() {
	canvas := l.view.Canvas()
	tm := l.view.Screen().Theme()

	x, y := l.Pos()
	w, h := l.Size()

	bg := RealColor(tm, l.bg, ColorEditBack)
	if l.Active() {
		bg = RealColor(tm, l.bg, ColorEditActiveBack)
	}
	canvas.FillRect(x, y, w, h, term.Cell{Bg: bg, Ch: ' '})
	l.redrawText(canvas, tm)
	l.redrawScrolls(canvas, tm)
}

func (l *TextView) home() {
	l.topLine = 0
}

func (l *TextView) end() {
	height := l.outputHeight()

	if l.virtualHeight <= height {
		return
	}

	if l.topLine+height >= l.virtualHeight {
		return
	}

	l.topLine = l.virtualHeight - height
}

func (l *TextView) moveUp(dy int) {
	if l.topLine == 0 {
		return
	}

	if l.topLine <= dy {
		l.topLine = 0
	} else {
		l.topLine -= dy
	}
}

func (l *TextView) moveDown(dy int) {
	end := l.topLine + l.outputHeight()

	if end >= l.virtualHeight {
		return
	}

	if l.topLine+dy+l.outputHeight() >= l.virtualHeight {
		l.topLine = l.virtualHeight - l.outputHeight()
	} else {
		l.topLine += dy
	}
}

func (l *TextView) moveLeft() {
	if l.wordWrap || l.leftShift == 0 {
		return
	}

	l.leftShift--
}

func (l *TextView) moveRight() {
	if l.wordWrap {
		return
	}

	if l.leftShift+l.width-1 >= l.virtualWidth {
		return
	}

	l.leftShift++
}

func (l *TextView) processMouseClick(ev Event) bool {
	if ev.Key != term.MouseLeft {
		return false
	}

	dx := ev.X - l.x
	dy := ev.Y - l.y
	yy := l.outputHeight()

	// cursor is not on any scrollbar
	if dx != l.width-1 && dy != l.height-1 {
		return false
	}
	// wordwrap mode does not have horizontal scroll
	if l.wordWrap && dx != l.width-1 {
		return false
	}
	// corner in not wordwrap mode
	if !l.wordWrap && dx == l.width-1 && dy == l.height-1 {
		return false
	}

	// vertical scroll bar
	if dx == l.width-1 {
		if dy == 0 {
			l.moveUp(1)
		} else if dy == yy-1 {
			l.moveDown(1)
		} else {
			newPos := ItemByThumbPosition(dy, l.virtualHeight-yy+1, yy)
			if newPos >= 0 {
				l.topLine = newPos
			}
		}

		return true
	}

	// horizontal scrollbar
	if dx == 0 {
		l.moveLeft()
	} else if dx == l.width-2 {
		l.moveRight()
	} else {
		newPos := ItemByThumbPosition(dx, l.virtualWidth-l.width+2, l.width-1)
		if newPos >= 0 {
			l.leftShift = newPos
		}
	}

	return true
}

/*
ProcessEvent processes all events come from the control parent. If a control
processes an event it should return true. If the method returns false it means
that the control do not want or cannot process the event and the caller sends
the event to the control parent
*/
func (l *TextView) ProcessEvent(event Event) bool {
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
			l.moveUp(1)
			return true
		case term.KeyArrowDown:
			l.moveDown(1)
			return true
		case term.KeyArrowLeft:
			l.moveLeft()
			return true
		case term.KeyArrowRight:
			l.moveRight()
			return true
		case term.KeyPgup:
			l.moveUp(l.outputHeight())
		case term.KeyPgdn:
			l.moveDown(l.outputHeight())
		default:
			return false
		}
	case EventMouse:
		return l.processMouseClick(event)
	}

	return false
}

// own methods

func (l *TextView) calculateVirtualSize() {
	w := l.width - 1
	l.virtualWidth = l.width - 1
	l.virtualHeight = 0

	l.lengths = make([]int, len(l.lines))
	for idx, str := range l.lines {
		if l.multicolor {
			str = UnColorizeText(str)
		}

		sz := xs.Len(str)
		if l.wordWrap {
			n := sz / w
			r := sz % w
			l.virtualHeight += n
			if r != 0 {
				l.virtualHeight++
			}
		} else {
			l.virtualHeight++
			if sz > l.virtualWidth {
				l.virtualWidth = sz
			}
		}
		l.lengths[idx] = sz
	}
}

func (l *TextView) SetText(text []string) {
	l.lines = make([]string, len(text))
	copy(l.lines, text)

	l.calculateVirtualSize()
}

// MultiColored returns if the TextView checks and applies any
// color related tags inside its text. If MultiColores is
// false then text is displayed as is.
// To read about available color tags, please see ColorParser
func (l *TextView) MultiColored() bool {
	return l.multicolor
}

// SetMultiColored changes how the TextView output its text: as is
// or parse and apply all internal color tags
func (l *TextView) SetMultiColored(multi bool) {
	if l.multicolor != multi {
		l.multicolor = multi
		l.calculateVirtualSize()
	}
}

func (l *TextView) posToItemNo(pos int) int {
	id := 0
	for idx, item := range l.lengths {
		if l.virtualWidth >= item {
			pos--
		} else {
			pos -= item / l.virtualWidth
			if item%l.virtualWidth != 0 {
				pos--
			}
		}

		if pos <= 0 {
			id = idx
			break
		}
	}

	return id
}

func (l *TextView) itemNoToPos(id int) int {
	pos := 0
	for i := 0; i < id; i++ {
		if l.virtualWidth >= l.lengths[i] {
			pos++
		} else {
			pos += l.lengths[i] / l.virtualWidth
			if l.lengths[i]%l.virtualWidth != 0 {
				pos++
			}
		}
	}

	return pos
}

func (l *TextView) WordWrap() bool {
	return l.wordWrap
}

func (l *TextView) recalculateTopLine() {
	currLn := l.topLine

	if l.wordWrap {
		l.topLine = l.itemNoToPos(currLn)
	} else {
		l.topLine = l.posToItemNo(currLn)
	}
}

func (l *TextView) SetWordWrap(wrap bool) {
	if wrap != l.wordWrap {
		l.wordWrap = wrap
		l.calculateVirtualSize()
		l.recalculateTopLine()
		l.Repaint()
	}
}

func (l *TextView) LoadFile(filename string) bool {
	l.lines = make([]string, 0)

	file, err := os.Open(filename)
	if err != nil {
		return false
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		l.lines = append(l.lines, line)
	}

	l.calculateVirtualSize()

	return true
}
