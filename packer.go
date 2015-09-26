package clui

/*
Dynamic control layout support for a Window. Each Window can turn on
its own layout mode: dynamic or manual. Layout mode must be selected
when a Window is created - one cannot change layout mode if a Window
has any control or Packer already.

A Packer can be horizontal or vertical. It sets the direction in which
the Packer adds every new child. E.g, in horizontal mode a Packer has
1-Contol height(real height is the height of the biggest control) and
unlimited width - width is increased after a new control is added.

To start dynamic mode, call AddPacker method of Window - it creates
and returns a main Window's Packer - a Window can have only one Packer.
Every Packer implements AddPacker method to create inner containers.
It allows to build a complex dialogs. Do not forget to call Window
method PackEnd when adding Controls and Packers is done. PackEnd
completes layouting process, calculates Packer's children sizes and
positions. It is not possible to add a new child with dynamic layout,
only a child with manual defined position and size can be added after
calling PackEnd.

Packer is a container that holds either Packers or controls - it is not
possible to mix controls and Packers in one container. All child
object are automatically placed inside a Packer and while application is
running the packer manages its children size and positions.

Set up padding sizes before adding new children. By default left/right
and top/bottom indents, vertical and horizontal paddings are 0. If you
turn on Packer border and indents equal 0 then indents are automatically
changed to 1 to avoid drawing controls on the Packer frame.
Turn on Packer border to imitate Frame control.

To add a child to Packer, use either AddPacker for new container or
PackControl(or shorthand methods like PackLabel etc) for new Control.

One should not use other methods of Packer like SetConstraints, SetSize
directly - these methods are used by Composer and parent Window to
react to its parent Window resize.
*/
type Container struct {
	posX, posY    int
	width, height int
	title         string
	id            WinId
	border        BorderStyle
	textColor     Color
	backColor     Color

	minW, minH int

	// pack support
	pack            PackType
	parent          Packer
	controls        []WinId
	view            Window
	padSide, padTop int
	padX, padY      int
	scale           int
	children        []WinId
	packers         []Packer
	lastX, lastY    int
}

func NewContainer(view Window, parent Packer, id WinId, x, y, width, height int, props Props) *Container {
	if view == nil {
		panic("Packer view cannot be nil")
	}
	f := new(Container)
	f.id = id
	f.SetPos(x, y)
	f.SetSize(width, height)
	f.border = props.Border
	f.parent = parent
	f.packers = make([]Packer, 0)
	f.view = view

	f.minW, f.minH = 1, 1
	f.padTop, f.padSide, f.padX, f.padY = 0, 0, 0, 0

	if f.border != BorderNone {
		f.padTop++
		f.padSide++
	}

	f.lastX, f.lastY = -1, -1

	return f
}

func (f *Container) SetText(title string) {
	f.title = title
}

func (f *Container) GetText() string {
	return f.title
}

func (f *Container) GetId() WinId {
	return f.id
}

// Internal method used by PackEnd method.
// The Packer and its inner containers calculate their
// sizes depending on their children constraints and
// eventually sets all the Packers constraints
func (f *Container) CalculateSize() (int, int) {
	if len(f.packers) == 0 {
		f.SetConstraints(f.width, f.height)
		return f.width, f.height
	} else {
		w, h := 0, 0
		for _, p := range f.packers {
			wp, hp := p.CalculateSize()
			if f.pack == PackHorizontal {
				w += wp
				if h < hp {
					h = hp
				}
			} else {
				h += hp
				if w < wp {
					w = wp
				}
			}
		}
		if f.border != BorderNone {
			w += 2
			h += 2
		}
		f.SetSize(w, h)
		f.SetConstraints(w, h)
		return w, h
	}
}

func (f *Container) GetSize() (int, int) {
	return f.width, f.height
}

func (f *Container) GetConstraints() (int, int) {
	return f.minW, f.minH
}

func (f *Container) SetConstraints(minW, minH int) {
	if minW >= 1 {
		f.minW = minW
	}
	if minH >= 1 {
		f.minH = minH
	}

	if f.width < minW || f.height < minH {
		f.SetSize(minW, minH)
	}
}

