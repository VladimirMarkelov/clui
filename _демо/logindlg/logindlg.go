package main

import (
	ui "../.."
	мИнт "../../пакИнтерфейсы"
)

func createView() {
	view := ui.AddWindow(0, 0, 30, 7, "Диалог логина")
	view.SetPack(мИнт.Vertical)
	view.SetGaps(0, 1)
	view.SetPaddings(2, 2)

	frmOpts := ui.CreateFrame(view, 1, 1, мИнт.BorderNone, мИнт.Fixed)
	frmOpts.SetPack(мИнт.Horizontal)
	cbCheck := ui.CreateCheckBox(frmOpts, мИнт.AutoSize, "Используйте обратный вызов для тестирования данных", мИнт.Fixed)

	ui.CreateLabel(view, мИнт.AutoSize, мИнт.AutoSize, "Корректный вход", мИнт.Fixed)

	frmCreds := ui.CreateFrame(view, 1, 1, мИнт.BorderNone, мИнт.Fixed)
	frmCreds.SetPack(мИнт.Horizontal)
	frmCreds.SetGaps(1, 0)
	ui.CreateLabel(frmCreds, мИнт.AutoSize, мИнт.AutoSize, "Имя:", мИнт.Fixed)
	edUser := ui.CreateEditField(frmCreds, 8, "", 1)
	ui.CreateLabel(frmCreds, мИнт.AutoSize, мИнт.AutoSize, "Пароль:", мИнт.Fixed)
	edPass := ui.CreateEditField(frmCreds, 8, "", 1)

	lbRes := ui.CreateLabel(view, мИнт.AutoSize, мИнт.AutoSize, "Результат:", мИнт.Fixed)

	frmBtns := ui.CreateFrame(view, 1, 1, мИнт.BorderNone, мИнт.Fixed)
	frmBtns.SetPack(мИнт.Horizontal)
	btnDlg := ui.CreateButton(frmBtns, мИнт.AutoSize, 4, "Войти", мИнт.Fixed)
	btnQuit := ui.CreateButton(frmBtns, мИнт.AutoSize, 4, "Выход", мИнт.Fixed)
	ui.CreateFrame(frmBtns, 1, 1, мИнт.BorderNone, 1)

	ui.ActivateControl(view, edUser)

	btnDlg.OnClick(func(ev мИнт.ИСобытие) {
		dlg := ui.CreateLoginDialog(
			"Введите данные:",
			edUser.Title(),
		)

		if cbCheck.State() == 1 {
			dlg.OnCheck(func(u, p string) bool {
				return u == edUser.Title() && p == edPass.Title()
			})
		} else {
			dlg.OnCheck(nil)
		}

		dlg.OnClose(func() {
			if dlg.Action == ui.LoginCanceled {
				lbRes.SetTitle("Результат:\nДиалог отменён")
				return
			}

			if dlg.Action == ui.LoginInvalid {
				lbRes.SetTitle("Результат:\nНеправильный логин или пароль")
				return
			}

			if dlg.Action == ui.LoginOk {
				if cbCheck.State() == 1 {
					lbRes.SetTitle("Результат:\nУспешный вход")
				} else {
					lbRes.SetTitle("Результат:\nВошёл: [" + dlg.Username + ":" + dlg.Password + "]")
				}
				return
			}
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
