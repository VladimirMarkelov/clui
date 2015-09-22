package clui

import (
	"fmt"
	xs "github.com/huandu/xstrings"
)

// Represents an object visible content. Used by Composer to keep all screen info and by Window
type FrameBuffer struct {
	buffer [][]Symbol
	w, h   int
}

// Returns current width and height of a buffer
func (fb *FrameBuffer) GetSize() (int, int) {
	return fb.w, fb.h
}

// Creates new buffer
func NewFrameBuffer(w, h int) *FrameBuffer {
	if w < 5 || h < 5 {
		panic(fmt.Sprintf("Invalid buffer size: %vx%v.", w, h))
	}

	c := new(FrameBuffer)

	c.w = w
	c.h = h

	c.buffer = make([][]Symbol, h)
	for i := 0; i < h; i++ {
		c.buffer[i] = make([]Symbol, w)
	}

	return c
}

// Fills buffers with color (the buffer is filled with spaces with defined background color)
func (fb *FrameBuffer) Clear(bg Color) {
	s := Symbol{ch: ' ', fg: ColorWhite, bg: bg}
	for y := 0; y < fb.h; y++ {
		for x := 0; x < fb.w; x++ {
			fb.buffer[y][x] = s
		}
	}
}

func (fb *FrameBuffer) GetSymbol(x, y int) (Symbol, bool) {
	if x >= fb.w || x < 0 || y >= fb.h || y < 0 {
		return Symbol{ch: ' '}, false
	} else {
		return fb.buffer[y][x], true
	}
}

// Draws a text on a buffer, the method does not do any checks, never use it directly
func (fb *FrameBuffer) rawTextOut(x, y int, text string, fg, bg Color) {
	dx := 0
	for _, char := range text {
		s := Symbol{ch: char, fg: fg, bg: bg}
		fb.buffer[y][x+dx] = s
		dx++
	}
}

// Draws a character on a buffer. There is no check, so never use it directly
func (fb *FrameBuffer) rawRuneOut(x, y int, r rune, fg, bg Color) {
	s := Symbol{ch: r, fg: fg, bg: bg}
	fb.buffer[y][x] = s
}

// Draws text line on buffer
func (fb *FrameBuffer) PutText(x, y int, text string, fg, bg Color) {
	// TODO: check for invalid x and y, and cutting must pay attention to x and y
	width := fb.w
	fb.rawTextOut(x, y, CutText(text, width), fg, bg)
}

func countBits(bt int) int {
	cnt := 0
	for i := uint(0); i < 16; i++ {
		if (1<<i)&bt != 0 {
			cnt++
		}
	}

	return cnt
}

// --------------------------------------------------------------
// view related output
// --------------------------------------------------------------

// Draws a frame around a Window. Window cannot be borderless, active Windows always has double border and others have single border
func (fb *FrameBuffer) DrawBorder(view Window, tm *ThemeManager, fg, bg Color) {
	var bs BorderStyle
	if view.GetActive() {
		bs = BorderDouble
	} else {
		bs = BorderSingle
	}

	var cH, cV, cUL, cUR, cDL, cDR rune
	if bs == BorderSingle {
		cH = tm.GetSysObject(ObjSingleBorderHLine)
		cV = tm.GetSysObject(ObjSingleBorderVLine)
		cUL = tm.GetSysObject(ObjSingleBorderULCorner)
		cUR = tm.GetSysObject(ObjSingleBorderURCorner)
		cDL = tm.GetSysObject(ObjSingleBorderDLCorner)
		cDR = tm.GetSysObject(ObjSingleBorderDRCorner)
	} else {
		cH = tm.GetSysObject(ObjDoubleBorderHLine)
		cV = tm.GetSysObject(ObjDoubleBorderVLine)
		cUL = tm.GetSysObject(ObjDoubleBorderULCorner)
		cUR = tm.GetSysObject(ObjDoubleBorderURCorner)
		cDL = tm.GetSysObject(ObjDoubleBorderDLCorner)
		cDR = tm.GetSysObject(ObjDoubleBorderDRCorner)
	}

	w := fb.w
	h := fb.h

	fb.rawRuneOut(0, 0, cUL, fg, bg)
	fb.rawRuneOut(0, h-1, cDL, fg, bg)
	fb.rawRuneOut(w-1, 0, cUR, fg, bg)
	fb.rawRuneOut(w-1, h-1, cDR, fg, bg)

	for x := 1; x < w-1; x++ {
		fb.rawRuneOut(x, 0, cH, fg, bg)
		fb.rawRuneOut(x, h-1, cH, fg, bg)
	}
	for y := 1; y < h-1; y++ {
		fb.rawRuneOut(0, y, cV, fg, bg)
		fb.rawRuneOut(w-1, y, cV, fg, bg)
	}
}

// Draws Window buttons in its title (only if a window has border)
func (fb *FrameBuffer) DrawBorderIcons(view Window, tm *ThemeManager, fg, bg Color) {
	bi := view.GetBorderIcons()
	cnt := countBits(int(bi))
	bs := view.GetBorderStyle()

	if bs == BorderNone {
		return
	}

	x := fb.w - 2
	if cnt == 0 || x <= cnt {
		return
	}

	cOpen := tm.GetSysObject(ObjIconOpen)
	cClose := tm.GetSysObject(ObjIconClose)
	cHide := tm.GetSysObject(ObjIconMinimize)
	cDestroy := tm.GetSysObject(ObjIconDestroy)

	fb.rawRuneOut(x, 0, cClose, fg, bg)
	x--
	if bi&IconClose != 0 {
		fb.rawRuneOut(x, 0, cDestroy, fg, bg)
		x--
	}
	if bi&IconBottom != 0 {
		fb.rawRuneOut(x, 0, cHide, fg, bg)
		x--
	}
	fb.rawRuneOut(x, 0, cOpen, fg, bg)
}

// Draws Window title on its frame (only if a window has border)
func (fb *FrameBuffer) DrawTitle(view Window, fg, bg Color) {
	bs := view.GetBorderStyle()
	if bs == BorderNone {
		return
	}

	title := view.GetTitle()
	w := fb.w
	if bs != BorderNone {
		w -= 2
	}

	bi := view.GetBorderIcons()
	cnt := countBits(int(bi))
	w -= cnt + 2

	if w < 1 {
		return
	}

	text := Ellipsize(title, w)
	fb.rawTextOut(1, 0, text, fg, bg)
}

// Draw a text inside a Window.
// X and Y are local Window coordinates(coordinate system starts from top left Window corner - it is 0,0)
func (fb *FrameBuffer) DrawText(view Window, x, y int, text string, fg, bg Color) {
	bs := view.GetBorderStyle()
	dx, dy := 0, 0
	w := fb.w
	h := fb.h
	if bs != BorderNone {
		dx, dy = 1, 1
		w -= 2
		h -= 2
	}

	if x >= w || y >= h || y < 0 || x+len(text) < 0 {
		return
	}

	length := xs.Len(text)
	if x+length < 0 {
		return
	}

	if x < 0 {
		text = xs.Slice(text, -x, -1)
	}

	if x+length >= w {
		text = xs.Slice(text, 0, w-x)
	}

	fb.rawTextOut(x+dx, y+dy, text, fg, bg)
}
