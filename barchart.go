package clui

import (
	"fmt"
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
)

type BarData struct {
	Value float64
	Title string
	Fg    term.Attribute
	Bg    term.Attribute
	Ch    rune
}

type BarDataCell struct {
	Item     string
	Id       int
	Value    float64
	MaxValue float64
	Fg       term.Attribute
	Bg       term.Attribute
	Ch       rune
}

/*
Label is a decorative control that can display text in horizontal
or vertical direction. Other available text features are alignment
and multi-line ability. Text can be single- or multi-colored with
tags inside the text. Multi-colored strings have limited support
of alignment feature: if text is longer than Label width the text
is always left aligned
*/
type BarChart struct {
	ControlBase
	data        []BarData
	autosize    bool
	gap         int
	barWidth    int
	legendWidth int
	valueWidth  int
	showMarks   bool
	showTicks   bool
	showTitles  bool
	onDrawCell  func(*BarDataCell)
}

/*
NewLabel creates a new label.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
w and h - are minimal size of the control.
title - is Label title.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func NewBarChart(view View, parent Control, w, h int, scale int) *BarChart {
	c := new(BarChart)

	if w == AutoSize {
		w = 10
	}
	if h == AutoSize {
		h = 5
	}

	c.view = view
	c.parent = parent

	c.SetSize(w, h)
	c.SetConstraints(w, h)
	c.tabSkip = true
	c.showTitles = true
	c.barWidth = 3
	c.data = make([]BarData, 0)

	if parent != nil {
		parent.AddChild(c, scale)
	}

	return c
}

// Repaint draws the control on its View surface
func (b *BarChart) Repaint() {
	canvas := b.view.Canvas()
	tm := b.view.Screen().Theme()

	fg, bg := RealColor(tm, b.fg, ColorBarChartText), RealColor(tm, b.bg, ColorBarChartBack)
	canvas.FillRect(b.x, b.y, b.width, b.height, term.Cell{Ch: ' ', Fg: fg, Bg: bg})

	if len(b.data) == 0 {
		return
	}

	b.drawRulers(tm, fg, bg)
	b.drawValues(tm, fg, bg)
	b.drawLegend(tm, fg, bg)
	b.drawBars(tm, fg, bg)
}

func (b *BarChart) barHeight() int {
	if b.showTitles {
		return b.height - 2
	} else {
		return b.height
	}
}

func (b *BarChart) drawBars(tm Theme, fg, bg term.Attribute) {
	if len(b.data) == 0 {
		return
	}

	start, width := b.calculateBarArea()
	if width < 2 {
		return
	}

	barW := b.calculateBarWidth()
	if barW == 0 {
		return
	}

	coeff, max := b.calculateMultiplier()
	if coeff == 0.0 {
		return
	}

	h := b.barHeight()
	pos := start
	canvas := b.view.Canvas()
	parts := []rune(tm.SysObject(ObjBarChart))

	for idx, d := range b.data {
		if pos+barW >= start+width {
			break
		}

		fColor, bColor := d.Fg, d.Bg
		ch := d.Ch
		if fColor == ColorDefault {
			fColor = fg
		}
		if bColor == ColorDefault {
			bColor = bg
		}
		if ch == 0 {
			ch = parts[0]
		}

		barH := int(d.Value * coeff)
		if b.onDrawCell == nil {
			cell := term.Cell{Ch: ch, Fg: fg, Bg: bg}
			canvas.FillRect(b.x+pos, b.y+h-barH, barW, barH, cell)
		} else {
			cellDef := BarDataCell{Item: d.Title, Id: idx, Value: d.Value, MaxValue: max, Fg: fg, Bg: bg, Ch: ch}
			for dy := 0; dy < barH; dy++ {
				req := cellDef
				b.onDrawCell(&req)
				cell := term.Cell{Ch: req.Ch, Fg: req.Fg, Bg: req.Bg}
				for dx := 0; dx < barW; dx++ {
					canvas.PutSymbol(b.x+pos+dx, b.y+h-1-dy, cell)
				}
			}
		}

		if b.showTitles {
			if b.showTicks {
				c := parts[7]
				canvas.PutSymbol(b.x+pos+barW/2, b.y+h, term.Cell{Ch: c, Bg: bg, Fg: fg})
			}
			var s string
			shift := 0
			if xs.Len(d.Title) > barW {
				s = CutText(d.Title, barW)
			} else {
				shift, s = AlignText(d.Title, barW, AlignCenter)
			}
			canvas.PutText(b.x+pos+shift, b.y+h+1, s, fg, bg)
		}

		pos += barW + b.gap
	}
}

func (b *BarChart) drawLegend(tm Theme, fg, bg term.Attribute) {
	pos, width := b.calculateBarArea()
	if pos+width >= b.width-3 {
		return
	}

	canvas := b.view.Canvas()
	for idx, d := range b.data {
		if idx >= b.height {
			break
		}

		canvas.PutSymbol(b.x+pos+width, b.y+idx, term.Cell{Ch: d.Ch, Fg: d.Fg, Bg: d.Bg})
		s := CutText(fmt.Sprintf(" - %v", d.Title), b.legendWidth)
		canvas.PutText(b.x+pos+width+1, b.y+idx, s, fg, bg)
	}
}

func (b *BarChart) drawValues(tm Theme, fg, bg term.Attribute) {
	pos, _ := b.calculateBarArea()
	if pos == 0 {
		return
	}

	h := b.barHeight()
	coeff, max := b.calculateMultiplier()
	if max == coeff {
		return
	}

	canvas := b.view.Canvas()
	dy := 0
	format := fmt.Sprintf("%%%vf", b.valueWidth)
	for dy < h-1 {
		v := float64(h-dy) / float64(h) * max
		s := fmt.Sprintf(format, v)
		s = CutText(s, b.valueWidth)
		canvas.PutText(b.x, b.y+dy, s, fg, bg)

		dy += 2
	}
}

func (b *BarChart) drawRulers(tm Theme, fg, bg term.Attribute) {
	if b.valueWidth <= 0 && b.legendWidth <= 0 && !b.showTitles {
		return
	}

	pos, vWidth := b.calculateBarArea()

	parts := []rune(tm.SysObject(ObjBarChart))
	h := b.barHeight()

	if pos > 0 {
		pos--
		vWidth++
	}

	// horizontal and vertical lines, corner
	cH, cV, cC := parts[1], parts[2], parts[5]
	canvas := b.view.Canvas()

	if pos > 0 {
		for dy := 0; dy < h; dy++ {
			canvas.PutSymbol(b.x+pos, b.y+dy, term.Cell{Ch: cV, Fg: fg, Bg: bg})
		}
	}
	if b.showTitles {
		for dx := 0; dx < vWidth; dx++ {
			canvas.PutSymbol(b.x+pos+dx, b.y+h, term.Cell{Ch: cH, Fg: fg, Bg: bg})
		}
	}
	if pos > 0 && b.showTitles {
		canvas.PutSymbol(b.x+pos, b.y+h, term.Cell{Ch: cC, Fg: fg, Bg: bg})
	}
}

func (b *BarChart) calculateBarArea() (int, int) {
	w := b.width
	pos := 0

	if b.valueWidth < w/2 {
		w = w - b.valueWidth - 1
		pos = b.valueWidth + 1
	}

	if b.legendWidth < w/2 {
		w -= b.legendWidth
	}

	return pos, w
}

func (b *BarChart) calculateBarWidth() int {
	if len(b.data) == 0 {
		return 0
	}

	if !b.autosize {
		return b.barWidth
	}

	w := b.width
	if b.valueWidth < w/2 {
		w = w - b.valueWidth - 1
	}
	if b.legendWidth < w/2 {
		w -= b.legendWidth
	}

	dataCount := len(b.data)
	minSize := dataCount*b.barWidth + (dataCount-1)*b.gap
	if minSize >= w {
		return b.barWidth
	}

	return (w - (dataCount-1)*b.gap) / dataCount
}

func (b *BarChart) calculateMultiplier() (float64, float64) {
	if len(b.data) == 0 {
		return 0, 0
	}

	h := b.barHeight()
	if h <= 1 {
		return 0, 0
	}

	max := b.data[0].Value
	for _, val := range b.data {
		if val.Value > max {
			max = val.Value
		}
	}

	if max == 0 {
		return 0, 0
	}

	return float64(h) / max, max
}

func (b *BarChart) AddData(val BarData) {
	b.data = append(b.data, val)
}

func (b *BarChart) ClearData() {
	b.data = make([]BarData, 0)
}

func (b *BarChart) SetData(data []BarData) {
	b.data = make([]BarData, len(data))
	copy(b.data, data)
}

func (b *BarChart) AutoSize() bool {
	return b.autosize
}

func (b *BarChart) SetAutoSize(auto bool) {
	b.autosize = auto
}

func (b *BarChart) Gap() int {
	return b.gap
}

func (b *BarChart) SetGap(gap int) {
	b.gap = gap
}

func (b *BarChart) MinBarSize() int {
	return b.barWidth
}

func (b *BarChart) SetMinBarSize(size int) {
	b.barWidth = size
}

func (b *BarChart) ValueWidth() int {
	return b.valueWidth
}

func (b *BarChart) SetValueWidth(width int) {
	b.valueWidth = width
}

func (b *BarChart) ShowTitles() bool {
	return b.showTitles
}

func (b *BarChart) SetShowTitles(show bool) {
	b.showTitles = show
}

func (b *BarChart) LegendWidth() int {
	return b.legendWidth
}

func (b *BarChart) SetLegendWidth(width int) {
	b.legendWidth = width
}

func (b *BarChart) OnDrawCell(fn func(*BarDataCell)) {
	b.onDrawCell = fn
}

func (b *BarChart) ShowMarks() bool {
	return b.showMarks
}

func (b *BarChart) SetShowMarks(show bool) {
	b.showMarks = show
}
