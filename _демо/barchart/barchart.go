package main

import (
	ui "../.."
	мИнт "../../пакИнтерфейсы"
	мСоб "../../пакСобытия"
)

func customColored(d *ui.BarDataCell) {
	part := d.TotalMax / 3
	if d.ID%2 == 0 {
		if d.Value <= part {
			d.Fg = мИнт.ColorGreen
		} else if d.Value > 2*part {
			d.Fg = мИнт.ColorRed
		} else {
			d.Fg = мИнт.ColorBlue
		}
	} else {
		d.Ch = '#'
		if d.Value <= part {
			d.Fg = мИнт.ColorGreenBold
		} else if d.Value > 2*part {
			d.Fg = мИнт.ColorRedBold
		} else {
			d.Fg = мИнт.ColorBlueBold
		}
	}
}

func createView() *ui.BarChart {

	view := ui.AddWindow(0, 0, 10, 7, "Пример графиков")
	bch := ui.CreateBarChart(view, 40, 12, 1)

	frmChk := ui.CreateFrame(view, 8, 5, мИнт.BorderNone, мИнт.Fixed)
	frmChk.SetPack(мИнт.Vertical)
	chkTitles := ui.CreateCheckBox(frmChk, мИнт.AutoSize, "Показать имена", мИнт.Fixed)
	chkMarks := ui.CreateCheckBox(frmChk, мИнт.AutoSize, "Показать штрихи", мИнт.Fixed)
	chkTitles.SetState(1)
	chkLegend := ui.CreateCheckBox(frmChk, мИнт.AutoSize, "Показать легенду", мИнт.Fixed)
	chkValues := ui.CreateCheckBox(frmChk, мИнт.AutoSize, "Показать значения", мИнт.Fixed)
	chkValues.SetState(1)
	chkFixed := ui.CreateCheckBox(frmChk, мИнт.AutoSize, "Фиксированная ширина", мИнт.Fixed)
	chkGap := ui.CreateCheckBox(frmChk, мИнт.AutoSize, "Без зазоров", мИнт.Fixed)
	chkMulti := ui.CreateCheckBox(frmChk, мИнт.AutoSize, "МультиЦвет", мИнт.Fixed)
	chkCustom := ui.CreateCheckBox(frmChk, мИнт.AutoSize, "Заданный цвет", мИнт.Fixed)

	ui.ActivateControl(view, chkTitles)

	chkTitles.OnChange(func(state int) {
		if state == 0 {
			chkMarks.SetEnabled(false)
			bch.SetShowTitles(false)
			ev := &мСоб.Event{}
			ev.TypeSet(мИнт.EventRedraw)
			ui.PutEvent(ev)
		} else if state == 1 {
			chkMarks.SetEnabled(true)
			bch.SetShowTitles(true)
			ev := &мСоб.Event{}
			ev.TypeSet(мИнт.EventRedraw)
			ui.PutEvent(ev)
		}
	})
	chkMarks.OnChange(func(state int) {
		if state == 0 {
			bch.SetShowMarks(false)
			ev := &мСоб.Event{}
			ev.TypeSet(мИнт.EventRedraw)
			ui.PutEvent(ev)
		} else if state == 1 {
			bch.SetShowMarks(true)
			ev := &мСоб.Event{}
			ev.TypeSet(мИнт.EventRedraw)
			ui.PutEvent(ev)
		}
	})
	chkLegend.OnChange(func(state int) {
		if state == 0 {
			bch.SetLegendWidth(0)
			ev := &мСоб.Event{}
			ev.TypeSet(мИнт.EventRedraw)
			ui.PutEvent(ev)
		} else if state == 1 {
			bch.SetLegendWidth(10)
			ev := &мСоб.Event{}
			ev.TypeSet(мИнт.EventRedraw)
			ui.PutEvent(ev)
		}
	})
	chkValues.OnChange(func(state int) {
		if state == 0 {
			bch.SetValueWidth(0)
			ev := &мСоб.Event{}
			ev.TypeSet(мИнт.EventRedraw)
			ui.PutEvent(ev)
		} else if state == 1 {
			bch.SetValueWidth(5)
			ev := &мСоб.Event{}
			ev.TypeSet(мИнт.EventRedraw)
			ui.PutEvent(ev)
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
			ev := &мСоб.Event{}
			ev.TypeSet(мИнт.EventRedraw)
			ui.PutEvent(ev)
		} else if state == 1 {
			d := []ui.BarData{
				{Value: 80, Title: "80%", Fg: мИнт.ColorBlue},
				{Value: 50, Title: "50%", Fg: мИнт.ColorGreen, Ch: 'X'},
				{Value: 150, Title: ">100%", Fg: мИнт.ColorYellow},
			}
			bch.SetData(d)
			ev := &мСоб.Event{}
			ev.TypeSet(мИнт.EventRedraw)
			ui.PutEvent(ev)
		}
	})
	chkFixed.OnChange(func(state int) {
		if state == 0 {
			bch.SetAutoSize(true)
		} else if state == 1 {
			bch.SetAutoSize(false)
		}
		ev := &мСоб.Event{}
		ev.TypeSet(мИнт.EventRedraw)
		ui.PutEvent(ev)
	})
	chkGap.OnChange(func(state int) {
		if state == 1 {
			bch.SetBarGap(0)
		} else if state == 0 {
			bch.SetBarGap(1)
		}
		ev := &мСоб.Event{}
		ev.TypeSet(мИнт.EventRedraw)
		ui.PutEvent(ev)
	})
	chkCustom.OnChange(func(state int) {
		if state == 0 {
			bch.OnDrawCell(nil)
		} else if state == 1 {
			bch.OnDrawCell(customColored)
		}
		ev := &мСоб.Event{}
		ev.TypeSet(мИнт.EventRedraw)
		ui.PutEvent(ev)
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
