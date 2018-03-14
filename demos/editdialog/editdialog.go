package main

import (
	ui "github.com/Viv1k/clui"

	termbox "github.com/nsf/termbox-go"
)

func CreateBox() {
	dlg := ui.CreateConfirmationEditDialog(
		"<c:blue>"+"SearchBox:",
		"Enter search text. Pressing enter would print your input in debug.log")

	dlg.OnClose(func() {
		// write input test to debug.log only when enter is pressed
		if dlg.Result() == ui.DialogButton1 {
			ui.Logger().Println("result", dlg.EditResult())
		}
	})
}

func mainLoop() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	// BUG, if you don't provide any window and create a dialog box
	// app crashes on pressing any button of dialog
	window := ui.AddWindow(0, 0, 10, 7, "EditDialog Demo")
	window.OnKeyDown(func(event ui.Event) bool {
		switch event.Key {
		case termbox.KeySpace:
			CreateBox()
		}
		return true
	})

	// CreateBox()

	ui.MainLoop()
}

func main() {
	mainLoop()
}
