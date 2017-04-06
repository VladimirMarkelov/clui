package main

import (
	"fmt"
	ui "github.com/VladimirMarkelov/clui"
)

func createView() *ui.TableView {

	view := ui.AddWindow(0, 0, 10, 7, "TableView Demo")
	bch := ui.CreateTableView(view, 25, 12, 1)
	ui.ActivateControl(view, bch)

	return bch
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	b := createView()
	b.SetShowLines(true)
	b.SetShowRowNumber(true)
	b.SetRowCount(15)
	cols := []ui.Column{
		ui.Column{Title: "Text", Width: 5, Alignment: ui.AlignLeft},
		ui.Column{Title: "Number", Width: 10, Alignment: ui.AlignRight},
		ui.Column{Title: "Misc", Width: 12, Alignment: ui.AlignCenter},
		ui.Column{Title: "Long", Width: 50, Alignment: ui.AlignLeft},
		ui.Column{Title: "Last", Width: 8, Alignment: ui.AlignLeft},
	}
	b.SetColumns(cols)
	b.OnDrawCell(func(info *ui.ColumnDrawInfo) {
		info.Text = fmt.Sprintf("%v:%v", info.Row, info.Col)
	})

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
