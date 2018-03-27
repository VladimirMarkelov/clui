package clui

import (
	"fmt"
	// xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
)

/*
SparkChart is a chart that represents a live data that
is continuously added to the chart. Or it can be static
element that displays predefined set of data - in this
case it looks like BarChart. At a moment SparkChart
keeps only th enumber of last data that is enough to
fill the control area. So, if you enlarge the control,
it will show partially filled area until it gets new data.
SparkChart displays vertical axis with values on the chart left
if ValueWidth greater than 0, horizontal axis with bar titles.
Maximum peaks(maximum of the the data that control keeps)
can be hilited with different color.
By default the data is autoscaled to make the highest bar
fit the full height of the control. But it maybe useful
to disable autoscale and set the Top value to have more
handy diagram. E.g, for CPU load in % you can set
AutoScale to false and Top value to 100.
Note: negative and zero values are displayed as empty bar
*/
type SparkChart struct {
	BaseControl
	data         []float64
	valueWidth   int
	hiliteMax    bool
	maxFg, maxBg term.Attribute
	topValue     float64
	autosize     bool
}

/*
NewSparkChart creates a new spark chart.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
w and h - are minimal size of the control.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func CreateSparkChart(parent Control, w, h int, scale int) *SparkChart {
	c := new(SparkChart)
	c.BaseControl = NewBaseControl()

	if w == AutoSize {
		w = 10
	}
	if h == AutoSize {
		h = 5
	}

	c.parent = parent

	c.SetSize(w, h)
	c.SetConstraints(w, h)
	c.tabSkip = true
	c.hiliteMax = true
	c.autosize = true
	c.data = make([]float64, 0)
	c.SetScale(scale)

	if parent != nil {
		parent.AddChild(c)
	}

	return c
}

// Repaint draws the control on its View surface
func (b *SparkChart) Draw() {
	if b.hidden {
		return
	}

	b.mtx.RLock()
	defer b.mtx.RUnlock()

	PushAttributes()
	defer PopAttributes()

	fg, bg := RealColor(b.fg, ColorSparkChartText), RealColor(b.bg, ColorSparkChartBack)
	SetTextColor(fg)
	SetBackColor(bg)
	FillRect(b.x, b.y, b.width, b.height, ' ')

	if len(b.data) == 0 {
		return
	}

	b.drawValues()
	b.drawBars()
}

func (b *SparkChart) drawBars() {
	if len(b.data) == 0 {
		return
	}

	start, width := b.calculateBarArea()
	if width < 2 {
		return
	}

	coeff, max := b.calculateMultiplier()
	if coeff == 0.0 {
		return
	}

	PushAttributes()
	defer PopAttributes()

	h := b.height
	pos := b.x + start

	mxFg, mxBg := RealColor(b.maxFg, ColorSparkChartMaxText), RealColor(b.maxBg, ColorSparkChartMaxBack)
	brFg, brBg := RealColor(b.fg, ColorSparkChartBarText), RealColor(b.bg, ColorSparkChartBarBack)
	parts := []rune(SysObject(ObjSparkChart))

	var dt []float64
	if len(b.data) > width {
		dt = b.data[len(b.data)-width:]
	} else {
		dt = b.data
	}

	for _, d := range dt {
		barH := int(d * coeff)

		if barH <= 0 {
			pos++
			continue
		}

		f, g := brFg, brBg
		if b.hiliteMax && max == d {
			f, g = mxFg, mxBg
		}
		SetTextColor(f)
		SetBackColor(g)
		FillRect(pos, b.y+h-barH, 1, barH, parts[0])

		pos++
	}
}

func (b *SparkChart) drawValues() {
	if b.valueWidth <= 0 {
		return
	}

	pos, _ := b.calculateBarArea()
	if pos == 0 {
		return
	}

	h := b.height
	coeff, max := b.calculateMultiplier()
	if max == coeff {
		return
	}
	if !b.autosize || b.topValue == 0 {
		max = b.topValue
	}

	dy := 0
	format := fmt.Sprintf("%%%v.2f", b.valueWidth)
	for dy < h-1 {
		v := float64(h-dy) / float64(h) * max
		s := fmt.Sprintf(format, v)
		s = CutText(s, b.valueWidth)
		DrawRawText(b.x, b.y+dy, s)

		dy += 2
	}
}

func (b *SparkChart) calculateBarArea() (int, int) {
	w := b.width
	pos := 0

	if b.valueWidth < w/2 {
		w = w - b.valueWidth
		pos = b.valueWidth
	}

	return pos, w
}

func (b *SparkChart) calculateMultiplier() (float64, float64) {
	if len(b.data) == 0 {
		return 0, 0
	}

	h := b.height
	if h <= 1 {
		return 0, 0
	}

	max := b.data[0]
	for _, val := range b.data {
		if val > max {
			max = val
		}
	}

	if max == 0 {
		return 0, 0
	}

	if b.autosize || b.topValue == 0 {
		return float64(h) / max, max
	}
	return float64(h) / b.topValue, max
}

// AddData appends a new bar to a chart
func (b *SparkChart) AddData(val float64) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.data = append(b.data, val)

	_, width := b.calculateBarArea()
	if len(b.data) > width {
		b.data = b.data[len(b.data)-width:]
	}
}

// ClearData removes all bar from chart
func (b *SparkChart) ClearData() {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.data = make([]float64, 0)
}

// SetData assign a new bar list to a chart
func (b *SparkChart) SetData(data []float64) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.data = make([]float64, len(data))
	copy(b.data, data)

	_, width := b.calculateBarArea()
	if len(b.data) > width {
		b.data = b.data[len(b.data)-width:]
	}
}

// ValueWidth returns the width of the area at the left of
// chart used to draw values. Set it to 0 to turn off the
// value panel
func (b *SparkChart) ValueWidth() int {
	return b.valueWidth
}

// SetValueWidth changes width of the value panel on the left
func (b *SparkChart) SetValueWidth(width int) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.valueWidth = width
}

// Top returns the value of the top of a chart. The value is
// used only if autosize is off to scale all the data
func (b *SparkChart) Top() float64 {
	return b.topValue
}

// SetTop sets the theoretical highest value of data flow
// to scale the chart
func (b *SparkChart) SetTop(top float64) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.topValue = top
}

// AutoScale returns whether spark chart scales automatically
// depending on displayed data or it scales using Top value
func (b *SparkChart) AutoScale() bool {
	return b.autosize
}

// SetAutoScale changes the way of scaling the data flow
func (b *SparkChart) SetAutoScale(auto bool) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.autosize = auto
}

// HilitePeaks returns whether chart draws maximum peaks
// with different color
func (b *SparkChart) HilitePeaks() bool {
	return b.hiliteMax
}

// SetHilitePeaks enables or disables hiliting maximum
// values with different colors
func (b *SparkChart) SetHilitePeaks(hilite bool) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.hiliteMax = hilite
}
