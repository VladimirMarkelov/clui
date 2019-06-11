package main

import (
	ui "github.com/VladimirMarkelov/clui"
)

func createView() {
	view := ui.AddWindow(0, 0, 10, 7, "Button` Demo")
	view.SetTitleButtons(ui.ButtonMaximize | ui.ButtonClose)

	frmViews := ui.CreateFrame(view, 8, 5, ui.BorderNone, ui.Fixed)
	frmViews.SetPack(ui.Horizontal)
	frmFull := ui.CreateFrame(frmViews, 8, 5, ui.BorderThin, ui.Fixed)
	frmFull.SetPack(ui.Vertical)
	frmFull.SetTitle("Full")
	frmHalf := ui.CreateFrame(frmViews, 8, 5, ui.BorderThin, ui.Fixed)
	frmHalf.SetPack(ui.Vertical)
	frmHalf.SetTitle("Half")
	frmNone := ui.CreateFrame(frmViews, 8, 5, ui.BorderThin, ui.Fixed)
	frmNone.SetPack(ui.Vertical)
	frmNone.SetTitle("None")

	btnF1 := ui.CreateButton(frmFull, ui.AutoSize, 4, "First", ui.Fixed)
	btnF2 := ui.CreateButton(frmFull, ui.AutoSize, 4, "Second", ui.Fixed)
	btnF3 := ui.CreateButton(frmFull, ui.AutoSize, 4, "Quit", ui.Fixed)
	btnF1.SetShadowType(ui.ShadowFull)
	btnF2.SetShadowType(ui.ShadowFull)
	btnF3.SetShadowType(ui.ShadowFull)
	btnH1 := ui.CreateButton(frmHalf, ui.AutoSize, 4, "First", ui.Fixed)
	btnH2 := ui.CreateButton(frmHalf, ui.AutoSize, 4, "Second", ui.Fixed)
	btnH3 := ui.CreateButton(frmHalf, ui.AutoSize, 4, "Quit", ui.Fixed)
	btnH1.SetShadowType(ui.ShadowHalf)
	btnH2.SetShadowType(ui.ShadowHalf)
	btnH3.SetShadowType(ui.ShadowHalf)
	btnN1 := ui.CreateButton(frmNone, ui.AutoSize, 4, "First", ui.Fixed)
	btnN2 := ui.CreateButton(frmNone, ui.AutoSize, 4, "Second", ui.Fixed)
	btnN3 := ui.CreateButton(frmNone, ui.AutoSize, 4, "Quit", ui.Fixed)
	btnN1.SetShadowType(ui.ShadowNone)
	btnN2.SetShadowType(ui.ShadowNone)
	btnN3.SetShadowType(ui.ShadowNone)

	btnF3.OnClick(func(ev ui.Event) {
		go ui.Stop()
	})
	btnH3.OnClick(func(ev ui.Event) {
		go ui.Stop()
	})
	btnN3.OnClick(func(ev ui.Event) {
		go ui.Stop()
	})

	ui.ActivateControl(view, btnF1)
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
