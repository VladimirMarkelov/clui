package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	"strings"
	"unicode"
)

type attr struct {
	text term.Attribute
	back term.Attribute
}

type rect struct {
	x, y, w, h int
}

/*
Canvas is a 'graphical' engine to draw primitives.
*/
type Canvas struct {
	width     int
	height    int
	textColor term.Attribute
	backColor term.Attribute
	clipX     int
	clipY     int
	clipW     int
	clipH     int
	attrStack []attr
	clipStack []rect
}

var (
	canvas *Canvas
)

func initCanvas() bool {
	err := term.Init()
	if err != nil {
		return false
	}
	term.SetInputMode(term.InputEsc | term.InputMouse)

	canvas = new(Canvas)
	Reset()

	return true
}

// PushAttributes saves the current back and fore colors. Useful when used with
// PopAttributes: you can save colors then change them to anything you like and
// as the final step just restore original colors
func PushAttributes() {
	p := attr{text: canvas.textColor, back: canvas.backColor}
	canvas.attrStack = append(canvas.attrStack, p)
}

// PopAttributes restores saved with PushAttributes colors. Function does
// nothing if there is no saved colors
func PopAttributes() {
	if len(canvas.attrStack) == 0 {
		return
	}
	a := canvas.attrStack[len(canvas.attrStack)-1]
	canvas.attrStack = canvas.attrStack[:len(canvas.attrStack)-1]
	SetTextColor(a.text)
	SetBackColor(a.back)
}

// PushClip saves the current clipping window
func PushClip() {
	c := rect{x: canvas.clipX, y: canvas.clipY, w: canvas.clipW, h: canvas.clipH}
	canvas.clipStack = append(canvas.clipStack, c)
}

// PopClip restores saved with PushClip clipping window
func PopClip() {
	if len(canvas.clipStack) == 0 {
		return
	}
	c := canvas.clipStack[len(canvas.clipStack)-1]
	canvas.clipStack = canvas.clipStack[:len(canvas.clipStack)-1]
	SetClipRect(c.x, c.y, c.w, c.h)
}

// Reset reinitializes canvas: set clipping rectangle to the whole
// terminal window, clears clip and color saved data, sets colors
// to default ones
func Reset() {
	canvas.width, canvas.height = term.Size()
	canvas.clipX, canvas.clipY = 0, 0
	canvas.clipW, canvas.clipH = canvas.width, canvas.height
	canvas.textColor = ColorWhite
	canvas.backColor = ColorBlack

	canvas.attrStack = make([]attr, 0)
	canvas.clipStack = make([]rect, 0)
}

// InClipRect returns true if x and y position is inside current clipping
// rectangle
func InClipRect(x, y int) bool {
	return x >= canvas.clipX && y >= canvas.clipY &&
		x < canvas.clipX+canvas.clipW &&
		y < canvas.clipY+canvas.clipH
}

func clip(x, y, w, h int) (cx int, cy int, cw int, ch int) {
	if x+w < canvas.clipX || x > canvas.clipX+canvas.clipW ||
		y+h < canvas.clipY || y > canvas.clipY+canvas.clipH {
		return 0, 0, 0, 0
	}

	if x < canvas.clipX {
		w = w - (canvas.clipX - x)
		x = canvas.clipX
	}
	if y < canvas.clipY {
		h = h - (canvas.clipY - y)
		y = canvas.clipY
	}
	if x+w > canvas.clipX+canvas.clipW {
		w = canvas.clipW - (x - canvas.clipX)
	}
	if y+h > canvas.clipY+canvas.clipH {
		h = canvas.clipH - (y - canvas.clipY)
	}

	return x, y, w, h
}

// Flush makes termbox to draw everything to screen
func Flush() {
	term.Flush()
}

// SetSize sets the new Canvas size. If new size does not
// equal old size then Canvas is recreated and cleared
// with default colors. Both Canvas width and height must
// be greater than 2
func SetScreenSize(width int, height int) {
	if canvas.width == width && canvas.height == height {
		return
	}

	canvas.width = width
	canvas.height = height

	canvas.clipStack = make([]rect, 0)
	SetClipRect(0, 0, width, height)
}

// Size returns current Canvas size
func ScreenSize() (width int, height int) {
	return canvas.width, canvas.height
}

// SetCursorPos sets text caret position. Used by controls like EditField
func SetCursorPos(x int, y int) {
	term.SetCursor(x, y)
}

// PutChar sets value for the Canvas cell: rune and its colors. Returns result of
// operation: e.g, if the symbol position is outside Canvas the operation fails
// and the function returns false
func PutChar(x, y int, r rune) bool {
	if InClipRect(x, y) {
		term.SetCell(x, y, r, canvas.textColor, canvas.backColor)
		return true
	}

	return false
}

func putCharUnsafe(x, y int, r rune) {
	term.SetCell(x, y, r, canvas.textColor, canvas.backColor)
}

