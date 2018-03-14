package clui

import (
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
			ctrl.ProcessEvent(Event{Type: EventActivate, X: 0})
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
			ctrl.ProcessEvent(Event{Type: EventActivate, X: 1})
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
func IsMouseClickEvent(ev Event) bool {
	if ev.Type == EventClick {
		return true
	}
	if ev.Type == EventMouse && ev.Key == term.MouseLeft {
		return true
	}

	return false
}

// FindFirstControl returns the first child for that fn returns true.
// The function is used to find active or tab-stop control
func FindFirstControl(parent Control, fn func(Control) bool) Control {
	for _, child := range parent.Children() {
		if fn(child) {
			return child
		}

		ch := FindFirstControl(child, fn)
		if ch != nil {
			return ch
		}
	}

	return nil
}

// FindLastControl returns the first child for that fn returns true.
// The function is used by TAB processing method if a user goes backwards
// with TAB key - not supported now
func FindLastControl(parent Control, fn func(Control) bool) Control {
	var last Control
	for _, child := range parent.Children() {
		if fn(child) {
			last = child
		}

		ch := FindLastControl(child, fn)
		if ch != nil {
			last = ch
		}
	}

	return last
}

// ActiveControl returns the active child of the parent or nil if no child is
// active
func ActiveControl(parent Control) Control {
	fnActive := func(c Control) bool {
		return c.Active()
	}
	return FindFirstControl(parent, fnActive)
}

func _nextControl(parent Control, curr, prev Control, foundPrev, next bool) (bool, Control) {
	found := foundPrev
	if parent == curr {
		if next {
			found = true
		} else {
			return false, prev
		}
	}

	p := prev
	for _, ctrl := range parent.Children() {
		if ctrl == curr {
			if next {
				found = true
				continue
			} else {
				return found, p
			}
		}

		if ctrl.Enabled() && ctrl.TabStop() && ctrl.Visible() {
			if found {
				return found, ctrl
			} else if !next {
				p = ctrl
			}
		}

		fnd, nn := _nextControl(ctrl, curr, p, found, next)
		if nn != nil {
			return fnd, nn
		}
		found = fnd
	}

	return found, nil
}

// NextControl returns the next or previous child (depends on next parameter)
// that has tab-stop feature on. Used by library when processing TAB key
func NextControl(parent Control, curr Control, next bool) Control {
	fnTab := func(c Control) bool {
		return c.TabStop() && c.Visible()
	}

	var defControl Control
	if next {
		defControl = FindFirstControl(parent, fnTab)
	} else {
		defControl = FindLastControl(parent, fnTab)
	}

	if defControl == nil {
		return nil
	}
	if curr == nil {
		return defControl
	}

	_, cNext := _nextControl(parent, curr, nil, false, next)
	if cNext == nil {
		cNext = defControl
	}

	return cNext
}

// SendEventToChild tries to find a child control that should recieve the evetn
// For mouse click events it looks for a control at coordinates of event,
// makes it active, and then sends the event to it.
// If it is not mouse click event then it looks for the first active child and
// sends the event to it if it is not nil
func SendEventToChild(parent Control, ev Event) bool {
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
		return child.ProcessEvent(ev)
	}

	return false
}
