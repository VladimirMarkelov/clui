package clui

/*
Decorative control - frame with optional title. All area inside a frame is transparent.
Frame can be used as spacer element in dynamic layout - set border to BorderNone and use
that control in any place where a spacer is required
*/
type Frame struct {
	ControlBase
	border   BorderStyle
	children []Control
	pack     PackType
}

func NewFrame(view View, parent Control, width, height int, bs BorderStyle, scale int) *Frame {
	f := new(Frame)
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
	var chars string = ""
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

func (f *Frame) AddChild(c Control, scale int) {
	if f.view.ChildExists(c) {
		panic("Frame: Cannot add the same control twice")
	}

	c.SetScale(scale)
	f.children = append(f.children, c)
	f.RecalculateConstraints()
	f.view.RegisterControl(c)
}

func (f *Frame) Children() []Control {
	return f.children
}

func (f *Frame) SetPack(pk PackType) {
	if len(f.children) > 0 {
		panic("Control already has children")
	}

	f.pack = pk
}

func (f *Frame) Pack() PackType {
	return f.pack
}
