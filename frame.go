package clui

/*
Decorative control - frame with optional title. All area inside a frame is transparent.
Frame can be used as spacer element in dynamic layout - set border to BorderNone and use
that control in any place where a spacer is required
*/
type Frame struct {
	posX, posY    int
	width, height int
	title         string
	anchor        Anchor
	id            WinId
	enabled       bool
	align         Align
	active        bool
	border        BorderStyle
	visible       bool
	textColor     Color
	backColor     Color
	scale         int

	minW, minH int
}

func NewFrame(parent Window, id WinId, x, y, width, height int, props Props) *Frame {
	f := new(Frame)
	f.id = id
	f.SetEnabled(true)
	f.SetPos(x, y)
	f.SetSize(width, height)
	f.border = props.Border
	f.anchor = props.Anchors
	f.visible = true
	f.minW, f.minH = 1, 1

	return f
}

func (f *Frame) SetText(title string) {
	f.title = title
}

func (f *Frame) GetText() string {
	return f.title
}

func (f *Frame) GetId() WinId {
	return f.id
}

func (f *Frame) GetSize() (int, int) {
	return f.width, f.height
}

func (f *Frame) GetConstraints() (int, int) {
	return f.minW, f.minH
}

func (f *Frame) SetConstraints(minW, minH int) {
	if minW >= 1 {
		f.minW = minW
	}
	if minH >= 1 {
		f.minH = minH
	}
}

func (f *Frame) SetSize(width, height int) {
	width, height = ApplyConstraints(f, width, height)
	f.width = width
	f.height = height
}

func (f *Frame) GetPos() (int, int) {
	return f.posX, f.posY
}

func (f *Frame) SetPos(x, y int) {
	f.posX = x
	f.posY = y
}

func (f *Frame) Redraw(canvas Canvas) {
	tm := canvas.Theme()

	x, y := f.GetPos()
	w, h := f.GetSize()

	fg, bg := f.textColor, f.backColor

	if fg == ColorDefault {
		fg = tm.GetSysColor(ColorActiveText)
	}
	if bg == ColorDefault {
		bg = tm.GetSysColor(ColorViewBack)
	}

	canvas.DrawFrame(x, y, w, h, f.border, fg, bg)

	if f.title != "" {
		text := Ellipsize(f.title, w-2)
		canvas.DrawText(x+1, y, w-2, text, fg, bg)
	}
}

func (f *Frame) GetEnabled() bool {
	return f.enabled
}

func (f *Frame) SetEnabled(active bool) {
	// nothing to do
}

func (f *Frame) SetAlign(align Align) {
	// nothing
}

func (f *Frame) GetAlign() Align {
	return f.align
}

func (f *Frame) SetAnchors(anchor Anchor) {
	f.anchor = anchor
}

func (f *Frame) GetAnchors() Anchor {
	return f.anchor
}

func (f *Frame) GetActive() bool {
	return f.active
}

func (f *Frame) SetActive(active bool) {
	f.active = active
}

func (f *Frame) GetTabStop() bool {
	return false
}

func (f *Frame) SetTabStop(tab bool) {
	// nothing
}

func (f *Frame) ProcessEvent(event Event) bool {
	return false
}

func (f *Frame) SetVisible(visible bool) {
	f.visible = visible
}

func (f *Frame) GetVisible() bool {
	return f.visible
}

func (f *Frame) GetColors() (Color, Color) {
	return f.textColor, f.backColor
}

func (f *Frame) SetTextColor(clr Color) {
	f.textColor = clr
}

func (f *Frame) SetBackColor(clr Color) {
	f.backColor = clr
}

func (f *Frame) HideChildren() {
	// nothing to do
}

func (f *Frame) GetScale() int {
	return f.scale
}

func (f *Frame) SetScale(scale int) {
	f.scale = scale
}
