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
	ui "github.com/VladimirMarkelov/clui"
	term "github.com/nsf/termbox-go"
	"strconv"
)

func updateProgress(value string, pb *ui.ProgressBar) {
	v, _ := strconv.Atoi(value)
	pb.SetValue(v)
}

func changeTheme(c *ui.Composer, lb *ui.ListBox, btn *ui.Button, tp int) {
	items := c.Theme().ThemeNames()
	dlgType := ui.SelectDialogRadio
	if tp == 1 {
		dlgType = ui.SelectDialogList
	}

	curr := -1
	for i, tName := range items {
		if tName == c.Theme().CurrentTheme() {
			curr = i
			break
		}
	}

	selDlg := ui.NewSelectDialog(c, "Choose a theme", items, curr, dlgType)
	selDlg.OnClose(func() {
		switch selDlg.Result() {
		case ui.DialogButton1:
			idx := selDlg.Value()
			lb.AddItem(fmt.Sprintf("Selected item: %v", selDlg.Value()))
			lb.SelectItem(lb.ItemCount() - 1)
			if idx != -1 {
				c.Theme().SetCurrentTheme(items[idx])
			}
		}

		btn.SetEnabled(true)
		// ask the composer to repaint all windows
		c.PutEvent(ui.Event{Type: ui.EventRedraw})
	})
}

func createView(c *ui.Composer) {

	view := c.CreateView(0, 0, 20, 7, "Theme Manager Demo")

	frmLeft := ui.NewFrame(view, view, 8, 4, ui.BorderNone, 1)
	frmLeft.SetPack(ui.Vertical)
	frmLeft.SetPaddings(1, 1, ui.DoNotChange, 1)

	frmTheme := ui.NewFrame(view, frmLeft, 8, 1, ui.BorderNone, ui.DoNotScale)
	frmTheme.SetPaddings(ui.DoNotChange, ui.DoNotChange, 1, ui.DoNotChange)
	checkBox := ui.NewCheckBox(view, frmTheme, ui.AutoSize, "Use ListBox", ui.DoNotScale)
	btnTheme := ui.NewButton(view, frmTheme, ui.AutoSize, 4, "Select theme", ui.DoNotScale)
	ui.NewFrame(view, frmLeft, 1, 1, ui.BorderNone, 1)

	frmPb := ui.NewFrame(view, frmLeft, 8, 1, ui.BorderNone, ui.DoNotScale)
	ui.NewLabel(view, frmPb, 1, 1, "[", ui.DoNotScale)
	pb := ui.NewProgressBar(view, frmPb, 20, 1, 1)
	pb.SetLimits(0, 10)
	pb.SetTitle("{{value}} of {{max}}")
	ui.NewLabel(view, frmPb, 1, 1, "]", ui.DoNotScale)

	edit := ui.NewEditField(view, frmLeft, 5, "0", ui.DoNotScale)

	frmEdit := ui.NewFrame(view, frmLeft, 8, 1, ui.BorderNone, ui.DoNotScale)
	frmEdit.SetPaddings(1, 1, 1, ui.DoNotChange)
	btnSet := ui.NewButton(view, frmEdit, ui.AutoSize, 4, "Set", ui.DoNotScale)
	btnStep := ui.NewButton(view, frmEdit, ui.AutoSize, 4, "Step", ui.DoNotScale)
	ui.NewFrame(view, frmEdit, 1, 1, ui.BorderNone, 1)
	btnQuit := ui.NewButton(view, frmEdit, ui.AutoSize, 4, "Quit", ui.DoNotScale)

	logBox := ui.NewListBox(view, view, 28, 5, ui.DoNotScale)

	view.ActivateControl(edit)

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
		changeTheme(c, logBox, btnTheme, tp)
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
		c.PutEvent(ui.Event{Type: ui.EventRedraw})
	})
	btnQuit.OnClick(func(ev ui.Event) {
		go c.Stop()
	})
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	c := ui.InitLibrary()
	defer c.Close()

	c.Theme().SetThemePath("themes")

	createView(c)

	// start event processing loop - the main core of the library
	c.MainLoop()
}

func main() {
	mainLoop()
}
