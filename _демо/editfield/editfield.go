package main

import (
	ui "../.."
	мИнт "../../пакИнтерфейсы"
	мСоб "../../пакСобытия"
)

func createView() {
	view := ui.AddWindow(0, 0, 10, 7, "Пример редактора")
	view.SetTitleButtons(мИнт.ButtonMaximize | мИнт.ButtonClose)

	frmChk := ui.CreateFrame(view, 8, 5, мИнт.BorderNone, мИнт.Fixed)
	frmChk.SetPack(мИнт.Vertical)
	frmChk.SetPaddings(1, 1)
	frmChk.SetGaps(1, 1)
	ui.CreateLabel(frmChk, мИнт.AutoSize, мИнт.AutoSize, "Введите пароль:", мИнт.Fixed)
	edFld := ui.CreateEditField(frmChk, 20, "", мИнт.Fixed)
	edFld.SetPasswordMode(true)
	chkPass := ui.CreateCheckBox(frmChk, мИнт.AutoSize, "Показать пароль", мИнт.Fixed)

	ui.ActivateControl(view, edFld)

	chkPass.OnChange(func(state int) {
		if state == 1 {
			edFld.SetPasswordMode(false)
			ev := &мСоб.Event{}
			ev.TypeSet(мИнт.EventRedraw)
			ui.PutEvent(ev)
		} else if state == 0 {
			edFld.SetPasswordMode(true)
			ev := &мСоб.Event{}
			ev.TypeSet(мИнт.EventRedraw)
			ui.PutEvent(ev)
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
