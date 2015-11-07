package clui

import (
	"fmt"
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
)

/*
FrameBuffer represents an object visible content. Used by Composer to
keep all screen info and by View. All methods of FrameBuffer use relative
corrdinate system that starts at left top FrameBuffer corner as 0,0
*/
type FrameBuffer struct {
	buffer [][]term.Cell
	w, h   int
}

// NewFrameBuffer creates new buffer. Width and height of a new buffer cannot be less than 3
func NewFrameBuffer(w, h int) *FrameBuffer {
	if w < 3 || h < 3 {
		panic(fmt.Sprintf("Invalid size: %vx%v.", w, h))
	}

	c := new(FrameBuffer)
	c.SetSize(w, h)

	return c
}

// Size returns current size
func (fb *FrameBuffer) Size() (width int, height int) {
	return fb.w, fb.h
}

/*
SetSize sets the new FrameBuffer size. If new size does not equal old size then
FrameBuffer is recreated and cleared with default colors. Both FrameBuffer width and
height must be greater than 2
*/
func (fb *FrameBuffer) SetSize(w, h int) {
	if w == fb.w && h == fb.h {
		return
	}

	if w < 3 || h < 3 {
		panic(fmt.Sprintf("Invalid size: %vx%v.", w, h))
	}

	fb.w, fb.h = w, h

	fb.buffer = make([][]term.Cell, h)
	for i := 0; i < h; i++ {
		fb.buffer[i] = make([]term.Cell, w)
	}
}

// Clear fills FrameBuffer with given background color
func (fb *FrameBuffer) Clear(bg term.Attribute) {
	s := term.Cell{Ch: ' ', Fg: ColorWhite, Bg: bg}
	for y := 0; y < fb.h; y++ {
		for x := 0; x < fb.w; x++ {
			fb.buffer[y][x] = s
		}
	}
}

// FillRect fills area of FrameBuffer with user-defined rune and colors
func (fb *FrameBuffer) FillRect(x, y, w, h int, s term.Cell) {
	if x < 0 {
		w += x
		x = 0
	}
	if y < 0 {
		h += y
		y = 0
	}
	if x+w >= fb.w {
		w = fb.w - x
	}
	if y+h >= fb.h {
		h = fb.h - y
	}

	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			fb.buffer[yy][xx] = s
		}
	}
}

// Symbol returns current FrameBuffer cell value at given coordinates.
// If coordinates are outside FrameBuffer ok is false
func (fb *FrameBuffer) Symbol(x, y int) (term.Cell, bool) {
	if x >= fb.w || x < 0 || y >= fb.h || y < 0 {
		return term.Cell{Ch: ' '}, false
	}

	return fb.buffer[y][x], true
}

// PutSymbol sets value for the FrameBuffer cell: rune and its colors. Returns result of operation: e.g, if the symbol position is outside FrameBuffer the operation fails and the function returns false
func (fb *FrameBuffer) PutSymbol(x, y int, s term.Cell) bool {
	if x < 0 || x >= fb.w || y < 0 || y >= fb.h {
		return false
	}

	fb.buffer[y][x] = s
	return true
}

// PutChar sets value for the FrameBuffer cell: rune and its colors. Returns result of operation: e.g, if the symbol position is outside FrameBuffer the operation fails and the function returns false
func (fb *FrameBuffer) PutChar(x, y int, c rune, fg, bg term.Attribute) bool {
	if x < 0 || x >= fb.w || y < 0 || y >= fb.h {
		return false
	}

	fb.buffer[y][x] = term.Cell{Ch: c, Fg: fg, Bg: bg}
	return true
}

// PutText draws horizontal string on FrameBuffer clipping by FrameBuffer boundaries. x and y are starting point, text is a string to display, fg and bg are text and background attributes
func (fb *FrameBuffer) PutText(x, y int, text string, fg, bg term.Attribute) {
	width := fb.w

	if (x < 0 && xs.Len(text) <= -x) || x >= fb.w || y < 0 || y >= fb.h {
		return
	}

	if x < 0 {
		xx := -x
		x = 0
		text = xs.Slice(text, xx, -1)
	}
	text = CutText(text, width)

	dx := 0
	for _, char := range text {
		s := term.Cell{Ch: char, Fg: fg, Bg: bg}
		if y >= 0 && y < fb.h && x+dx >= 0 && x+dx < fb.w {
			fb.buffer[y][x+dx] = s
		}
		dx++
	}
}

// PutVerticalText draws vertical string on FrameBuffer clipping by
// FrameBuffer boundaries. x and y are starting point, text is a string
// to display, fg and bg are text and background attributes
func (fb *FrameBuffer) PutVerticalText(x, y int, text string, fg, bg term.Attribute) {
	height := fb.h

	if (y < 0 && xs.Len(text) <= -y) || x < 0 || y < 0 || x >= fb.w {
		return
	}

	if y < 0 {
		yy := -y
		y = 0
		text = xs.Slice(text, yy, -1)
	}
	text = CutText(text, height)

	dy := 0
	for _, char := range text {
		s := term.Cell{Ch: char, Fg: fg, Bg: bg}
		fb.buffer[y+dy][x] = s
		dy++
	}
}

