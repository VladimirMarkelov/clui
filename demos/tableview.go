package main

import (
	"fmt"
	ui "github.com/VladimirMarkelov/clui"
)

func createView(c *ui.Composer) *ui.TableView {

	view := c.CreateView(0, 0, 10, 7, "TableView Demo")
	bch := ui.NewTableView(view, view, 25, 12, 1)

	return bch
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	c := ui.InitLibrary()
	defer c.Close()

	b := createView(c)
	b.SetShowLines(true)
	b.SetShowRowNumber(true)
	b.SetRowCount(5)
	cols := []ui.Column{
		ui.Column{Title: "Text", Width: 5, Alignment: ui.AlignLeft},
		ui.Column{Title: "Number", Width: 10, Alignment: ui.AlignRight},
	}
	b.SetColumns(cols)
	b.OnDrawCell(func(info *ui.ColumnDrawInfo) {
		info.Text = fmt.Sprintf("%v:%v", info.Row, info.Col)
	})

	// start event processing loop - the main core of the library
	c.MainLoop()
}

func main() {
	mainLoop()
}
