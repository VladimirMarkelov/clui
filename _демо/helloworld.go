package main

import (
	ui ".."
	мИнт "../пакИнтерфейсы"
)

func main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	view := ui.AddWindow(0, 0, 10, 7, "Привет, мир!")

	btnQuit := ui.CreateButton(view, 15, 4, "Привет", 1)
	btnQuit.OnClick(func(ev мИнт.ИСобытие) {
		go ui.Stop()
	})

	ui.MainLoop()
}
