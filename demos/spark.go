package main

import (
	ui "github.com/VladimirMarkelov/clui"
	"math/rand"
	"time"
)

func createView(c *ui.Composer) *ui.SparkChart {

	view := c.CreateView(0, 0, 10, 7, "BarChart Demo")
	bch := ui.NewSparkChart(view, view, 25, 12, 1)
	bch.SetTop(20)

	frmChk := ui.NewFrame(view, view, 8, 5, ui.BorderNone, ui.DoNotScale)
	frmChk.SetPack(ui.Vertical)
	chkValues := ui.NewCheckBox(view, frmChk, ui.AutoSize, "Show Values", ui.DoNotScale)
	chkValues.SetState(0)
	chkHilite := ui.NewCheckBox(view, frmChk, ui.AutoSize, "Hilite peaks", ui.DoNotScale)
	chkHilite.SetState(1)
	chkAuto := ui.NewCheckBox(view, frmChk, ui.AutoSize, "Auto scale", ui.DoNotScale)
	chkAuto.SetState(1)

	chkValues.OnChange(func(state int) {
		if state == 0 {
			bch.SetValueWidth(0)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			bch.SetValueWidth(5)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkHilite.OnChange(func(state int) {
		if state == 0 {
			bch.SetHilitePeaks(false)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			bch.SetHilitePeaks(true)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkAuto.OnChange(func(state int) {
		if state == 0 {
			bch.SetAutoScale(false)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			bch.SetAutoScale(true)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})

	return bch
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	c := ui.InitLibrary()
	defer c.Close()

	b := createView(c)
	b.SetData([]float64{1, 2, 3, 4, 5, 6, 6, 7, 5, 8, 9})

	ticker := time.NewTicker(time.Millisecond * 200).C
	go func() {
		for {
			select {
			case <-ticker:
				b.AddData(float64(rand.Int31n(20)))
				c.PutEvent(ui.Event{Type: ui.EventRedraw})
			}
		}
	}()

	// start event processing loop - the main core of the library
	c.MainLoop()
}

func main() {
	mainLoop()
}
