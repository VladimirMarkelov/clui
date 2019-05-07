package clui

import (
	мКнст "./пакКонстанты"
	term "github.com/nsf/termbox-go"
)

// ThumbPosition returns a scrollbar thumb position depending
// on currently active item(itemNo), total number of items
// (itemCount), and length/height of the scrollbar(length)
// including arrows. Returns position in interval of (1..lenght-2)
// or -1 if the thumb is not visible
func ThumbPosition(itemNo, itemCount, length int) int {
	length -= 2
	if itemNo < 0 {
		return -1
	}
	if itemNo >= itemCount-1 {
		return length - 1
	}

	if length < 4 {
		return 0
	}

	ydiff := int(float32(itemNo) / float32(itemCount-1.0) * float32(length-1))
	return ydiff
}

// ItemByThumbPosition calculates item number by scrollbar
// thumb position. Position - thumb position inside scrollbar,
// itemCount - total number of items, lenght - lenght or heigth
// of scrollbar. Return -1 if it is not possible to calculate:
// e.g, itemCount equals zero
func ItemByThumbPosition(position, itemCount, length int) int {
	length -= 2
	if position < 1 {
		return -1
	}
	if itemCount < 1 {
		return -1
	}
	if itemCount == 1 {
		return 0
	}

	newPos := int(float32(itemCount-1)*float32(position-1)/float32(length-1) + 0.9)

	if newPos < 0 {
		newPos = 0
	} else if newPos >= itemCount {
		newPos = itemCount - 1
	}

	return newPos
}

// ChildAt returns the children of parent control that is at absolute
// coordinates x, y. Returns nil if x, y are outside parent control and
// returns parent if no child is at x, y
func ChildAt(parent Control, x, y int) Control {
	px, py := parent.Pos()
	pw, ph := parent.Size()
	if px > x || py > y || px+pw <= x || py+ph <= y {
		return nil
	}

	if len(parent.Children()) == 0 {
		return parent
	}

	var ctrl Control
	ctrl = parent
	for _, child := range parent.Children() {
		if !child.Visible() {
			continue
		}

		check := ChildAt(child, x, y)
		if check != nil {
			ctrl = check
			break
		}
	}

	return ctrl
}

// DeactivateControls makes all children of parent inactive
func DeactivateControls(parent Control) {
	for _, ctrl := range parent.Children() {
		if ctrl.Active() {
			ctrl.SetActive(false)
			ctrl.ProcessEvent(мКнст.Event{Type: мКнст.EventActivate, X: 0})
		}

		DeactivateControls(ctrl)
	}
}

// ActivateControl makes control active and disables all other children of
// the parent. Returns true if control was found and activated
func ActivateControl(parent, control Control) bool {
	DeactivateControls(parent)
	res := false
	ctrl := FindChild(parent, control)
	if ctrl != nil {
		res = true
		if !ctrl.Active() {
			ctrl.ProcessEvent(мКнст.Event{Type: мКнст.EventActivate, X: 1})
			ctrl.SetActive(true)
		}
	}

	return res
}

// FindChild returns control if it is a child of the parent and nil otherwise
func FindChild(parent, control Control) Control {
	var res Control

	if parent == control {
		return parent
	}

	for _, ctrl := range parent.Children() {
		if ctrl == control {
			res = ctrl
			break
		}

		res = FindChild(ctrl, control)
		if res != nil {
			break
		}
	}

	return res
}

// IsMouseClickEvent returns if a user action can be treated as mouse click.
func IsMouseClickEvent(ev мКнст.Event) bool {
	if ev.Type == мКнст.EventClick {
		return true
	}
	if ev.Type == мКнст.EventMouse && ev.Key == term.MouseLeft {
		return true
	}

	return false
}

// FindFirstControl returns the first child for that fn returns true.
// The function is used to find active or tab-stop control
func FindFirstControl(parent Control, fn func(Control) bool) Control {
	linear := getLinearControlList(parent, fn)
	if len(linear) == 0 {
		return nil
	}

	return linear[0]
}