func (f *Container) SetSize(width, height int) {
	f.width = width
	f.height = height
}

func (f *Container) GetPos() (int, int) {
	return f.posX, f.posY
}

func (f *Container) SetPos(x, y int) {
	f.posX = x
	f.posY = y
}

func (f *Container) Redraw(canvas Canvas) {
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

	if f.packers != nil && len(f.packers) > 0 {
		for _, pck := range f.packers {
			pck.Redraw(canvas)
		}
	}
}

func (f *Container) ProcessEvent(event Event) bool {
	return false
}

func (f *Container) GetColors() (Color, Color) {
	return f.textColor, f.backColor
}

func (f *Container) SetTextColor(clr Color) {
	f.textColor = clr
}

func (f *Container) SetBackColor(clr Color) {
	f.backColor = clr
}

func (f *Container) GetContainer() Packer {
	return f.parent
}

func (f *Container) repositionPackers() {
	x, y := f.GetPos()
	pSide, pTop, pX, pY := f.GetPaddings()
	x += pSide
	y += pTop
	pType := f.GetPackType()

	for _, pck := range f.packers {
		pck.SetPos(x, y)
		w, h := pck.GetSize()
		if pType == PackHorizontal {
			x += pX + w
		} else {
			y += pY + h
		}

		pck.RepositionChildren()
	}
}

func (f *Container) repositionControls() {
	x, y := f.GetPos()
	pSide, pTop, pX, pY := f.GetPaddings()
	x += pSide
	y += pTop
	pType := f.GetPackType()

	w, h := f.GetSize()
	minW, minH := f.GetConstraints()
	dx, dy := w-minW, h-minH
	w -= 2 * pSide
	h -= 2 * pTop

	pScale, scaled := f.getScaleStats()
	left := dx
	if pType == PackVertical {
		left = dy
	}

	for _, ctrlId := range f.children {
		ctrl := f.view.GetControl(ctrlId)
		ctrl.SetPos(x, y)
		cw, ch := ctrl.GetConstraints()
		sc := ctrl.GetScale()

		if pType == PackHorizontal {
			ch = h
		} else {
			cw = w
		}

		if sc <= DoNotScale {
			if pType == PackHorizontal {
				x += pX + cw
			} else {
				y += pY + ch
			}
		} else {
			delta := 0
			if pType == PackHorizontal {
				if scaled == 1 {
					delta = left
				} else {
					delta = dx * sc / pScale
				}
				cw += delta
				x += pX + cw
			} else {
				if scaled == 1 {
					delta = left
				} else {
					delta = dy * sc / pScale
				}
				ch += delta
				y += pY + ch
			}
			left -= delta
			scaled--
		}

		ctrl.SetSize(cw, ch)
	}
}

func (f *Container) RepositionChildren() {
	if len(f.packers) != 0 {
		f.repositionPackers()
	} else {
		f.repositionControls()
	}
}

func (f *Container) GetNextPosition() (int, int) {
	if f.lastY < 0 {
		return f.padSide, f.padTop
	}
	return f.lastX, f.lastY
}

func (f *Container) SetNextPosition(x, y int) {
	f.lastX, f.lastY = x, y
}

func (f *Container) GetPackType() PackType {
	return f.pack
}

func (f *Container) SetPaddings(pSide, pTop, pX, pY int) {
	if len(f.children) > 0 {
		panic("Cannot change padding if a child is added")
	}

	if pSide != DoNotChange {
		f.padSide = pSide
	}
	if pTop != DoNotChange {
		f.padTop = pTop
	}
	if pX != DoNotChange {
		f.padX = pX
	}
	if pY != DoNotChange {
		f.padY = pSide
	}
}

func (f *Container) GetPaddings() (int, int, int, int) {
	return f.padSide, f.padTop, f.padX, f.padY
}

