package clui

/*
Frame is a decorative control and container - frame with optional title.
All area inside a frame is transparent. Frame can be used as spacer element
- set border to BorderNone and use that control in any place where a spacer
is required
*/
type Frame struct {
	ControlBase
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
func NewFrame(view View, parent Control, width, height int, bs BorderStyle, scale int) *Frame {
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
	f.view = view
	f.parent = parent
	f.SetTabStop(false)

	f.padX, f.padY = 0, 0
	if bs != BorderNone {
		f.padSide, f.padTop = 1, 1
	} else {
		f.padSide, f.padTop = 0, 0
	}

	if parent != nil {
		parent.AddChild(f, scale)
	}

	return f
}

func (f *Frame) repaintChildren() {
	for _, ctrl := range f.children {
		ctrl.Repaint()
	}
}

// Repaint draws the control on its View surface
func (f *Frame) Repaint() {
	f.repaintChildren()

	if f.border == BorderNone {
		return
	}

	tm := f.view.Screen().Theme()
	canvas := f.view.Canvas()

	x, y := f.Pos()
	w, h := f.Size()

	fg, bg := RealColor(tm, f.fg, ColorViewText), RealColor(tm, f.bg, ColorViewBack)
	var chars string
	if f.border == BorderDouble {
		chars = tm.SysObject(ObjDoubleBorder)
	} else {
		chars = tm.SysObject(ObjSingleBorder)
	}

	canvas.DrawFrame(x, y, w, h, fg, bg, chars)

	if f.title != "" {
		text := Ellipsize(f.title, w-2)
		canvas.PutText(x+1, y, text, fg, bg)
	}
}

// RecalculateConstraints used by containers to recalculate new minimal size
// depending on its children constraints after a new child is added
func (f *Frame) RecalculateConstraints() {
	width, height := f.Constraints()
	minW, minH := CalculateMinimalSize(f)

	newW, newH := width, height
	if minW > newW {
		newW = minW
	}
	if minH > newH {
		newH = minH
	}

	if newW != width || newH != height {
		f.SetConstraints(newW, newH)
	}

	if f.parent != nil {
		f.parent.RecalculateConstraints()
	} else if f.view != nil {
		f.view.RecalculateConstraints()
	}
}

// AddChild adds a new child to a container. After adding
// a new child the frame automatically recalculates its
// its minimal size
func (f *Frame) AddChild(c Control, scale int) {
	if f.view.ChildExists(c) {
		// Frame: Cannot add the same control twice
		return
	}

	c.SetScale(scale)
	f.children = append(f.children, c)
	f.RecalculateConstraints()
	f.view.RegisterControl(c)
}

// Children returns the list of container child controls
func (f *Frame) Children() []Control {
	return f.children
}

// SetPack changes the direction of children packing.
// Changing pack type on the fly is not always possible:
// it does nothing if a frame already contains children
func (f *Frame) SetPack(pk PackType) {
	if len(f.children) > 0 {
		// Control already has children
		return
	}

	f.pack = pk
}

// Pack returns direction in which a container packs
// its children: horizontal or vertical
func (f *Frame) Pack() PackType {
	return f.pack
}
