package clui

// Checks new control size against its constraints. If a size is less than constraint then the constraint value is used.
// Return adjusted width and height
func ApplyConstraints(c Control, width, height int) (int, int) {
	minw, minh := c.GetConstraints()
	if minw > 0 && width < minw {
		width = minw
	}
	if minh > 0 && height < minh {
		height = minh
	}

	return width, height
}

// Calculates position of the next control using width and heigth of the current one, the last control position, and packer paddings
func CalculateNextPosition(p Packer, w, h int) (int, int, int, int) {
	lastX, lastY := p.GetNextPosition()
	padSide, padTop, padX, padY := p.GetPaddings()
	cW, cH := p.GetConstraints()

	if lastX == -1 || lastY == -1 {
		lastX = padSide
		lastY = padTop
	} else {
		cW -= 2 * padSide
		cH -= 2 * padTop
	}

	if p.GetPackType() == PackHorizontal {
		lastX += w + padX
		if lastX-padX > cW {
			cW = lastX - padX
		}
		if cH < h {
			cH = h
		}
	} else if p.GetPackType() == PackVertical {
		lastY += h + padY
		if lastY-padY-1 > cH {
			cH = lastY - padY - 1
		}
		if cW < w {
			cW = w
		}
	}

	return lastX, lastY, cW + 2*padSide, cH + 2*padTop
}

// Adds a control to a packer into current position, then calculates the next control position depending on the control constraints. Changes the packer constraints after adding the control to keep the packer minimal size enough for all its children
func PackControl(p Packer, c Control, scale int) Control {
	objX, objY := p.GetNextPosition()
	pX, pY := p.GetPos()

	c.SetScale(scale)
	wo, ho := c.GetSize()
	c.SetConstraints(wo, ho)
	c.SetPos(objX+pX, objY+pY)

	objW, objH := c.GetConstraints()
	lastX, lastY, width, height := CalculateNextPosition(p, objW, objH)

	p.SetNextPosition(lastX, lastY)
	p.SetConstraints(width, height)

	return c
}

/*
Helper funtions to make creation of controls easier.
Used for manual/fixed layout
*/

func CreateLabel(parent Window, posX, posY, width int, text string, props Props) *Label {
	id := parent.GetNextControlId()
	lbl := NewLabel(parent, id, posX, posY, width, text, props)
	parent.AddControl(lbl)
	return lbl
}

func CreateEditField(parent Window, posX, posY, width int, text string, props Props) *EditField {
	id := parent.GetNextControlId()
	edit := NewEditField(parent, id, posX, posY, width, text, props)
	parent.AddControl(edit)
	return edit
}

func CreateListBox(parent Window, posX, posY, width, height int, props Props) *ListBox {
	id := parent.GetNextControlId()
	lbox := NewListBox(parent, id, posX, posY, width, height, props)
	parent.AddControl(lbox)
	return lbox
}

func CreateFrame(parent Window, posX, posY, width, height int, title string, props Props) *Frame {
	id := parent.GetNextControlId()
	frm := NewFrame(parent, id, posX, posY, width, height, props)
	if title != "" {
		frm.SetText(title)
	}
	parent.AddControl(frm)
	return frm
}

func CreateButton(parent Window, posX, posY, width, height int, text string, props Props) *Button {
	id := parent.GetNextControlId()
	btn := NewButton(parent, id, posX, posY, width, height, text, props)
	parent.AddControl(btn)
	return btn
}

func CreateCheckbox(parent Window, posX, posY, width, height int, text string, props Props) *CheckBox {
	id := parent.GetNextControlId()
	chk := NewCheckBox(parent, id, posX, posY, width, height, text, props)
	parent.AddControl(chk)
	return chk
}

func CreateRadioGroup(parent Window, posX, posY, width, height int, text string, props Props) *Radio {
	id := parent.GetNextControlId()
	rg := NewRadio(parent, id, posX, posY, width, height, text, props)
	parent.AddControl(rg)
	return rg
}

