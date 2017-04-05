package main

import (
	ui "github.com/VladimirMarkelov/clui"
)

func main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	view := ui.AddWindow(0, 0, 10, 7, "Hello World!")

	btnQuit := ui.CreateButton(view, 15, 4, "Hi", 1)
	btnQuit.OnClick(func(ev ui.Event) {
		go ui.Stop()
	})

	ui.MainLoop()
}
