package clui

// Decorative control to display horizontal or vertical text line. Text can be aligned within Label bounding box. Multiline output is not supported
type Label struct {
	posX, posY    int
	width, height int
	title         string
	anchor        Anchor
	id            WinId
	enabled       bool
	align         Align
	active        bool
	direction     Direction
	visible       bool
	textColor     Color
	backColor     Color
	scale         int

	minW, minH int
}

func NewLabel(parent Window, id WinId, x, y, length int, text string, props Props) *Label {
	l := new(Label)
	l.SetText(text)
	l.SetEnabled(true)
	l.SetPos(x, y)
	l.direction = props.Dir
	l.id = id
	if props.Dir == DirHorizontal {
		l.SetSize(length, 1)
	} else {
		l.SetSize(1, length)
	}
	l.visible = true
	l.minW, l.minH = 1, 1

	return l
}

func (l *Label) SetText(title string) {
	l.title = title
}

func (l *Label) GetText() string {
	return l.title
}

func (l *Label) GetId() WinId {
	return l.id
}

func (l *Label) GetSize() (int, int) {
	return l.width, l.height
}

func (l *Label) GetConstraints() (int, int) {
	return l.minW, l.minH
}

func (l *Label) SetConstraints(minW, minH int) {
	if minW >= 1 {
		l.minW = minW
	}
	if minH >= 1 {
		l.minH = minH
	}
}

func (l *Label) SetSize(width, height int) {
	width, height = ApplyConstraints(l, width, height)
	l.width = width
	l.height = height
}

func (l *Label) GetPos() (int, int) {
	return l.posX, l.posY
}

func (l *Label) SetPos(x, y int) {
	l.posX = x
	l.posY = y
}

func (l *Label) Redraw(canvas Canvas) {
	x, y := l.GetPos()
	w, h := l.GetSize()

	tm := canvas.Theme()

	fg, bg := l.textColor, l.backColor
	if fg == ColorDefault {
		if l.enabled {
			fg = tm.GetSysColor(ColorControlText)
		} else {
			fg = tm.GetSysColor(ColorGrayText)
		}
	}
	if bg == ColorDefault {
		bg = tm.GetSysColor(ColorViewBack)
	}

	if l.direction == DirHorizontal {
		canvas.ClearRect(x, y, w, 1, bg)
		canvas.DrawText(x, y, w, l.GetText(), fg, bg)
	} else {
		canvas.ClearRect(x, y, 1, h, bg)
		canvas.DrawVerticalText(x, y, h, l.GetText(), fg, bg)
	}
}

func (l *Label) GetEnabled() bool {
	return l.enabled
}

func (l *Label) SetEnabled(active bool) {
	l.enabled = active
}

func (l *Label) SetAlign(align Align) {
	l.align = align
}

func (l *Label) GetAlign() Align {
	return l.align
}

func (l *Label) SetAnchors(anchor Anchor) {
	l.anchor = anchor
}

func (l *Label) GetAnchors() Anchor {
	return l.anchor
}

func (l *Label) GetActive() bool {
	return l.active
}

func (l *Label) SetActive(active bool) {
	l.active = active
}

func (l *Label) GetTabStop() bool {
	return false
}

func (l *Label) SetTabStop(tab bool) {
	// nothing
}

func (l *Label) ProcessEvent(event Event) bool {
	return false
}

func (l *Label) SetVisible(visible bool) {
	l.visible = visible
}

func (l *Label) GetVisible() bool {
	return l.visible
}

func (l *Label) GetColors() (Color, Color) {
	return l.textColor, l.backColor
}

func (l *Label) SetTextColor(clr Color) {
	l.textColor = clr
}

func (l *Label) SetBackColor(clr Color) {
	l.backColor = clr
}

func (l *Label) HideChildren() {
	// nothing to do
}

func (l *Label) GetScale() int {
	return l.scale
}

func (l *Label) SetScale(scale int) {
	l.scale = scale
}