func (f *Container) getScaleStats() (int, int) {
	dScale, scaled := 0, 0

	if len(f.packers) > 0 {
		for _, pck := range f.packers {
			sc := pck.GetScale()
			if sc > DoNotScale {
				dScale += sc
				scaled++
			}
		}
	} else if len(f.children) > 0 {
		for _, childId := range f.children {
			child := f.view.GetControl(childId)
			sc := child.GetScale()
			if sc > DoNotScale {
				dScale += sc
				scaled++
			}
		}
	}

	return dScale, scaled
}

func (f *Container) resizePackers(dx, dy int) {
	pSide, pTop, _, _ := f.GetPaddings()
	w, h := f.GetConstraints()
	h = h + dy - 2*pTop
	w = w + dx - 2*pSide
	origDx, origDy := dx, dy
	dScale, scaled := f.getScaleStats()
	pType := f.GetPackType()

	for _, pck := range f.packers {
		pScale := pck.GetScale()

		pW, pH := pck.GetConstraints()
		if pScale <= DoNotScale {
			if pType == PackHorizontal {
				pck.SetSize(pW, h)
			} else {
				pck.SetSize(w, pH)
			}
			continue
		}

		dpX, dpY := pScale*origDx/dScale, pScale*origDy/dScale

		if pType == PackHorizontal {
			pH = h
			if scaled == 1 {
				pW += dx
			} else {
				pW += dpX
				dx -= dpX
			}
		} else {
			pW = w
			if scaled == 1 {
				pH += dy
			} else {
				pH += dpY
				dy -= dpY
			}
		}
		pck.SetSize(pW, pH)
		pck.ResizeChidren(dx, dy)
		scaled--
	}
}

func (f *Container) ResizeChidren(dx, dy int) {
	if len(f.packers) != 0 {
		f.resizePackers(dx, dy)
	}
}

func (f *Container) GetScale() int {
	return f.scale
}

func (f *Container) SetScale(scale int) {
	f.scale = scale
}

func (f *Container) SetPackType(pt PackType) {
	if len(f.controls) > 0 {
		panic("Cannot enable pack mode if a packer contains any control")
	}

	if pt == PackFixed {
		panic("Useless pack type Fixed - use manual layout of view instead")
	}

	f.pack = pt
}

func (f *Container) AddPack(pt PackType, scale int) Packer {
	if len(f.children) > 0 {
		panic("Cannot add a packer - the packer already contains controls")
	}

	if pt == PackFixed {
		panic("Fixed layout is not allowed for inner packers")
	}

	pid := f.view.GetNextControlId()
	x, y := f.GetNextPosition()
	p := NewContainer(f.view, f, pid, x, y, AutoSize, AutoSize, Props{})
	p.SetPackType(pt)
	p.SetScale(scale)
	f.packers = append(f.packers, p)

	return p
}

func (f *Container) PackControl(c Control, scale int) Control {
	if len(f.packers) > 0 {
		panic("Cannot add a control - the packer already contains packers")
	}

	PackControl(f, c, scale)
	f.children = append(f.children, c.GetId())
	return c
}

// ---- syntax sugar for control packing ---------------

func (f *Container) PackLabel(width int, text string, scale int, props Props) *Label {
	id := f.view.GetNextControlId()
	lbl := NewLabel(f.view, id, -1, -1, width, text, props)
	f.view.AddControl(lbl)

	if len(f.packers) > 0 {
		pack := f.AddPack(f.GetPackType(), 1)
		pack.PackControl(lbl, scale)
	} else {
		f.PackControl(lbl, scale)
	}

	return lbl
}

func (f *Container) PackEditField(width int, text string, scale int, props Props) *EditField {
	id := f.view.GetNextControlId()
	edit := NewEditField(f.view, id, -1, -1, width, text, props)
	f.view.AddControl(edit)

	if len(f.packers) > 0 {
		pack := f.AddPack(f.GetPackType(), 1)
		pack.PackControl(edit, scale)
	} else {
		f.PackControl(edit, scale)
	}

	return edit
}

