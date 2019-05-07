package main

import (
	ui "../.."
	мИнт "../../пакИнтерфейсы"
)

func main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	view := ui.AddWindow(0, 0, 10, 7, "Hello World!")
	view.SetPack(мИнт.Vertical)

	frmResize := ui.CreateFrame(view, 8, 6, ui.BorderNone, мИнт.Fixed)
	frmResize.SetTitle("FrameTop")
	frmResize.SetPack(мИнт.Horizontal)
	btn1 := ui.CreateButton(frmResize, 8, 5, "Button 1", 1)
	btn2 := ui.CreateButton(frmResize, 8, 5, "Button 2", 1)
	btn3 := ui.CreateButton(frmResize, 8, 5, "Button 3", 1)

	frmBtns := ui.CreateFrame(view, 8, 5, ui.BorderNone, ui.Fixed)
	frmBtns.SetPack(ui.Horizontal)
	frmBtns.SetTitle("FrameBottom")

	btnHide1 := ui.CreateButton(frmBtns, 8, 4, "Hide 1", 1)
	btnHide1.OnClick(func(ev ui.Event) {
		if btn1.Visible() {
			btnHide1.SetTitle("Show 1")
			ui.ActivateControl(view, btn1)
			btn1.SetVisible(false)
		} else {
			btnHide1.SetTitle("Hide 1")
			btn1.SetVisible(true)
		}
	})
	btnHide2 := ui.CreateButton(frmBtns, 8, 4, "Hide 2", 1)
	btnHide2.OnClick(func(ev ui.Event) {
		if btn2.Visible() {
			btnHide2.SetTitle("Show 2")
			ui.ActivateControl(view, btn2)
			btn2.SetVisible(false)
		} else {
			btnHide2.SetTitle("Hide 2")
			btn2.SetVisible(true)
		}
	})
	btnHide3 := ui.CreateButton(frmBtns, 8, 4, "Hide 3", 1)
	btnHide3.OnClick(func(ev ui.Event) {
		if btn3.Visible() {
			btnHide3.SetTitle("Show 3")
			ui.ActivateControl(view, btn3)
			btn3.SetVisible(false)
		} else {
			btnHide3.SetTitle("Hide 3")
			btn3.SetVisible(true)
		}
	})

	btnQuit := ui.CreateButton(frmBtns, 8, 4, "Quit", 1)
	btnQuit.OnClick(func(ev ui.Event) {
		go ui.Stop()
	})

	ui.MainLoop()
}
