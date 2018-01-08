package main

import (
	ui "github.com/VladimirMarkelov/clui"
)

func createView() {
	view := ui.AddWindow(0, 0, 10, 7, "EditField Demo")

	frmChk := ui.CreateFrame(view, 8, 5, ui.BorderNone, ui.Fixed)
	frmChk.SetPack(ui.Vertical)
    frmChk.SetPaddings(1, 1)
    frmChk.SetGaps(1, 1)
    ui.CreateLabel(frmChk, ui.AutoSize, ui.AutoSize, "Enter password:", ui.Fixed)
	edFld := ui.CreateEditField(frmChk, 20, "", ui.Fixed)
    edFld.SetPasswordMode(true)
	chkPass := ui.CreateCheckBox(frmChk, ui.AutoSize, "Show Password", ui.Fixed)

	ui.ActivateControl(view, edFld)

	chkPass.OnChange(func(state int) {
		if state == 1 {
			edFld.SetPasswordMode(false)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 0 {
			edFld.SetPasswordMode(true)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	createView()

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
