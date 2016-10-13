package clui

import (
	xs "github.com/huandu/xstrings"
)

/*
Frame is a decorative control and container - frame with optional title.
All area inside a frame is transparent. Frame can be used as spacer element
- set border to BorderNone and use that control in any place where a spacer
is required
*/
type Frame struct {
	BaseControl
	border   BorderStyle
	children []Control
	pack     PackType
}

/*
NewFrame creates a new frame.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width and heigth - are minimal size of the control.
bs - type of border: no border, single or double.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func CreateFrame(parent Control, width, height int, bs BorderStyle, scale int) *Frame {
	f := new(Frame)

	if width == AutoSize {
		width = 5
	}
	if height == AutoSize {
		height = 3
	}

	f.SetSize(width, height)
	f.SetConstraints(width, height)
	f.border = bs
	f.parent = parent
	f.SetTabStop(false)
	f.scale = scale

	f.gapX, f.gapY = 0, 0
	if bs == BorderNone {
		f.padX, f.padY = 0, 0
	} else {
		f.padX, f.padY = 1, 1
	}

	if parent != nil {
		parent.AddChild(f)
	}

	return f
}

// Repaint draws the control on its View surface
func (f *Frame) Draw() {
	PushAttributes()
	defer PopAttributes()

	if f.border == BorderNone {
		f.DrawChildren()
		return
	}

	x, y := f.Pos()
	w, h := f.Size()

	fg, bg := RealColor(f.fg, ColorViewText), RealColor(f.bg, ColorViewBack)

	SetTextColor(fg)
	SetBackColor(bg)
	DrawFrame(x, y, w, h, f.border)

	if f.title != "" {
		str := f.title
		raw := UnColorizeText(str)
		if xs.Len(raw) > w-2 {
			str = SliceColorized(str, 0, w-2-3) + "..."
		}
		DrawText(x+1, y, str)
	}

	f.DrawChildren()
}
