package main

import (
	ui "github.com/VladimirMarkelov/clui"
	"math/rand"
	"time"
)

func createView() *ui.SparkChart {

	view := ui.AddWindow(0, 0, 10, 7, "BarChart Demo")
	bch := ui.CreateSparkChart(view, 25, 12, 1)
	bch.SetTop(20)

	frmChk := ui.CreateFrame(view, 8, 5, ui.BorderNone, ui.Fixed)
	frmChk.SetPack(ui.Vertical)
	chkValues := ui.CreateCheckBox(frmChk, ui.AutoSize, "Show Values", ui.Fixed)
	chkValues.SetState(0)
	chkHilite := ui.CreateCheckBox(frmChk, ui.AutoSize, "Hilite peaks", ui.Fixed)
	chkHilite.SetState(1)
	chkAuto := ui.CreateCheckBox(frmChk, ui.AutoSize, "Auto scale", ui.Fixed)
	chkAuto.SetState(1)

	ui.ActivateControl(view, chkValues)

	chkValues.OnChange(func(state int) {
		if state == 0 {
			bch.SetValueWidth(0)
		} else if state == 1 {
			bch.SetValueWidth(5)
		}
		ui.PutEvent(ui.Event{Type: ui.EventRedraw})
	})
	chkHilite.OnChange(func(state int) {
		if state == 0 {
			bch.SetHilitePeaks(false)
		} else if state == 1 {
			bch.SetHilitePeaks(true)
		}
		ui.PutEvent(ui.Event{Type: ui.EventRedraw})
	})
	chkAuto.OnChange(func(state int) {
		if state == 0 {
			bch.SetAutoScale(false)
		} else if state == 1 {
			bch.SetAutoScale(true)
		}
		ui.PutEvent(ui.Event{Type: ui.EventRedraw})
	})

	return bch
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	b := createView()
	b.SetData([]float64{1, 2, 3, 4, 5, 6, 6, 7, 5, 8, 9})

	ticker := time.NewTicker(time.Millisecond * 200).C
	go func() {
		for {
			select {
			case <-ticker:
				b.AddData(float64(rand.Int31n(20)))
				ui.PutEvent(ui.Event{Type: ui.EventRedraw})
			}
		}
	}()

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
