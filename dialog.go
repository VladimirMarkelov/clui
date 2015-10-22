package clui

import (
	term "github.com/nsf/termbox-go"
)

type OkCancelDialog struct {
	view    View
	parent  *Composer
	result  int
	onClose func()
}

func NewOkCancelDialog(c *Composer, title, question string, buttonTitles []string) *OkCancelDialog {
	dlg := new(OkCancelDialog)

	cw, ch := term.Size()

	dlg.parent = c
	dlg.view = c.CreateView(cw/2-12, ch/2-8, 20, 10, title)
	dlg.view.SetModal(true)
	dlg.view.SetPack(Vertical)
	NewFrame(dlg.view, dlg.view, 1, 1, BorderNone, DoNotScale)

	fbtn := NewFrame(dlg.view, dlg.view, 1, 1, BorderNone, 1)
	NewFrame(dlg.view, fbtn, 1, 1, BorderNone, DoNotScale)
	lb := NewLabel(dlg.view, fbtn, 10, 3, question, 1)
	NewFrame(dlg.view, fbtn, 1, 1, BorderNone, DoNotScale)
	lb.SetMultiline(true)

	NewFrame(dlg.view, dlg.view, 1, 1, BorderNone, DoNotScale)
	frm1 := NewFrame(dlg.view, dlg.view, 16, 4, BorderNone, DoNotScale)
	NewFrame(dlg.view, frm1, 1, 1, BorderNone, 1)
	bText := "OK"
	if len(buttonTitles) > 0 {
		bText = buttonTitles[0]
	}
	btnOk := NewButton(dlg.view, frm1, AutoSize, AutoSize, bText, DoNotScale)
	if len(buttonTitles) > 1 {
		bText = buttonTitles[1]
	} else {
		bText = "Cancel"
	}
	NewFrame(dlg.view, frm1, 1, 1, BorderNone, 1)
	btnCancel := NewButton(dlg.view, frm1, AutoSize, AutoSize, bText, DoNotScale)
	NewFrame(dlg.view, frm1, 1, 1, BorderNone, 1)
	dlg.view.ActivateControl(btnCancel)

	dlg.view.OnClose(func(ev Event) {
		if dlg.result == DialogAlive {
			dlg.result = DialogClosed
			c.DestroyView(dlg.view)
			if dlg.onClose != nil {
				go dlg.onClose()
			}
		}
	})
	btnOk.OnClick(func(ev Event) {
		dlg.result = DialogOK
		c.DestroyView(dlg.view)
		if dlg.onClose != nil {
			go dlg.onClose()
		}
	})
	btnCancel.OnClick(func(ev Event) {
		dlg.result = DialogCancel
		c.DestroyView(dlg.view)
		if dlg.onClose != nil {
			go dlg.onClose()
		}
	})

	return dlg
}

func (d *OkCancelDialog) OnClose(fn func()) {
	d.onClose = fn
}

func (d *OkCancelDialog) Result() int {
	return d.result
}