func CreateComboBox(parent Window, posX, posY, width int, text string, props Props) *EditField {
	id := parent.GetNextControlId()
	cbox := NewComboBox(parent, id, posX, posY, width, text, props)
	parent.AddControl(cbox)
	return cbox
}

func CreateProgressBar(parent Window, posX, posY, width, height, min, max int, props Props) *ProgressBar {
	id := parent.GetNextControlId()
	pb := NewProgressBar(parent, id, posX, posY, width, height, min, max, props)
	parent.AddControl(pb)
	return pb
}

func CreateTextScroll(parent Window, posX, posY, width, height int, props Props) *TextScroll {
	id := parent.GetNextControlId()
	scr := NewTextScroll(parent, id, posX, posY, width, height, props)
	parent.AddControl(scr)
	return scr
}

/*
Helper funtions to make creation of controls easier.
Used for dynamic layout
*/

func PackLabel(parent Window, pack Packer, width int, text string, scale int, props Props) *Label {
	id := parent.GetNextControlId()
	lbl := NewLabel(parent, id, -1, -1, width, text, props)
	parent.AddControl(lbl)
	pack.PackControl(lbl, scale)
	return lbl
}

func PackEditField(parent Window, pack Packer, width int, text string, scale int, props Props) *EditField {
	id := parent.GetNextControlId()
	edit := NewEditField(parent, id, -1, -1, width, text, props)
	parent.AddControl(edit)
	pack.PackControl(edit, scale)
	return edit
}

func PackListBox(parent Window, pack Packer, width, height int, scale int, props Props) *ListBox {
	id := parent.GetNextControlId()
	lbox := NewListBox(parent, id, -1, -1, width, height, props)
	parent.AddControl(lbox)
	pack.PackControl(lbox, scale)
	return lbox
}

func PackFrame(parent Window, pack Packer, width, height int, title string, scale int, props Props) *Frame {
	id := parent.GetNextControlId()
	frm := NewFrame(parent, id, -1, -1, width, height, props)
	if title != "" {
		frm.SetText(title)
	}
	parent.AddControl(frm)
	pack.PackControl(frm, scale)
	return frm
}

func PackButton(parent Window, pack Packer, width, height int, text string, scale int, props Props) *Button {
	id := parent.GetNextControlId()
	btn := NewButton(parent, id, -1, -1, width, height, text, props)
	parent.AddControl(btn)
	pack.PackControl(btn, scale)
	return btn
}

func PackCheckBox(parent Window, pack Packer, width int, text string, scale int, props Props) *CheckBox {
	id := parent.GetNextControlId()
	chk := NewCheckBox(parent, id, -1, -1, width, 1, text, props)
	parent.AddControl(chk)
	pack.PackControl(chk, scale)
	return chk
}

func PackRadioGroup(parent Window, pack Packer, width, height int, text string, scale int, props Props) *Radio {
	id := parent.GetNextControlId()
	rg := NewRadio(parent, id, -1, -1, width, height, text, props)
	parent.AddControl(rg)
	pack.PackControl(rg, scale)
	return rg
}

func PackComboBox(parent Window, pack Packer, width int, text string, scale int, props Props) *EditField {
	id := parent.GetNextControlId()
	cbox := NewComboBox(parent, id, -1, -1, width, text, props)
	parent.AddControl(cbox)
	pack.PackControl(cbox, scale)
	return cbox
}

func PackProgressBar(parent Window, pack Packer, width, height, min, max int, scale int, props Props) *ProgressBar {
	id := parent.GetNextControlId()
	pb := NewProgressBar(parent, id, -1, -1, width, height, min, max, props)
	parent.AddControl(pb)
	pack.PackControl(pb, scale)
	return pb
}

func PackTextScroll(parent Window, pack Packer, width, height int, scale int, props Props) *TextScroll {
	id := parent.GetNextControlId()
	scr := NewTextScroll(parent, id, -1, -1, width, height, props)
	parent.AddControl(scr)
	pack.PackControl(scr, scale)
	return scr
}
