package main

import (
	"fmt"
	ui "github.com/VladimirMarkelov/clui"
)

func createView() *ui.TableView {

	view := ui.AddWindow(0, 0, 10, 7, "TableView Demo")
	bch := ui.CreateTableView(view, 25, 12, 1)
	ui.ActivateControl(view, bch)

	return bch
}

const rowCount = 15

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	b := createView()
	b.SetShowLines(true)
	b.SetShowRowNumber(true)
	b.SetRowCount(rowCount)
	cols := []ui.Column{
		ui.Column{Title: "Text", Width: 5, Alignment: ui.AlignLeft},
		ui.Column{Title: "Number", Width: 10, Alignment: ui.AlignRight},
		ui.Column{Title: "Misc", Width: 12, Alignment: ui.AlignCenter},
		ui.Column{Title: "Long", Width: 50, Alignment: ui.AlignLeft},
		ui.Column{Title: "Last", Width: 8, Alignment: ui.AlignLeft},
	}
	b.SetColumns(cols)
	colCount := len(cols)

	values := make([]string, rowCount*colCount)
	for r := 0; r < rowCount; r++ {
		for c := 0; c < colCount; c++ {
			values[r*colCount+c] = fmt.Sprintf("%v:%v", r, c)
		}
	}

	b.OnDrawCell(func(info *ui.ColumnDrawInfo) {
		info.Text = values[info.Row*colCount+info.Col]
	})

	b.OnAction(func(ev ui.TableEvent) {
		btns := []string{"Close", "Dismiss"}
		var action string
		switch ev.Action {
		case ui.TableActionSort:
			action = "Sort table"
		case ui.TableActionEdit:
			c := ev.Col
			r := ev.Row
			oldVal := values[r*colCount+c]
			dlg := ui.CreateEditDialog(
				fmt.Sprintf("Editing value: %s", oldVal), "New value", oldVal,
			)
			dlg.OnClose(func() {
				switch dlg.Result() {
				case ui.DialogButton1:
					newText := dlg.EditResult()
					values[r*colCount+c] = newText
					ui.PutEvent(ui.Event{Type: ui.EventRedraw})
				}
			})
			return
		case ui.TableActionNew:
			action = "Add new row"
		case ui.TableActionDelete:
			action = "Delete row"
		default:
			action = "Unknown action"
		}

		dlg := ui.CreateConfirmationDialog(
			"<c:blue>"+action,
			"Click any button or press <c:yellow>SPACE<c:> to close the dialog",
			btns, ui.DialogButton1)
		dlg.OnClose(func() {})
	})

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