/*
PutColorizedText draws multicolor string on Canvas clipping by Canvas boundaries.
Multiline is not supported.
Various parts of text can be colorized with html-like tags. Every tag must start with '<'
followed by tag type and colon(without space between them), atrribute in human redable form,
and closing '>'.
Available tags are:
'f' & 't' - sets new text color
'b' - sets new background color
Available attributes (it is possible to write a few attributes for one tag separated with space):
empty string - reset the color to default value (that is passed as argument)
'bold' or 'bright' - bold or brigther color(depends on terminal)
'underline' and 'underlined' - underined text(not every terminal can do it)
'reversed' - reversed text and background
other available attributes are color names: black, red, green, yellow, blue, magenta, cyan, white.

Example: PutColorizedText(0, 0, 10, "<t:red bold><b:yellow>E<t:>xample, ColorBlack, ColorWhite, Horizontal)
It displays red letter 'C' on a yellow background, then switch text color to default one and draws
other letters in black text and yellow background colors. Default background color is not used, so
it can be set as ColroDefault in a method call
*/
func (fb *FrameBuffer) PutColorizedText(x, y, max int, text string, fg, bg term.Attribute, dir Direction) {
	var dx, dy int
	if dir == Horizontal {
		dx = 1
	} else {
		dy = 1
	}

	parser := NewColorParser(text, fg, bg)
	elem := parser.NextElement()
	for elem.Type != ElemEndOfText && max > 0 {
		if elem.Type == ElemPrintable {
			fb.PutChar(x, y, elem.Ch, elem.Fg, elem.Bg)
			x += dx
			y += dy
			max--
		}

		elem = parser.NextElement()
	}
}

/*
DrawFrame paints a frame inside FrameBuffer with optional border
rune set(by default, in case of border is empty string, the rune set
equals "─│┌┐└┘" - single border). The inner area of frame is not
filled - in other words it is transparent
*/
func (fb *FrameBuffer) DrawFrame(x, y, w, h int, fg, bg term.Attribute, frameChars string) {
	if h < 1 || w < 1 {
		return
	}

	if frameChars == "" {
		frameChars = "─│┌┐└┘"
	}

	parts := []rune(frameChars)
	if len(parts) < 6 {
		panic("Invalid theme: single border")
	}

	H, V, UL, UR, DL, DR := parts[0], parts[1], parts[2], parts[3], parts[4], parts[5]

	if h == 1 {
		for xx := x; xx < x+w; xx++ {
			fb.PutChar(xx, y, H, fg, bg)
		}
		return
	}
	if w == 1 {
		for yy := y; yy < y+h; yy++ {
			fb.PutChar(x, yy, V, fg, bg)
		}
		return
	}

	fb.PutChar(x, y, UL, fg, bg)
	fb.PutChar(x, y+h-1, DL, fg, bg)
	fb.PutChar(x+w-1, y, UR, fg, bg)
	fb.PutChar(x+w-1, y+h-1, DR, fg, bg)

	for xx := x + 1; xx < x+w-1; xx++ {
		fb.PutChar(xx, y, H, fg, bg)
		fb.PutChar(xx, y+h-1, H, fg, bg)
	}
	for yy := y + 1; yy < y+h-1; yy++ {
		fb.PutChar(x, yy, V, fg, bg)
		fb.PutChar(x+w-1, yy, V, fg, bg)
	}
}

/*
SetCursorPos sets text caret position. In opposite to other FrameBuffer
methods, x and y - are absolute console coordinates. Use negative values
if you want to hide the caret. Used by controls like EditField
*/
func (fb *FrameBuffer) SetCursorPos(x, y int) {
	term.SetCursor(x, y)
}

/*
DrawScroll paints a scroll bar inside FrameBuffer.
x, y - start position.
w, h - width and height (if h equals 1 then horizontal scroll is drawn
and vertical otherwise).
pos - thumb position.
fgScroll, bgScroll - scroll bar main attributes.
fgThumb, bgThumb - thumb colors.
scrollChars  - rune set(by default, in case of is is empty string, the
rune set equals "░■▲▼")
*/
func (fb *FrameBuffer) DrawScroll(x, y, w, h, pos int, fgScroll, bgScroll, fgThumb, bgThumb term.Attribute, scrollChars string) {
	if w < 1 || h < 1 {
		return
	}

	if scrollChars == "" {
		scrollChars = "░■▲▼◄►"
	}

	parts := []rune(scrollChars)
	chLine, chCursor, chUp, chDown := parts[0], parts[1], parts[2], parts[3]
	chLeft, chRight := '◄', '►'
	if len(parts) > 4 {
		chLeft, chRight = parts[4], parts[5]
	}

	if h == 1 {
		fb.PutChar(x, y, chLeft, fgScroll, bgScroll)
		fb.PutChar(x+w-1, y, chRight, fgScroll, bgScroll)

		if w > 2 {
			for xx := 1; xx < w-1; xx++ {
				fb.PutChar(x+xx, y, chLine, fgScroll, bgScroll)
			}
		}

		if pos != -1 {
			fb.PutChar(x+pos, y, chCursor, fgThumb, bgThumb)
		}
	} else {
		fb.PutChar(x, y, chUp, fgScroll, bgScroll)
		fb.PutChar(x, y+h-1, chDown, fgScroll, bgScroll)

		if h > 2 {
			for yy := 1; yy < h-1; yy++ {
				fb.PutChar(x, y+yy, chLine, fgScroll, bgScroll)
			}
		}

		if pos != -1 {
			fb.PutChar(x, y+pos, chCursor, fgThumb, bgThumb)
		}
	}
}
