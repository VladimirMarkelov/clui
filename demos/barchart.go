package main

import (
	ui "github.com/VladimirMarkelov/clui"
)

func customColored(d *ui.BarDataCell) {
	part := d.TotalMax / 3
	if d.ID%2 == 0 {
		if d.Value <= part {
			d.Fg = ui.ColorGreen
		} else if d.Value > 2*part {
			d.Fg = ui.ColorRed
		} else {
			d.Fg = ui.ColorBlue
		}
	} else {
		d.Ch = '#'
		if d.Value <= part {
			d.Fg = ui.ColorGreenBold
		} else if d.Value > 2*part {
			d.Fg = ui.ColorRedBold
		} else {
			d.Fg = ui.ColorBlueBold
		}
	}
}

func createView() *ui.BarChart {

	view := ui.AddWindow(0, 0, 10, 7, "BarChart Demo")
	bch := ui.CreateBarChart(view, 40, 12, 1)

	frmChk := ui.CreateFrame(view, 8, 5, ui.BorderNone, ui.Fixed)
	frmChk.SetPack(ui.Vertical)
	chkTitles := ui.CreateCheckBox(frmChk, ui.AutoSize, "Show Titles", ui.Fixed)
	chkMarks := ui.CreateCheckBox(frmChk, ui.AutoSize, "Show Marks", ui.Fixed)
	chkTitles.SetState(1)
	chkLegend := ui.CreateCheckBox(frmChk, ui.AutoSize, "Show Legend", ui.Fixed)
	chkValues := ui.CreateCheckBox(frmChk, ui.AutoSize, "Show Values", ui.Fixed)
	chkValues.SetState(1)
	chkFixed := ui.CreateCheckBox(frmChk, ui.AutoSize, "Fixed Width", ui.Fixed)
	chkGap := ui.CreateCheckBox(frmChk, ui.AutoSize, "No Gap", ui.Fixed)
	chkMulti := ui.CreateCheckBox(frmChk, ui.AutoSize, "MultiColored", ui.Fixed)
	chkCustom := ui.CreateCheckBox(frmChk, ui.AutoSize, "Custom Colors", ui.Fixed)

	ui.ActivateControl(view, chkTitles)

	chkTitles.OnChange(func(state int) {
		if state == 0 {
			chkMarks.SetEnabled(false)
			bch.SetShowTitles(false)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			chkMarks.SetEnabled(true)
			bch.SetShowTitles(true)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkMarks.OnChange(func(state int) {
		if state == 0 {
			bch.SetShowMarks(false)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			bch.SetShowMarks(true)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkLegend.OnChange(func(state int) {
		if state == 0 {
			bch.SetLegendWidth(0)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			bch.SetLegendWidth(10)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkValues.OnChange(func(state int) {
		if state == 0 {
			bch.SetValueWidth(0)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			bch.SetValueWidth(5)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkMulti.OnChange(func(state int) {
		if state == 0 {
			d := []ui.BarData{
				{Value: 80, Title: "80%"},
				{Value: 50, Title: "50%"},
				{Value: 150, Title: ">100%"},
			}
			bch.SetData(d)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			d := []ui.BarData{
				{Value: 80, Title: "80%", Fg: ui.ColorBlue},
				{Value: 50, Title: "50%", Fg: ui.ColorGreen, Ch: 'X'},
				{Value: 150, Title: ">100%", Fg: ui.ColorYellow},
			}
			bch.SetData(d)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkFixed.OnChange(func(state int) {
		if state == 0 {
			bch.SetAutoSize(true)
		} else if state == 1 {
			bch.SetAutoSize(false)
		}
		ui.PutEvent(ui.Event{Type: ui.EventRedraw})
	})
	chkGap.OnChange(func(state int) {
		if state == 1 {
			bch.SetBarGap(0)
		} else if state == 0 {
			bch.SetBarGap(1)
		}
		ui.PutEvent(ui.Event{Type: ui.EventRedraw})
	})
	chkCustom.OnChange(func(state int) {
		if state == 0 {
			bch.OnDrawCell(nil)
		} else if state == 1 {
			bch.OnDrawCell(customColored)
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
	b.SetBarGap(1)
	d := []ui.BarData{
		{Value: 80, Title: "80%"},
		{Value: 50, Title: "50%"},
		{Value: 150, Title: ">100%"},
	}
	b.SetData(d)
	b.SetValueWidth(5)
	b.SetAutoSize(true)

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
