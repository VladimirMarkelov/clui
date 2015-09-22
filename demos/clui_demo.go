/*
The Demo shows how to create Windows with fixed(manual
layout) and dynamic one. Simple forms are easier to create
in manual layout(check method createManualView).

Below there are two Windows that look and act the same
way. The only difference: it is not easy to relocate buttons
on Window resize in the same way as it is done in Window with
dynamic layout. So, Buttons in manual layout are moved in
a bit different way. All other features are the same.

Demo includes:
    - How to use Control's events (ComboBox and Button ones)
    - How to create manual layout
    - How to use Packers
    - How to intialize and run the application
    - How to stop the application
*/
package main

import (
	"fmt"
	ui "github.com/VladimirMarkelov/clui"
	"strconv"
)

func updateProgress(value string, pb *ui.ProgressBar) {
	v, _ := strconv.Atoi(value)
	pb.SetValue(v)
}

func createDynamicView(c *ui.Composer) {
	emptyProp := ui.Props{}

	wnd := c.CreateWindow(1, 15, ui.AutoSize, ui.AutoSize, "Dynamic Layout", emptyProp)

	packMain := wnd.AddPack(ui.PackHorizontal)
	packLeft := packMain.AddPack(ui.PackVertical, 1)
	packLeft.SetBorderStyle(ui.BorderSingle)

	packRight := packMain.AddPack(ui.PackVertical, ui.DoNotScale)
	packRight.SetBorderStyle(ui.BorderDouble)
	ts := packRight.PackTextScroll(30, 10, 1, ui.Props{})

	// Here are a lot of additional Packers because it is not possible to mix
	// Packers and Controls inside one Packer
	packPb := packLeft.AddPack(ui.PackHorizontal, ui.DoNotScale)
	packPb.SetPaddings(ui.DoNotChange, 1, ui.DoNotChange, ui.DoNotChange)
	pb := packPb.PackProgressBar(10, 1, 0, 10, 1, emptyProp)
	packCb := packLeft.AddPack(ui.PackHorizontal, ui.DoNotScale)
	packCb.PackLabel(22, "Set ProgressBar Value", ui.DoNotScale, emptyProp)
	cb := packCb.PackComboBox(5, "0", 1, ui.Props{Text: "0|1|2|3|4|5|6|7|8|9|10"})
	cb.OnChange(func(ev ui.Event) {
		ts.AddItem(fmt.Sprintf("ComboBox changed to %v", ev.Msg))
	})
	packFiller := packLeft.AddPack(ui.PackHorizontal, 1)
	packFiller.PackFrame(1, 1, "", 1, emptyProp)
	packBtn := packLeft.AddPack(ui.PackHorizontal, ui.DoNotScale)
	btnSet := packBtn.PackButton(11, 3, "Set Value", ui.DoNotScale, emptyProp)
	btnSet.OnClick(func(ev ui.Event) {
		v := cb.GetText()
		ts.AddItem(fmt.Sprintf("New ProgressBar value: %v", v))
		updateProgress(v, pb)
	})
	packBtn.PackFrame(2, 1, "", 1, emptyProp)
	btnStep := packBtn.PackButton(6, 3, "Step", ui.DoNotScale, emptyProp)
	btnStep.OnClick(func(ev ui.Event) {
		go pb.Step()
		ts.AddItem("ProgressBar step")
	})
	packBtn.PackFrame(2, 1, "", 1, emptyProp)
	btnQuit := packBtn.PackButton(6, 3, "Quit", ui.DoNotScale, emptyProp)
	btnQuit.OnClick(func(ev ui.Event) {
		go c.Stop()
	})
	packBtn.PackFrame(1, 1, "", ui.DoNotScale, emptyProp)

	// Method must be called after all Window controls are added to it
	// Otherwise window won't display anything
	wnd.PackEnd()
}

func createManualView(c *ui.Composer) {
	emptyProp := ui.Props{}

	w, h := 64, 14
	wnd := c.CreateWindow(1, 0, w, h, "Manual Layout", emptyProp)
	wnd.SetConstraints(w, h)
	ui.CreateFrame(wnd, 0, 0, 30, h-2, "Task Progress", ui.Props{Border: ui.BorderSingle, Anchors: ui.AnchorAll})
	ui.CreateFrame(wnd, 30, 0, w-2-30, h-2, "Event List", ui.Props{Border: ui.BorderSingle, Anchors: ui.AnchorRight | ui.AnchorHeight})

	ts := ui.CreateTextScroll(wnd, 31, 1, w-2-30-2, h-4, ui.Props{Anchors: ui.AnchorRight | ui.AnchorHeight})

	pb := ui.CreateProgressBar(wnd, 1, 2, 30-2, 1, 0, 10, ui.Props{Anchors: ui.AnchorWidth})
	ui.CreateLabel(wnd, 1, 4, 22, "Set ProgressBar Value", emptyProp)
	cb := ui.CreateComboBox(wnd, 1+22, 4, 6, "0", ui.Props{Anchors: ui.AnchorWidth})
	var _ = pb
	cb.OnChange(func(ev ui.Event) {
		ts.AddItem(fmt.Sprintf("ComboBox changed to %v", ev.Msg))
	})

	btnSet := ui.CreateButton(wnd, 1, 6, 11, 3, "Set Value", ui.Props{Anchors: ui.AnchorBottom})
	btnSet.OnClick(func(ev ui.Event) {
		v := cb.GetText()
		ts.AddItem(fmt.Sprintf("New ProgressBar value: %v", v))
		updateProgress(v, pb)
	})

	btnStep := ui.CreateButton(wnd, 1+11+2, 6, 6, 3, "Step", ui.Props{Anchors: ui.AnchorBottom})
	btnStep.OnClick(func(ev ui.Event) {
		go pb.Step()
		ts.AddItem("ProgressBar step")
	})
	btnQuit := ui.CreateButton(wnd, 1+11+2+6+2, 6, 6, 3, "Quit", ui.Props{Anchors: ui.AnchorBottom | ui.AnchorRight})
	btnQuit.OnClick(func(ev ui.Event) {
		go c.Stop()
	})
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	var c ui.Composer
	c.Init()
	defer c.Close()

	createManualView(&c)
	createDynamicView(&c)

	c.RefreshScreen()
	// start event precessing loop - the main core of the library
	c.MainLoop()
}

func main() {
	mainLoop()
}
