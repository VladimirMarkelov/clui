package main

import (
	"fmt"
	ui "github.com/VladimirMarkelov/clui"
)

func createView() *ui.TextReader {

	view := ui.AddWindow(0, 0, 10, 7, "TableView Demo")
	bch := ui.CreateTextReader(view, 45, 24, 1)
	ui.ActivateControl(view, bch)

	return bch
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	b := createView()
    _ = b
    b.SetLineCount(50)
	b.OnDrawLine(func(ind int) string {
		return fmt.Sprintf("%03d line line line", ind+1)
	})

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
