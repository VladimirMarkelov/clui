package main

import (
	ui "../.."
	мИнт "../../пакИнтерфейсы"
)

func main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	view := ui.AddWindow(0, 0, 10, 7, "Привет, мир!")
	view.SetPack(мИнт.Vertical)

	frmResize := ui.CreateFrame(view, 8, 6, мИнт.BorderNone, мИнт.Fixed)
	frmResize.SetTitle("FrameTop")
	frmResize.SetPack(мИнт.Horizontal)
	btn1 := ui.CreateButton(frmResize, 8, 5, "Кнопка 1", 1)
	btn2 := ui.CreateButton(frmResize, 8, 5, "Кнопка 2", 1)
	btn3 := ui.CreateButton(frmResize, 8, 5, "Кнопка 3", 1)

	frmBtns := ui.CreateFrame(view, 8, 5, мИнт.BorderNone, мИнт.Fixed)
	frmBtns.SetPack(мИнт.Horizontal)
	frmBtns.SetTitle("FrameBottom")

	btnHide1 := ui.CreateButton(frmBtns, 8, 4, "Скрыть Кн1", 1)
	btnHide1.OnClick(func(ev мИнт.ИСобытие) {
		if btn1.Visible() {
			btnHide1.SetTitle("Показать Кн1")
			ui.ActivateControl(view, btn1)
			btn1.SetVisible(false)
		} else {
			btnHide1.SetTitle("Скрыть Кн1")
			btn1.SetVisible(true)
		}
	})
	btnHide2 := ui.CreateButton(frmBtns, 8, 4, "Скрыть Кн2", 1)
	btnHide2.OnClick(func(ev мИнт.ИСобытие) {
		if btn2.Visible() {
			btnHide2.SetTitle("Показать Кн2")
			ui.ActivateControl(view, btn2)
			btn2.SetVisible(false)
		} else {
			btnHide2.SetTitle("Скрыть Кн2")
			btn2.SetVisible(true)
		}
	})
	btnHide3 := ui.CreateButton(frmBtns, 8, 4, "Скрыть Кн3", 1)
	btnHide3.OnClick(func(ev мИнт.ИСобытие) {
		if btn3.Visible() {
			btnHide3.SetTitle("Показать Кн3")
			ui.ActivateControl(view, btn3)
			btn3.SetVisible(false)
		} else {
			btnHide3.SetTitle("Скрыть Кн3")
			btn3.SetVisible(true)
		}
	})

	btnQuit := ui.CreateButton(frmBtns, 8, 4, "Выход", 1)
	btnQuit.OnClick(func(ev мИнт.ИСобытие) {
		go ui.Stop()
	})

	ui.MainLoop()
}
