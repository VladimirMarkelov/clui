package clui

import (
	term "github.com/nsf/termbox-go"
	"strings"
)

// CalculateMinimalSize return the minimal width and height
// of the Control based on control's children minial sizes
// and control's paddings
func CalculateMinimalSize(c Control) (int, int) {
	w, h := c.Constraints()

	kids := c.Children()
	if len(kids) == 0 {
		return w, h
	}

	w, h = 0, 0
	pk := c.Pack()
	top, side, dx, dy := c.Paddings()

	for _, child := range kids {
		cw, ch := child.Constraints()

		if pk == Vertical {
			if h != 0 {
				h += dy
			}
			h += ch
			if cw > w {
				w = cw
			}
		} else {
			if w != 0 {
				w += dx
			}
			w += cw
			if ch > h {
				h = ch
			}
		}
	}

	w += side * 2
	h += top * 2

	return w, h
}

// CalculateTotalScale return sum of all
// control children scale coefficients
func CalculateTotalScale(c Control) int {
	scale := 0

	for _, child := range c.Children() {
		sc := child.Scale()
		if sc > DoNotScale {
			scale += sc
		}
	}

	return scale
}

// RepositionControls calculates position of all
// control children and moves them to a new positions.
// Call the funtion after the control is resized.
// dx and dy are position of the container control
// relative to parent View. Initial calculation start
// point is 0, 0. After calculation all childen are
// moved by dx and dy
func RepositionControls(dx, dy int, c Control) {
	if len(c.Children()) == 0 {
		return
	}

	scale := CalculateTotalScale(c)

	minW, minH := c.Constraints()
	width, height := c.Size()
	xDiff, yDiff := width-minW, height-minH
	pk := c.Pack()
	top, side, xx, yy := c.Paddings()

	calcShift := func(height int, ctrl Control) int {
		_, h := ctrl.Size()
		if h < height {
			return int((height - h) / 2)
		}

		return 0
	}

	currX, currY := top, side
	if pk == Vertical {
		delta := float64(yDiff) / float64(scale)
		curr := 0.0
		for _, child := range c.Children() {
			cs := child.Scale()
			cw, ch := child.Constraints()
			if cs != 0 {
				var dh int
				if cs == scale {
					dh = yDiff
				} else {
					curr += float64(cs) * delta
					dh = int(curr)
					curr = curr - float64(dh)
				}
				scale -= cs
				ch += dh
				yDiff -= dh
			}
			cw = width - 2*side
			child.SetSize(cw, ch)

			yShift := calcShift(ch, child)
			child.SetPos(currX+dx, currY+dy+yShift)
			currY += ch + yy
		}
	} else {
		delta := float64(xDiff) / float64(scale)
		curr := 0.0
		for _, child := range c.Children() {
			cs := child.Scale()
			cw, ch := child.Constraints()
			if cs != 0 {
				var dw int
				if cs == scale {
					dw = xDiff
				} else {
					curr += float64(cs) * delta
					dw = int(curr)
					curr = curr - float64(dw)
				}
				scale -= cs
				cw += dw
				// log.Printf("Child width %v: %v", child.Title(), cw)
				xDiff -= dw
			}
			ch = height - 2*top
			child.SetSize(cw, ch)

			yShift := calcShift(ch, child)
			child.SetPos(currX+dx, currY+dy+yShift)
			// log.Printf("ChildH %v pos %v:%v width: %v", child.Title(), currX+dx, currY+dy, cw)
			currX += cw + xx
		}
	}

	for _, child := range c.Children() {
		childX, childY := child.Pos()
		RepositionControls(childX, childY, child)
	}
}

// RealColor returns attribute that should be applied to an
// object. By default all attributes equal ColorDefault and
// the real color should be retrieved from the current theme.
// Attribute selection work this way: if color is not ColorDefault,
// it is returned as is, otherwise the function tries to load
// color from the theme.
// tm - the theme to retrieve color from
// clr - current object color
// id - color ID in theme
func RealColor(tm Theme, clr term.Attribute, id string) term.Attribute {
	if clr != ColorDefault {
		return clr
	}

	if clr == ColorDefault {
		return tm.SysColor(id)
	}

	return clr
}