func (f *Container) PackListBox(width, height int, scale int, props Props) *ListBox {
	id := f.view.GetNextControlId()
	lbox := NewListBox(f.view, id, -1, -1, width, height, props)
	f.view.AddControl(lbox)

	if len(f.packers) > 0 {
		pack := f.AddPack(f.GetPackType(), 1)
		pack.PackControl(lbox, scale)
	} else {
		f.PackControl(lbox, scale)
	}

	return lbox
}

func (f *Container) PackFrame(width, height int, title string, scale int, props Props) *Frame {
	id := f.view.GetNextControlId()
	frm := NewFrame(f.view, id, -1, -1, width, height, props)
	if title != "" {
		frm.SetText(title)
	}
	f.view.AddControl(frm)

	if len(f.packers) > 0 {
		pack := f.AddPack(f.GetPackType(), 1)
		pack.PackControl(frm, scale)
	} else {
		f.PackControl(frm, scale)
	}

	return frm
}

func (f *Container) PackButton(width, height int, text string, scale int, props Props) *Button {
	id := f.view.GetNextControlId()
	btn := NewButton(f.view, id, -1, -1, width, height, text, props)
	f.view.AddControl(btn)

	if len(f.packers) > 0 {
		pack := f.AddPack(f.GetPackType(), 1)
		pack.PackControl(btn, scale)
	} else {
		f.PackControl(btn, scale)
	}
	return btn
}

func (f *Container) PackCheckBox(width int, text string, scale int, props Props) *CheckBox {
	id := f.view.GetNextControlId()
	chk := NewCheckBox(f.view, id, -1, -1, width, 1, text, props)
	f.view.AddControl(chk)

	if len(f.packers) > 0 {
		pack := f.AddPack(f.GetPackType(), 1)
		pack.PackControl(chk, scale)
	} else {
		f.PackControl(chk, scale)
	}
	return chk
}

func (f *Container) PackRadioGroup(width, height int, text string, scale int, props Props) *Radio {
	id := f.view.GetNextControlId()
	rg := NewRadio(f.view, id, -1, -1, width, height, text, props)
	f.view.AddControl(rg)

	if len(f.packers) > 0 {
		pack := f.AddPack(f.GetPackType(), 1)
		pack.PackControl(rg, scale)
	} else {
		f.PackControl(rg, scale)
	}

	return rg
}

func (f *Container) PackComboBox(width int, text string, scale int, props Props) *EditField {
	id := f.view.GetNextControlId()
	cbox := NewComboBox(f.view, id, -1, -1, width, text, props)
	f.view.AddControl(cbox)

	if len(f.packers) > 0 {
		pack := f.AddPack(f.GetPackType(), 1)
		pack.PackControl(cbox, scale)
	} else {
		f.PackControl(cbox, scale)
	}

	return cbox
}

func (f *Container) PackProgressBar(width, height, min, max int, scale int, props Props) *ProgressBar {
	id := f.view.GetNextControlId()
	pb := NewProgressBar(f.view, id, -1, -1, width, height, min, max, props)
	f.view.AddControl(pb)

	if len(f.packers) > 0 {
		pack := f.AddPack(f.GetPackType(), 1)
		pack.PackControl(pb, scale)
	} else {
		f.PackControl(pb, scale)
	}

	return pb
}

func (f *Container) PackTextScroll(width, height int, scale int, props Props) *TextScroll {
	id := f.view.GetNextControlId()
	scr := NewTextScroll(f.view, id, -1, -1, width, height, props)
	f.view.AddControl(scr)

	if len(f.packers) > 0 {
		pack := f.AddPack(f.GetPackType(), 1)
		pack.PackControl(scr, scale)
	} else {
		f.PackControl(scr, scale)
	}

	return scr
}

func (f *Container) GetBorderStyle() BorderStyle {
	return f.border
}

func (f *Container) SetBorderStyle(bs BorderStyle) {
	f.border = bs
	if len(f.children) == 0 && len(f.packers) == 0 {
		// fix paddings if there are no children yet
		if f.padSide == 0 {
			f.padSide = 1
		}
		if f.padTop == 0 {
			f.padTop = 1
		}
	}
}

func (f *Container) View() Window {
	return f.view
}
