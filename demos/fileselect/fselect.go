package main

import (
	ui "github.com/VladimirMarkelov/clui"
)

func createView() {
	view := ui.AddWindow(0, 0, 30, 7, "File select")
	view.SetPack(ui.Vertical)
	view.SetGaps(0, 1)
	view.SetPaddings(2, 2)

	frmPath := ui.CreateFrame(view, 1, 1, ui.BorderNone, ui.Fixed)
	frmPath.SetPack(ui.Horizontal)
	ui.CreateLabel(frmPath, ui.AutoSize, ui.AutoSize, "Initial path", ui.Fixed)
	edPath := ui.CreateEditField(frmPath, 16, "", 1)

	frmMask := ui.CreateFrame(view, 1, 1, ui.BorderNone, ui.Fixed)
	frmMask.SetPack(ui.Horizontal)
	ui.CreateLabel(frmMask, ui.AutoSize, ui.AutoSize, "File masks", ui.Fixed)
	edMasks := ui.CreateEditField(frmMask, 16, "*", 1)

	frmOpts := ui.CreateFrame(view, 1, 1, ui.BorderNone, ui.Fixed)
	frmOpts.SetPack(ui.Horizontal)
	cbDir := ui.CreateCheckBox(frmOpts, ui.AutoSize, "Select directory", ui.Fixed)
	cbMust := ui.CreateCheckBox(frmOpts, ui.AutoSize, "Must exists", ui.Fixed)
	ui.CreateFrame(frmOpts, 1, 1, ui.BorderNone, 1)

	lblSelected := ui.CreateLabel(view, 30, 5, "Selected:", ui.Fixed)
	lblSelected.SetMultiline(true)

	frmBtns := ui.CreateFrame(view, 1, 1, ui.BorderNone, ui.Fixed)
	frmBtns.SetPack(ui.Horizontal)
	btnSet := ui.CreateButton(frmBtns, ui.AutoSize, 4, "Select", ui.Fixed)
	btnQuit := ui.CreateButton(frmBtns, ui.AutoSize, 4, "Quit", ui.Fixed)
	ui.CreateFrame(frmBtns, 1, 1, ui.BorderNone, 1)

	ui.ActivateControl(view, edMasks)

	btnSet.OnClick(func(ev ui.Event) {
		s := "Select "
		if cbDir.State() == 1 {
			s += "directory"
		} else {
			s += "file"
		}
		if cbMust.State() == 1 {
			s += "[X]"
		}
		dlg := ui.CreateFileSelectDialog(
			s,
			edMasks.Title(),
			edPath.Title(),
			cbDir.State() == 1,
			cbMust.State() == 1)
		dlg.OnClose(func() {
			if !dlg.Selected {
				lblSelected.SetTitle("Selected:\nNothing")
				return
			}

			var lb string
			if dlg.Exists {
				lb = "Selected existing"
			} else {
				lb = "Create new"
			}

			if cbDir.State() == 0 {
				lb += " file:\n"
			} else {
				lb += " directory:\n"
			}

			lb += dlg.FilePath
			lblSelected.SetTitle(lb)
		})
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

	createView()

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