// StringToColor returns attribute by its string description.
// Description is the list of attributes separated with
// spaces, plus or pipe symbols. You can use 8 base colors:
// black, white, red, green, blue, magenta, yellow, cyan
// and a few modifiers:
// bold or bright, underline or underlined, reverse
// Note: some terminals do not support all modifiers, e.g,
// Windows one understands only bold/bright - it makes the
// color brighter with the modidierA
// Examples: "red bold", "green+underline+bold"
func StringToColor(str string) term.Attribute {
	var parts []string
	if strings.ContainsRune(str, '+') {
		parts = strings.Split(str, "+")
	} else if strings.ContainsRune(str, '|') {
		parts = strings.Split(str, "|")
	} else if strings.ContainsRune(str, ' ') {
		parts = strings.Split(str, " ")
	} else {
		parts = append(parts, str)
	}

	var cmap = map[string]term.Attribute{
		"default":    term.ColorDefault,
		"black":      term.ColorBlack,
		"red":        term.ColorRed,
		"green":      term.ColorGreen,
		"yellow":     term.ColorYellow,
		"blue":       term.ColorBlue,
		"magenta":    term.ColorMagenta,
		"cyan":       term.ColorCyan,
		"white":      term.ColorWhite,
		"bold":       term.AttrBold,
		"bright":     term.AttrBold, // windows make color brighter when it is bold
		"underline":  term.AttrUnderline,
		"underlined": term.AttrUnderline,
		"reverse":    term.AttrReverse,
	}

	var clr term.Attribute
	for _, item := range parts {
		item = strings.Trim(item, " ")
		item = strings.ToLower(item)

		c, ok := cmap[item]
		if ok {
			clr |= c
		}
	}

	return clr
}

// ColorToString returns string representation of the attribute
func ColorToString(attr term.Attribute) string {
	var out string
	colors := []string{
		"", "black", "red", "green", "yellow",
		"blue", "magenta", "cyan", "white"}

	rawClr := attr & 15
	if rawClr < 8 {
		out += colors[rawClr] + " "
	}

	if attr&term.AttrBold != 0 {
		out += "bold "
	}
	if attr&term.AttrUnderline != 0 {
		out += "underline "
	}
	if attr&term.AttrReverse != 0 {
		out += "reverse "
	}

	return strings.TrimSpace(out)
}

// ThumbPosition returns a scrollbar thumb position depending
// on currently active item(itemNo), total number of items
// (itemCount), and length/height of the scrollbar(length)
// including arrows. Returns position in interval of (1..lenght-2)
// or -1 if the thumb is not visible
func ThumbPosition(itemNo, itemCount, length int) int {
	if itemNo < 0 {
		return -1
	}
	if itemNo >= itemCount-1 {
		return length - 2
	}

	if length < 4 {
		return 1
	}

	ydiff := int(float32(itemNo) / float32(itemCount-1.0) * float32(length-3))
	return ydiff + 1
}

// ItemByThumbPosition calculates item number by scrollbar
// thumb position. Position - thumb position inside scrollbar,
// itemCount - total number of items, lenght - lenght or heigth
// of scrollbar. Return -1 if it is not possible to calculate:
// e.g, itemCount equals zero
func ItemByThumbPosition(position, itemCount, length int) int {
	if position < 1 {
		return -1
	}
	if itemCount < 1 {
		return -1
	}
	if itemCount == 1 {
		return 1
	}

	newPos := int(float32(itemCount-1)*float32(position-1)/float32(length-3) + 0.9)

	if newPos < 0 {
		newPos = 0
	} else if newPos >= itemCount {
		newPos = itemCount - 1
	}

	return newPos
}
