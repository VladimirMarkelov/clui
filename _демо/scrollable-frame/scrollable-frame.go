package main

import (
	"fmt"
	ui "../.."
	мИнт "../../пакИнтерфейсы"
)

func main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	wnd := ui.AddWindow(0, 0, 60, мИнт.AutoSize, "Фрейм со скроллом")
	wnd.SetSizable(false)

	frm := ui.CreateFrame(wnd, 50, 12, мИнт.BorderNone, мИнт.Fixed)
	frm.SetPack(мИнт.Vertical)
	frm.SetScrollable(true)

	for i := 0; i < 10; i++ {
		label := fmt.Sprintf("Кнопка %d - нажимте для выхода", i)
		btn := ui.CreateButton(frm, 40, мИнт.AutoSize, label, 1)

		btn.OnClick(func(ev мИнт.ИСобытие) {
			go ui.Stop()
		})
	}

	ui.MainLoop()
}
