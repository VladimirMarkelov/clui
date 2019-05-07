package main

import (
	"math/rand"
	"time"

	ui "../.."
	мИнт "../../пакИнтерфейсы"
	мСоб "../../пакСобытия"
)

func createView() *ui.SparkChart {

	view := ui.AddWindow(0, 0, 10, 7, "Пример графика")
	bch := ui.CreateSparkChart(view, 25, 12, 1)
	bch.SetTop(20)

	frmChk := ui.CreateFrame(view, 8, 5, мИнт.BorderNone, мИнт.Fixed)
	frmChk.SetPack(мИнт.Vertical)
	chkValues := ui.CreateCheckBox(frmChk, мИнт.AutoSize, "Значения", мИнт.Fixed)
	chkValues.SetState(0)
	chkHilite := ui.CreateCheckBox(frmChk, мИнт.AutoSize, "Пики", мИнт.Fixed)
	chkHilite.SetState(1)
	chkAuto := ui.CreateCheckBox(frmChk, мИнт.AutoSize, "Авто масштаб", мИнт.Fixed)
	chkAuto.SetState(1)

	ui.ActivateControl(view, chkValues)

	chkValues.OnChange(func(state int) {
		if state == 0 {
			bch.SetValueWidth(0)
		} else if state == 1 {
			bch.SetValueWidth(5)
		}
		ev := &мСоб.Event{}
		ev.TypeSet(мИнт.EventRedraw)
		ui.PutEvent(ev)
	})
	chkHilite.OnChange(func(state int) {
		if state == 0 {
			bch.SetHilitePeaks(false)
		} else if state == 1 {
			bch.SetHilitePeaks(true)
		}
		ev := &мСоб.Event{}
		ev.TypeSet(мИнт.EventRedraw)
		ui.PutEvent(ev)
	})
	chkAuto.OnChange(func(state int) {
		if state == 0 {
			bch.SetAutoScale(false)
		} else if state == 1 {
			bch.SetAutoScale(true)
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
	b.SetData([]float64{1, 2, 3, 4, 5, 6, 6, 7, 5, 8, 9})

	ticker := time.NewTicker(time.Millisecond * 200).C
	go func() {
		for {
			select {
			case <-ticker:
				b.AddData(float64(rand.Int31n(20)))
				ev := &мСоб.Event{}
				ev.TypeSet(мИнт.EventRedraw)
				ui.PutEvent(ev)
			}
		}
	}()

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