// Symbol returns the character and its attributes by its coordinates
func Symbol(x, y int) (term.Cell, bool) {
	if x >= 0 && x < canvas.width && y >= 0 && y < canvas.height {
		cells := term.CellBuffer()
		return cells[y*canvas.width+x], true
	}
	return term.Cell{Ch: ' '}, false
}

// SetTextColor changes current text color
func SetTextColor(clr term.Attribute) {
	canvas.textColor = clr
}

// SetBackColor changes current background color
func SetBackColor(clr term.Attribute) {
	canvas.backColor = clr
}

func TextColor() term.Attribute {
	return canvas.textColor
}

func BackColor() term.Attribute {
	return canvas.backColor
}

// SetClipRect defines a new clipping rect. Maybe useful with PopClip and
// PushClip functions
func SetClipRect(x, y, w, h int) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x+w > canvas.width {
		w = canvas.width - x
	}
	if y+h > canvas.height {
		h = canvas.height - h
	}

	canvas.clipX = x
	canvas.clipY = y
	canvas.clipW = w
	canvas.clipH = h
}

// ClipRect returns the current clipping rectangle
func ClipRect() (x int, y int, w int, h int) {
	return canvas.clipX, canvas.clipY, canvas.clipW, canvas.clipH
}

// DrawHorizontalLine draws the part of the horizontal line that is inside
// current clipping rectangle
func DrawHorizontalLine(x, y, w int, r rune) {
	x, y, w, _ = clip(x, y, w, 1)
	if w == 0 {
		return
	}

	for i := x; i < x+w; i++ {
		putCharUnsafe(i, y, r)
	}
}

// DrawVerticalLine draws the part of the vertical line that is inside current
// clipping rectangle
func DrawVerticalLine(x, y, h int, r rune) {
	x, y, _, h = clip(x, y, 1, h)
	if h == 0 {
		return
	}

	for i := y; i < y+h; i++ {
		putCharUnsafe(x, i, r)
	}
}

// DrawText draws the part of text that is inside the current clipping
// rectangle. DrawText always paints colorized string. If you want to draw
// raw string then use DrawRawText function
func DrawText(x, y int, text string) {
	PushAttributes()
	defer PopAttributes()

	defText, defBack := TextColor(), BackColor()
	firstdrawn := InClipRect(x, y)

	parser := NewColorParser(text, defText, defBack)
	elem := parser.NextElement()
	for elem.Type != ElemEndOfText {
		if elem.Type == ElemPrintable {
			SetTextColor(elem.Fg)
			SetBackColor(elem.Bg)
			drawn := PutChar(x, y, elem.Ch)

			if unicode.Is(unicode.Scripts["Han"], elem.Ch) {
				x += 2
			} else {
				x += 1
			}

			if firstdrawn && !drawn {
				break
			}
		}

		elem = parser.NextElement()
	}
}

// DrawRawText draws the part of text that is inside the current clipping
// rectangle. DrawRawText always paints string as is - no color changes.
// If you want to draw string with color changing commands included then
// use DrawText function
func DrawRawText(x, y int, text string) {
	cx, cy, cw, ch := ClipRect()
	if x >= cx+cw || y < cy || y >= cy+ch {
		return
	}

	length := xs.Len(text)
	if x+length < cx {
		return
	}

	if x < cx {
		text = xs.Slice(text, cx-x, -1)
		length = length - (cx - x)
		x = cx
	}
	text = CutText(text, cw)

	dx := 0
	for _, ch := range text {
		putCharUnsafe(x+dx, y, ch)
		dx++
	}
}

// DrawTextVertical draws the part of text that is inside the current clipping
// rectangle. DrawTextVertical always paints colorized string. If you want to draw
// raw string then use DrawRawTextVertical function
func DrawTextVertical(x, y int, text string) {
	PushAttributes()
	defer PopAttributes()

	defText, defBack := TextColor(), BackColor()
	firstdrawn := InClipRect(x, y)

	parser := NewColorParser(text, defText, defBack)
	elem := parser.NextElement()
	for elem.Type != ElemEndOfText {
		if elem.Type == ElemPrintable {
			SetTextColor(elem.Fg)
			SetBackColor(elem.Bg)
			drawn := PutChar(x, y, elem.Ch)
			y += 1
			if firstdrawn && !drawn {
				break
			}
		}

		elem = parser.NextElement()
	}
}

// DrawRawTextVertical draws the part of text that is inside the current clipping
// rectangle. DrawRawTextVertical always paints string as is - no color changes.
// If you want to draw string with color changing commands included then
// use DrawTextVertical function
func DrawRawTextVertical(x, y int, text string) {
	cx, cy, cw, ch := ClipRect()
	if y >= cy+ch || x < cx || x >= cx+cw {
		return
	}

	length := xs.Len(text)
	if y+length < cy {
		return
	}

	if y < cy {
		text = xs.Slice(text, cy-y, -1)
		length = length - (cy - y)
		y = cy
	}
	text = CutText(text, ch)

	dy := 0
	for _, ch := range text {
		putCharUnsafe(x, y+dy, ch)
		dy++
	}
}

