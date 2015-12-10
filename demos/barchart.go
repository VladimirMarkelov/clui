package main

import (
	// "fmt"
	ui "github.com/VladimirMarkelov/clui"
	// term "github.com/nsf/termbox-go"
)

func createView(c *ui.Composer) *ui.BarChart {

	view := c.CreateView(0, 0, 10, 7, "BarChart Demo")
	bch := ui.NewBarChart(view, view, 30, 12, 1)

	return bch
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	c := ui.InitLibrary()
	defer c.Close()

	b := createView(c)
	b.SetGap(1)
	d := []ui.BarData{
		{Value: 80, Title: "80%"},
		{Value: 50, Title: "50%"},
		{Value: 150, Title: ">100%"},
	}
	b.SetData(d)
	b.SetValueWidth(5)
	b.SetAutoSize(true)

	// start event processing loop - the main core of the library
	c.MainLoop()
}

func main() {
	mainLoop()
}
