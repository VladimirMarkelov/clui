package clui

import (
	xs "github.com/huandu/xstrings"
	"math"
)

/*
Frame is a decorative control and container - frame with optional title.
All area inside a frame is transparent. Frame can be used as spacer element
- set border to BorderNone and use that control in any place where a spacer
is required
*/
type Frame struct {
	BaseControl
	border         BorderStyle
	children       []Control
	pack           PackType
	scrollable     bool
	lastScrollProp int
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
	f.BaseControl = NewBaseControl()

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

func (f *Frame) SetScrollable(scrollable bool) {
	f.scrollable = scrollable

	if scrollable {
		px, py := f.Paddings()

		if f.Pack() == Vertical {
			px += 1
		}

		if f.Pack() == Horizontal {
			py += 1
		}

		f.SetPaddings(px, py)
	}

	f.SetClipped(scrollable)
}

func (f *Frame) Scrollable() bool {
	return f.scrollable
}

// Repaint draws the control on its View surface
func (f *Frame) Draw() {
	if f.hidden {
		return
	}

	PushAttributes()
	defer PopAttributes()

	x, y, w, h := f.Clipper()

	if f.scrollable {

		_, fy := f.Pos()
		_, fh := f.Size()
		_, fpy := f.Paddings()

		var dist float64
		prop := 0
		ctrl := ActiveControl(f)

		if ctrl != nil {
			var frameProp float64

			_, ty := ctrl.Pos()

			dist = (float64(fy) + float64(fpy)) - float64(ty)
			dist = math.Sqrt(dist * dist)

			if dist > 0 {
				frameProp = (dist * 100) / float64(fh)
			}

			if frameProp > 0 {
				prop = int(math.Round((float64(h-2) / (100 / frameProp))))
			}

			f.lastScrollProp = prop
		}

		DrawScrollBar(x+w, y, 1, h, f.lastScrollProp)
	}

	fg, bg := RealColor(f.fg, f.Style(), ColorViewText), RealColor(f.bg, f.Style(), ColorViewBack)

	if f.border == BorderNone {
		if bg != ColorDefault {
			SetBackColor(bg)
			FillRect(x, y, w, h, ' ')
		}

		f.DrawChildren()
		return
	}

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

// ScrollTo in case of a scrollable frame this api will scroll the content
// without adjusting the clipper
func (f *Frame) ScrollTo(x int, y int) {
	if !f.scrollable {
		return
	}

	f.x = x
	f.y = y

	f.ResizeChildren()
	f.PlaceChildren()
}

func (f *Frame) ProcessEvent(ev Event) bool {
	if ev.Type != EventActivateChild || (!f.scrollable || ev.Target == nil) {
		return false
	}

	x, y := f.Pos()
	px, py := f.Paddings()

	cx, cy, cw, ch := f.Clipper()

	tw, th := ev.Target.Size()
	tx, ty := ev.Target.Pos()

	if ControlInRect(ev.Target, cx, cy, cw, ch) {
		return false
	}

	xx := x
	yy := y

	if (ty+th)-(py/2) > cy+ch {
		delta := (ty + th) - (cy + ch)
		yy = y - delta
	} else if ty < cy {
		delta := cy - ty
		yy = y + delta
	}

	if (tx+tw)-(px/2) > cx+cw {
		delta := (tx + tw) - (cx + cw)
		xx = (x - delta)
	} else if tx < cx {
		delta := cx - tx
		xx = x + delta
	}

	f.x = xx
	f.y = yy

	f.ResizeChildren()
	f.PlaceChildren()

	return false
}