// DrawFrame paints the frame without changing area inside it
func DrawFrame(x, y, w, h int, border BorderStyle) {
	var chars string
	if border == BorderThick {
		chars = SysObject(ObjDoubleBorder)
	} else if border == BorderThin {
		chars = SysObject(ObjSingleBorder)
	} else if border == BorderNone {
		chars = "      "
	} else {
		chars = "      "
	}

	parts := []rune(chars)
	H, V, UL, UR, DL, DR := parts[0], parts[1], parts[2], parts[3], parts[4], parts[5]

	if InClipRect(x, y) {
		putCharUnsafe(x, y, UL)
	}
	if InClipRect(x+w-1, y+h-1) {
		putCharUnsafe(x+w-1, y+h-1, DR)
	}
	if InClipRect(x, y+h-1) {
		putCharUnsafe(x, y+h-1, DL)
	}
	if InClipRect(x+w-1, y) {
		putCharUnsafe(x+w-1, y, UR)
	}

	var xx, yy, ww, hh int
	xx, yy, ww, _ = clip(x+1, y, w-2, 1)
	if ww > 0 {
		DrawHorizontalLine(xx, yy, ww, H)
	}
	xx, yy, ww, _ = clip(x+1, y+h-1, w-2, 1)
	if ww > 0 {
		DrawHorizontalLine(xx, yy, ww, H)
	}
	xx, yy, _, hh = clip(x, y+1, 1, h-2)
	if hh > 0 {
		DrawVerticalLine(xx, yy, hh, V)
	}
	xx, yy, _, hh = clip(x+w-1, y+1, 1, h-2)
	if hh > 0 {
		DrawVerticalLine(xx, yy, hh, V)
	}
}

// DrawScrollBar displays a scrollbar. pos is the position of the thumb.
// The function detects direction of the scrollbar automatically: if w is greater
// than h then it draws horizontal scrollbar and vertical otherwise
func DrawScrollBar(x, y, w, h, pos int) {
	xx, yy, ww, hh := clip(x, y, w, h)
	if ww < 1 || hh < 1 {
		return
	}

	PushAttributes()
	defer PopAttributes()

	fg, bg := RealColor(ColorDefault, "", ColorScrollText), RealColor(ColorDefault, "", ColorScrollBack)
	// TODO: add thumb styling
	// fgThumb, bgThumb := RealColor(ColorDefault, "", ColorThumbText), RealColor(ColorDefault, "", ColorThumbBack)
	SetTextColor(fg)
	SetBackColor(bg)

	parts := []rune(SysObject(ObjScrollBar))
	chLine, chThumb, chUp, chDown := parts[0], parts[1], parts[2], parts[3]
	chLeft, chRight := parts[4], parts[5]

	chStart, chEnd := chUp, chDown
	var dx, dy int
	if w > h {
		chStart, chEnd = chLeft, chRight
		dx = w - 1
		dy = 0
	} else {
		dx = 0
		dy = h - 1
	}

	if InClipRect(x, y) {
		putCharUnsafe(x, y, chStart)
	}
	if InClipRect(x+dx, y+dy) {
		putCharUnsafe(x+dx, y+dy, chEnd)
	}
	if xx == x && w > h {
		xx = x + 1
		ww--
	}
	if yy == y && w < h {
		yy = y + 1
		hh--
	}
	if xx+ww == x+w && w > h {
		ww--
	}
	if yy+hh == y+h && w < h {
		hh--
	}

	if w > h {
		DrawHorizontalLine(xx, yy, ww, chLine)
	} else {
		DrawVerticalLine(xx, yy, hh, chLine)
	}

	if pos >= 0 {
		if w > h {
			if pos < w-2 && InClipRect(x+1+pos, y) {
				putCharUnsafe(x+1+pos, y, chThumb)
			}
		} else {
			if pos < h-2 && InClipRect(x, y+1+pos) {
				putCharUnsafe(x, y+1+pos, chThumb)
			}
		}
	}
}

// FillRect paints the area with r character using the current colors
func FillRect(x, y, w, h int, r rune) {
	x, y, w, h = clip(x, y, w, h)
	if w < 1 || y < -1 {
		return
	}

	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			putCharUnsafe(xx, yy, r)
		}
	}
}

// TextExtent calculates the width and the height of the text
func TextExtent(text string) (int, int) {
	if text == "" {
		return 0, 0
	}
	parts := strings.Split(text, "\n")
	h := len(parts)
	w := 0
	for _, p := range parts {
		s := UnColorizeText(p)
		l := len(s)
		if l > w {
			w = l
		}
	}

	return h, w
}
