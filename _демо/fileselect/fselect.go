package main

import (
	ui "../.."
	мИнт "../../пакИнтерфейсы"
)

func createView() {
	view := ui.AddWindow(0, 0, 30, 7, "Выбор файла")
	view.SetPack(мИнт.Vertical)
	view.SetGaps(0, 1)
	view.SetPaddings(2, 2)

	frmPath := ui.CreateFrame(view, 1, 1, мИнт.BorderNone, мИнт.Fixed)
	frmPath.SetPack(мИнт.Horizontal)
	ui.CreateLabel(frmPath, мИнт.AutoSize, мИнт.AutoSize, "Начальный путь", мИнт.Fixed)
	edPath := ui.CreateEditField(frmPath, 16, "", 1)

	frmMask := ui.CreateFrame(view, 1, 1, мИнт.BorderNone, мИнт.Fixed)
	frmMask.SetPack(мИнт.Horizontal)
	ui.CreateLabel(frmMask, мИнт.AutoSize, мИнт.AutoSize, "Маска файла", мИнт.Fixed)
	edMasks := ui.CreateEditField(frmMask, 16, "*", 1)

	frmOpts := ui.CreateFrame(view, 1, 1, мИнт.BorderNone, мИнт.Fixed)
	frmOpts.SetPack(мИнт.Horizontal)
	cbDir := ui.CreateCheckBox(frmOpts, мИнт.AutoSize, "Выбор папки", мИнт.Fixed)
	cbMust := ui.CreateCheckBox(frmOpts, мИнт.AutoSize, "Должно присутствовать", мИнт.Fixed)
	ui.CreateFrame(frmOpts, 1, 1, мИнт.BorderNone, 1)

	lblSelected := ui.CreateLabel(view, 30, 5, "Выбрано:", мИнт.Fixed)
	lblSelected.SetMultiline(true)

	frmBtns := ui.CreateFrame(view, 1, 1, мИнт.BorderNone, мИнт.Fixed)
	frmBtns.SetPack(мИнт.Horizontal)
	btnSet := ui.CreateButton(frmBtns, мИнт.AutoSize, 4, "Выбрать", мИнт.Fixed)
	btnQuit := ui.CreateButton(frmBtns, мИнт.AutoSize, 4, "Выход", мИнт.Fixed)
	ui.CreateFrame(frmBtns, 1, 1, мИнт.BorderNone, 1)

	ui.ActivateControl(view, edMasks)

	btnSet.OnClick(func(ev мИнт.ИСобытие) {
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
				lblSelected.SetTitle("Выбрано:\nПусто")
				return
			}

			var lb string
			if dlg.Exists {
				lb = "Выбрано существующее"
			} else {
				lb = "Создано новое"
			}

			if cbDir.State() == 0 {
				lb += " файл:\n"
			} else {
				lb += " директория:\n"
			}

			lb += dlg.FilePath
			lblSelected.SetTitle(lb)
		})
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

	createView()

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
