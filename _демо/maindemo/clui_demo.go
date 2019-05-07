package main

/*
Demo includes:
    - How to intialize and run the application
    - How to stop the application
    - How to use Control's events (Button ones)
    - How to change theme on the fly
    - How to use dialogs
    - How to make composer refresh the screen
    - How to intercept Enter key(term.KeyCtrlM) in EditField(ListBox is the same)
*/

import (
	"fmt"
	"strconv"

	ui "../.."
	term "github.com/nsf/termbox-go"
	мИнт "../../пакИнтерфейсы"
	мСоб "../../пакСобытия"
)

func updateProgress(value string, pb *ui.ProgressBar) {
	v, _ := strconv.Atoi(value)
	pb.SetValue(v)
}

func changeTheme(lb *ui.ListBox, btn *ui.Button, tp int) {
	items := ui.ThemeNames()
	dlgType := мИнт.SelectDialogRadio
	if tp == 1 {
		dlgType = мИнт.SelectDialogList
	}

	curr := -1
	for i, tName := range items {
		if tName == ui.CurrentTheme() {
			curr = i
			break
		}
	}

	selDlg := ui.CreateSelectDialog("Выберите тему", items, curr, dlgType)
	selDlg.OnClose(func() {
		switch selDlg.Result() {
		case мИнт.DialogButton1:
			idx := selDlg.Value()
			lb.AddItem(fmt.Sprintf("Выбран пункт: %v", selDlg.Value()))
			lb.SelectItem(lb.ItemCount() - 1)
			if idx != -1 {
				ui.SetCurrentTheme(items[idx])
			}
		}

		btn.SetEnabled(true)
		// ask the composer to repaint all windows
		ev:=&мСоб.Event{}
		ev.TypeSet(мИнт.EventRedraw)
		ui.PutEvent(ev)
	})
}

func createView() {

	view := ui.AddWindow(0, 0, 20, 7, "Приме рменеджера тем")

	frmLeft := ui.CreateFrame(view, 8, 4, мИнт.BorderNone, 1)
	frmLeft.SetPack(мИнт.Vertical)
	frmLeft.SetGaps(мИнт.KeepValue, 1)
	frmLeft.SetPaddings(1, 1)

	frmTheme := ui.CreateFrame(frmLeft, 8, 1, мИнт.BorderNone, мИнт.Fixed)
	frmTheme.SetGaps(1, мИнт.KeepValue)
	checkBox := ui.CreateCheckBox(frmTheme, мИнт.AutoSize, "Использовать ListBox", мИнт.Fixed)
	btnTheme := ui.CreateButton(frmTheme, мИнт.AutoSize, 4, "Выберте тему", мИнт.Fixed)
	ui.CreateFrame(frmLeft, 1, 1, мИнт.BorderNone, 1)

	frmPb := ui.CreateFrame(frmLeft, 8, 1, мИнт.BorderNone, мИнт.Fixed)
	ui.CreateLabel(frmPb, 1, 1, "[", мИнт.Fixed)
	pb := ui.CreateProgressBar(frmPb, 20, 1, 1)
	pb.SetLimits(0, 10)
	pb.SetTitle("{{value}} of {{max}}")
	ui.CreateLabel(frmPb, 1, 1, "]", мИнт.Fixed)

	edit := ui.CreateEditField(frmLeft, 5, "0", мИнт.Fixed)

	frmEdit := ui.CreateFrame(frmLeft, 8, 1, мИнт.BorderNone, мИнт.Fixed)
	frmEdit.SetPaddings(1, 1)
	frmEdit.SetGaps(1, мИнт.KeepValue)
	btnSet := ui.CreateButton(frmEdit, мИнт.AutoSize, 4, "Установить", мИнт.Fixed)
	btnStep := ui.CreateButton(frmEdit, мИнт.AutoSize, 4, "Шаг", мИнт.Fixed)
	ui.CreateFrame(frmEdit, 1, 1, мИнт.BorderNone, 1)
	btnQuit := ui.CreateButton(frmEdit, мИнт.AutoSize, 4, "Выход", мИнт.Fixed)

	logBox := ui.CreateListBox(view, 28, 5, мИнт.Fixed)

	ui.ActivateControl(view, edit)

	edit.OnKeyPress(func(key term.Key, ch rune) bool {
		if key == term.KeyCtrlM {
			v := edit.Title()
			logBox.AddItem(fmt.Sprintf("Новое PB значение(KeyPress): %v", v))
			logBox.SelectItem(logBox.ItemCount() - 1)
			updateProgress(v, pb)
			return true
		}
		return false
	})
	btnTheme.OnClick(func(ev мИнт.ИСобытие) {
		btnTheme.SetEnabled(false)
		tp := checkBox.State()
		changeTheme(logBox, btnTheme, tp)
	})
	btnSet.OnClick(func(ev мИнт.ИСобытие) {
		v := edit.Title()
		logBox.AddItem(fmt.Sprintf("Новое значение ProgressBar: %v", v))
		logBox.SelectItem(logBox.ItemCount() - 1)
		updateProgress(v, pb)
	})
	btnStep.OnClick(func(ev мИнт.ИСобытие) {
		go pb.Step()
		logBox.AddItem("Шаг ProgressBar")
		logBox.SelectItem(logBox.ItemCount() - 1)
		ev=&мСоб.Event{}
		ev.TypeSet(мИнт.EventRedraw)
		ui.PutEvent(ev)
	})
	btnQuit.OnClick(func(ev мИнт.ИСобытие) {
		go ui.Stop()
	})
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	ui.SetThemePath("themes")

	createView()

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