// FindLastControl returns the first child for that fn returns true.
// The function is used by TAB processing method if a user goes backwards
// with TAB key - not supported now
func FindLastControl(parent Control, fn func(Control) bool) Control {
	linear := getLinearControlList(parent, fn)

	if len(linear) == 0 {
		return nil
	}

	return linear[len(linear)-1]
}

// ActiveControl returns the active child of the parent or nil if no child is
// active
func ActiveControl(parent Control) Control {
	fnActive := func(c Control) bool {
		return c.Active()
	}
	return FindFirstControl(parent, fnActive)
}

// FindFirstActiveControl returns the first active control of a parent
func FindFirstActiveControl(parent Control) Control {
	for _, curr := range getLinearControlList(parent, nil) {
		if curr.Active() {
			return curr
		}
	}
	return nil
}

func getLinearControlList(parent Control, fn func(Control) bool) []Control {
	result := []Control{}

	for _, curr := range parent.Children() {
		if fn != nil && fn(curr) {
			result = append(result, curr)
		}

		if len(curr.Children()) == 0 {
			continue
		}

		ch := getLinearControlList(curr, fn)
		if len(ch) != 0 {
			result = append(result, ch...)
		}
	}

	return result
}

// NextControl returns the next or previous child (depends on next parameter)
// that has tab-stop feature on. Used by library when processing TAB key
func NextControl(parent Control, curr Control, next bool) Control {
	fnTab := func(c Control) bool {
		isVisible := func() bool {
			ctrl := c.Parent()

			for ctrl != nil {
				if !ctrl.Visible() {
					return false
				}

				ctrl = ctrl.Parent()
			}

			return c.Visible()
		}

		return c.TabStop() && isVisible() && c.Enabled()
	}

	linear := getLinearControlList(parent, fnTab)

	if len(linear) == 0 {
		return nil
	}

	var pIndex, nIndex int

	for i, ch := range linear {
		if ch != curr {
			continue
		}

		pIndex = i - 1
		nIndex = i + 1
		break
	}

	if nIndex > len(linear)-1 {
		nIndex = 0
	}

	if pIndex < 0 {
		pIndex = len(linear) - 1
	}

	if next {
		return linear[nIndex]
	}
	return linear[pIndex]

}

// SendEventToChild tries to find a child control that should recieve the evetn
// For mouse click events it looks for a control at coordinates of event,
// makes it active, and then sends the event to it.
// If it is not mouse click event then it looks for the first active child and
// sends the event to it if it is not nil
func SendEventToChild(parent Control, ev мКнст.Event) bool {
	var child Control
	if IsMouseClickEvent(ev) {
		child = ChildAt(parent, ev.X, ev.Y)
		if child != nil && !child.Active() {
			ActivateControl(parent, child)
		}
	} else {
		child = ActiveControl(parent)
	}

	if child != nil && child != parent {
		ev.Target = child
		res := child.ProcessEvent(ev)

		if cparent := ClippedParent(child); cparent != nil && cparent != child {
			cparent.ProcessEvent(ev)
		}

		return res
	}

	return false
}

// CalcClipper calculates the clipper size based on the control's size, position
// and paddings
func CalcClipper(c Control) (int, int, int, int) {
	w, h := c.Size()
	x, y := c.Pos()
	px, py := c.Paddings()

	x = x + px
	y = y + py
	w = w - 2*px
	h = h - 2*py

	return x, y, w, h
}

// ClippedParent finds the first c parent with clipped flag
func ClippedParent(c Control) Control {
	var clipped Control

	ctrl := c.Parent()
	clipped = c

	for ctrl != nil {
		if ctrl.Clipped() {
			clipped = ctrl
			break
		}

		ctrl = ctrl.Parent()
	}

	return clipped
}

// ControlInRect returns true if c is within a given rect
func ControlInRect(c Control, x int, y int, w int, h int) bool {
	xx, yy := c.Pos()
	ww, hh := c.Size()

	return xx >= x && ww <= x+w && yy <= y+h &&
		yy+hh <= y+h && yy >= y && yy+h >= y
}
