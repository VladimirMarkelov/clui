package clui

import (
	"fmt"
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
)

// Represents an object visible content. Used by Composer to keep all screen info and by Window
type FrameBuffer struct {
	buffer [][]term.Cell
	w, h   int
}

// Creates new buffer
func NewFrameBuffer(w, h int) *FrameBuffer {
	if w < 5 || h < 5 {
		panic(fmt.Sprintf("Invalid size: %vx%v.", w, h))
	}

	c := new(FrameBuffer)
	c.SetSize(w, h)

	return c
}

// Returns current width and height of a buffer
func (fb *FrameBuffer) Size() (int, int) {
	return fb.w, fb.h
}

func (fb *FrameBuffer) SetSize(w, h int) {
	if w == fb.w && h == fb.h {
		return
	}

	if w < 5 || h < 5 {
		panic(fmt.Sprintf("Invalid size: %vx%v.", w, h))
	}

	fb.w, fb.h = w, h

	fb.buffer = make([][]term.Cell, h)
	for i := 0; i < h; i++ {
		fb.buffer[i] = make([]term.Cell, w)
	}
}

// Fills buffers with color (the buffer is filled with spaces with defined background color)
func (fb *FrameBuffer) Clear(bg term.Attribute) {
	s := term.Cell{Ch: ' ', Fg: ColorWhite, Bg: bg}
	for y := 0; y < fb.h; y++ {
		for x := 0; x < fb.w; x++ {
			fb.buffer[y][x] = s
		}
	}
}

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

func (fb *FrameBuffer) Symbol(x, y int) (term.Cell, bool) {
	if x >= fb.w || x < 0 || y >= fb.h || y < 0 {
		return term.Cell{Ch: ' '}, false
	} else {
		return fb.buffer[y][x], true
	}
}

func (fb *FrameBuffer) PutSymbol(x, y int, s term.Cell) bool {
	if x < 0 || x >= fb.w || y < 0 || y >= fb.h {
		return false
	}

	fb.buffer[y][x] = s
	return true
}

func (fb *FrameBuffer) PutChar(x, y int, c rune, fg, bg term.Attribute) bool {
	if x < 0 || x >= fb.w || y < 0 || y >= fb.h {
		return false
	}

	fb.buffer[y][x] = term.Cell{Ch: c, Fg: fg, Bg: bg}
	return true
}

// Draws text line on buffer
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
		fb.buffer[y][x+dx] = s
		dx++
	}
}

// Draws vertical text line on buffer
func (fb *FrameBuffer) PutTextVertical(x, y int, text string, fg, bg term.Attribute) {
	height := fb.h

	if (y < 0 && xs.Len(text) <= -y) || x >= fb.w || y > 0 || y >= fb.h {
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

// x, y - absolute console coordinates
func (fb *FrameBuffer) SetCursorPos(x, y int) {
	term.SetCursor(x, y)
}
