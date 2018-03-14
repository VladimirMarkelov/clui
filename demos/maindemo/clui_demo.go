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
package main

import (
	"fmt"
	"strconv"

	ui "github.com/VladimirMarkelov/clui"
	term "github.com/nsf/termbox-go"
)

func updateProgress(value string, pb *ui.ProgressBar) {
	v, _ := strconv.Atoi(value)
	pb.SetValue(v)
}

func changeTheme(lb *ui.ListBox, btn *ui.Button, tp int) {
	items := ui.ThemeNames()
	dlgType := ui.SelectDialogRadio
	if tp == 1 {
		dlgType = ui.SelectDialogList
	}

	curr := -1
	for i, tName := range items {
		if tName == ui.CurrentTheme() {
			curr = i
			break
		}
	}

	selDlg := ui.CreateSelectDialog("Choose a theme", items, curr, dlgType)
	selDlg.OnClose(func() {
		switch selDlg.Result() {
		case ui.DialogButton1:
			idx := selDlg.Value()
			lb.AddItem(fmt.Sprintf("Selected item: %v", selDlg.Value()))
			lb.SelectItem(lb.ItemCount() - 1)
			if idx != -1 {
				ui.SetCurrentTheme(items[idx])
			}
		}

		btn.SetEnabled(true)
		// ask the composer to repaint all windows
		ui.PutEvent(ui.Event{Type: ui.EventRedraw})
	})
}

func createView() {

	view := ui.AddWindow(0, 0, 20, 7, "Theme Manager Demo")

	frmLeft := ui.CreateFrame(view, 8, 4, ui.BorderNone, 1)
	frmLeft.SetPack(ui.Vertical)
	frmLeft.SetGaps(ui.KeepValue, 1)
	frmLeft.SetPaddings(1, 1)

	frmTheme := ui.CreateFrame(frmLeft, 8, 1, ui.BorderNone, ui.Fixed)
	frmTheme.SetGaps(1, ui.KeepValue)
	checkBox := ui.CreateCheckBox(frmTheme, ui.AutoSize, "Use ListBox", ui.Fixed)
	btnTheme := ui.CreateButton(frmTheme, ui.AutoSize, 4, "Select theme", ui.Fixed)
	ui.CreateFrame(frmLeft, 1, 1, ui.BorderNone, 1)

	frmPb := ui.CreateFrame(frmLeft, 8, 1, ui.BorderNone, ui.Fixed)
	ui.CreateLabel(frmPb, 1, 1, "[", ui.Fixed)
	pb := ui.CreateProgressBar(frmPb, 20, 1, 1)
	pb.SetLimits(0, 10)
	pb.SetTitle("{{value}} of {{max}}")
	ui.CreateLabel(frmPb, 1, 1, "]", ui.Fixed)

	edit := ui.CreateEditField(frmLeft, 5, "0", ui.Fixed)

	frmEdit := ui.CreateFrame(frmLeft, 8, 1, ui.BorderNone, ui.Fixed)
	frmEdit.SetPaddings(1, 1)
	frmEdit.SetGaps(1, ui.KeepValue)
	btnSet := ui.CreateButton(frmEdit, ui.AutoSize, 4, "Set", ui.Fixed)
	btnStep := ui.CreateButton(frmEdit, ui.AutoSize, 4, "Step", ui.Fixed)
	ui.CreateFrame(frmEdit, 1, 1, ui.BorderNone, 1)
	btnQuit := ui.CreateButton(frmEdit, ui.AutoSize, 4, "Quit", ui.Fixed)

	logBox := ui.CreateListBox(view, 28, 5, ui.Fixed)

	ui.ActivateControl(view, edit)

	edit.OnKeyPress(func(key term.Key) bool {
		if key == term.KeyCtrlM {
			v := edit.Title()
			logBox.AddItem(fmt.Sprintf("New PB value(KeyPress): %v", v))
			logBox.SelectItem(logBox.ItemCount() - 1)
			updateProgress(v, pb)
			return true
		}
		return false
	})
	btnTheme.OnClick(func(ev ui.Event) {
		btnTheme.SetEnabled(false)
		tp := checkBox.State()
		changeTheme(logBox, btnTheme, tp)
	})
	btnSet.OnClick(func(ev ui.Event) {
		v := edit.Title()
		logBox.AddItem(fmt.Sprintf("New ProgressBar value: %v", v))
		logBox.SelectItem(logBox.ItemCount() - 1)
		updateProgress(v, pb)
	})
	btnStep.OnClick(func(ev ui.Event) {
		go pb.Step()
		logBox.AddItem("ProgressBar step")
		logBox.SelectItem(logBox.ItemCount() - 1)
		ui.PutEvent(ui.Event{Type: ui.EventRedraw})
	})
	btnQuit.OnClick(func(ev ui.Event) {
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
