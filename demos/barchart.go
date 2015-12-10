package main

import (
	// "fmt"
	ui "github.com/VladimirMarkelov/clui"
	// term "github.com/nsf/termbox-go"
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

func createView(c *ui.Composer) *ui.BarChart {

	view := c.CreateView(0, 0, 10, 7, "BarChart Demo")
	bch := ui.NewBarChart(view, view, 40, 12, 1)

	frmChk := ui.NewFrame(view, view, 8, 5, ui.BorderNone, ui.DoNotScale)
	frmChk.SetPack(ui.Vertical)
	chkTitles := ui.NewCheckBox(view, frmChk, ui.AutoSize, "Show Titles", ui.DoNotScale)
	chkMarks := ui.NewCheckBox(view, frmChk, ui.AutoSize, "Show Marks", ui.DoNotScale)
	chkTitles.SetState(1)
	chkLegend := ui.NewCheckBox(view, frmChk, ui.AutoSize, "Show Legend", ui.DoNotScale)
	chkValues := ui.NewCheckBox(view, frmChk, ui.AutoSize, "Show Values", ui.DoNotScale)
	chkValues.SetState(1)
	chkFixed := ui.NewCheckBox(view, frmChk, ui.AutoSize, "Fixed Width", ui.DoNotScale)
	chkGap := ui.NewCheckBox(view, frmChk, ui.AutoSize, "No Gap", ui.DoNotScale)
	chkMulti := ui.NewCheckBox(view, frmChk, ui.AutoSize, "MultiColored", ui.DoNotScale)
	chkCustom := ui.NewCheckBox(view, frmChk, ui.AutoSize, "Custom Colors", ui.DoNotScale)

	chkTitles.OnChange(func(state int) {
		if state == 0 {
			chkMarks.SetEnabled(false)
			bch.SetShowTitles(false)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			chkMarks.SetEnabled(true)
			bch.SetShowTitles(true)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkMarks.OnChange(func(state int) {
		if state == 0 {
			bch.SetShowMarks(false)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			bch.SetShowMarks(true)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkLegend.OnChange(func(state int) {
		if state == 0 {
			bch.SetLegendWidth(0)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			bch.SetLegendWidth(10)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkValues.OnChange(func(state int) {
		if state == 0 {
			bch.SetValueWidth(0)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			bch.SetValueWidth(5)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
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
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			d := []ui.BarData{
				{Value: 80, Title: "80%", Fg: ui.ColorBlue},
				{Value: 50, Title: "50%", Fg: ui.ColorGreen, Ch: 'X'},
				{Value: 150, Title: ">100%", Fg: ui.ColorYellow},
			}
			bch.SetData(d)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkFixed.OnChange(func(state int) {
		if state == 0 {
			bch.SetAutoSize(true)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			bch.SetAutoSize(false)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkGap.OnChange(func(state int) {
		if state == 1 {
			bch.SetGap(0)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 0 {
			bch.SetGap(1)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	})
	chkCustom.OnChange(func(state int) {
		if state == 0 {
			bch.OnDrawCell(nil)
			c.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 1 {
			bch.OnDrawCell(customColored)
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
