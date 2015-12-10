package clui

import (
	"fmt"
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
)

// BarData is info about one bar in the chart. Every
// bar can be customized by setting its own colors and
// rune to draw the bar. Use ColorDefault for Fg and Bg,
// and 0 for Ch to draw with BarChart defaults
type BarData struct {
	Value float64
	Title string
	Fg    term.Attribute
	Bg    term.Attribute
	Ch    rune
}

// BarDataCell is used in callback to user to draw with
// customized colors and runes
type BarDataCell struct {
	// Title of the bar
	Item string
	// order number of the bar
	ID int
	// value of the bar that is currently drawn
	Value float64
	// maximum value of the bar
	BarMax float64
	// value of the highest bar
	TotalMax float64
	// Default attributes and rune to draw the bar
	Fg term.Attribute
	Bg term.Attribute
	Ch rune
}

/*
BarChart is a chart that represents grouped data with
rectangular bars. It can be monochrome - defaut behavior.
One can assign individual color to each bar and even use
custom drawn bars to display multicolored bars depending
on bar value.
All bars have the same width: either constant BarSize - in
case of AutoSize is false, or automatically calculated but
cannot be less than BarSize. Bars that do not fit the chart
area are not displayed.
BarChart displays vertical axis with values on the chart left
if ValueWidth greater than 0, horizontal axis with bar titles
if ShowTitles is true (to enable displaying marks on horizontal
axis, set ShowMarks to true), and chart legend on the right if
LegendWidth is greater than 3.
If LegendWidth is greater than half of the chart it is not
displayed. The same is applied to ValueWidth
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
	showTitles  bool
	onDrawCell  func(*BarDataCell)
}

/*
NewBarChart creates a new bar chart.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
w and h - are minimal size of the control.
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
	}
	return b.height
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
		if pos+barW > start+width {
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
			cell := term.Cell{Ch: ch, Fg: fColor, Bg: bColor}
			canvas.FillRect(b.x+pos, b.y+h-barH, barW, barH, cell)
		} else {
			cellDef := BarDataCell{Item: d.Title, ID: idx,
				Value: 0, BarMax: d.Value, TotalMax: max,
				Fg: fColor, Bg: bColor, Ch: ch}
			for dy := 0; dy < barH; dy++ {
				req := cellDef
				req.Value = max * float64(dy+1) / float64(h)
				b.onDrawCell(&req)
				cell := term.Cell{Ch: req.Ch, Fg: req.Fg, Bg: req.Bg}
				for dx := 0; dx < barW; dx++ {
					canvas.PutSymbol(b.x+pos+dx, b.y+h-1-dy, cell)
				}
			}
		}

		if b.showTitles {
			if b.showMarks {
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
	parts := []rune(tm.SysObject(ObjBarChart))
	defRune := parts[0]
	for idx, d := range b.data {
		if idx >= b.height {
			break
		}

		c := d.Ch
		if c == 0 {
			c = defRune
		}
		canvas.PutSymbol(b.x+pos+width, b.y+idx, term.Cell{Ch: c, Fg: d.Fg, Bg: d.Bg})
		s := CutText(fmt.Sprintf(" - %v", d.Title), b.legendWidth)
		canvas.PutText(b.x+pos+width+1, b.y+idx, s, fg, bg)
	}
}

func (b *BarChart) drawValues(tm Theme, fg, bg term.Attribute) {
	if b.valueWidth <= 0 {
		return
	}

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
	format := fmt.Sprintf("%%%v.2f", b.valueWidth)
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

	sz := (w - (dataCount-1)*b.gap) / dataCount
	if sz == 0 {
		sz = 1
	}

	return sz
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

// AddData appends a new bar to a chart
func (b *BarChart) AddData(val BarData) {
	b.data = append(b.data, val)
}

// ClearData removes all bar from chart
func (b *BarChart) ClearData() {
	b.data = make([]BarData, 0)
}

// SetData assign a new bar list to a chart
func (b *BarChart) SetData(data []BarData) {
	b.data = make([]BarData, len(data))
	copy(b.data, data)
}

// AutoSize returns whether automatic bar width
// calculation is on. If AutoSize is false then all
// bars have width BarWidth. If AutoSize is true then
// bar width is the maximum of three values: BarWidth,
// calculated width that makes all bars fit the
// bar chart area, and 1
func (b *BarChart) AutoSize() bool {
	return b.autosize
}

// SetAutoSize enables or disables automatic bar
// width calculation
func (b *BarChart) SetAutoSize(auto bool) {
	b.autosize = auto
}

// Gap returns width of visual gap between two adjacent bars
func (b *BarChart) Gap() int {
	return b.gap
}

// SetGap sets the space width between two adjacent bars
func (b *BarChart) SetGap(gap int) {
	b.gap = gap
}

// MinBarWidth returns current minimal bar width
func (b *BarChart) MinBarWidth() int {
	return b.barWidth
}

// SetMinBarWidth changes the minimal bar width
func (b *BarChart) SetMinBarWidth(size int) {
	b.barWidth = size
}

// ValueWidth returns the width of the area at the left of
// chart used to draw values. Set it to 0 to turn off the
// value panel
func (b *BarChart) ValueWidth() int {
	return b.valueWidth
}

// SetValueWidth changes width of the value panel on the left
func (b *BarChart) SetValueWidth(width int) {
	b.valueWidth = width
}

// ShowTitles returns if chart displays horizontal axis and
// bar titles under it
func (b *BarChart) ShowTitles() bool {
	return b.showTitles
}

// SetShowTitles turns on and off horizontal axis and bar titles
func (b *BarChart) SetShowTitles(show bool) {
	b.showTitles = show
}

// LegendWidth returns width of chart legend displayed at the
// right side of the chart. Set it to 0 to disable legend
func (b *BarChart) LegendWidth() int {
	return b.legendWidth
}

// SetLegendWidth sets new legend panel width
func (b *BarChart) SetLegendWidth(width int) {
	b.legendWidth = width
}

// OnDrawCell sets callback that allows to draw multicolored
// bars. BarChart sends the current attrubutes and rune that
// it is going to use to display as well as the current value
// of the bar. A user can change the values of BarDataCell
// depending on some external data or calculations - only
// changing colors and rune makes sense. Changing anything else
// does not affect the chart
func (b *BarChart) OnDrawCell(fn func(*BarDataCell)) {
	b.onDrawCell = fn
}

// ShowMarks returns if horizontal axis has mark under each
// bar. To show marks, ShowTitles must be enabled.
func (b *BarChart) ShowMarks() bool {
	return b.showMarks
}

// SetShowMarks turns on and off marks under horizontal axis
func (b *BarChart) SetShowMarks(show bool) {
	b.showMarks = show
}
