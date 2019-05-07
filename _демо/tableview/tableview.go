package main

import (
	"fmt"
	ui "../.."
	мИнт "../../пакИнтерфейсы"
)

func createView() *ui.TableView {

	view := ui.AddWindow(0, 0, 10, 7, "Пример таблицы")
	bch := ui.CreateTableView(view, 25, 12, 1)
	ui.ActivateControl(view, bch)

	return bch
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	b := createView()
	b.SetShowLines(true)
	b.SetShowRowNumber(true)
	b.SetRowCount(15)
	cols := []ui.Column{
		ui.Column{Title: "Текст", Width: 5, Alignment: мИнт.AlignLeft},
		ui.Column{Title: "Число", Width: 10, Alignment: мИнт.AlignRight},
		ui.Column{Title: "Разное", Width: 12, Alignment: мИнт.AlignCenter},
		ui.Column{Title: "Длина", Width: 50, Alignment: мИнт.AlignLeft},
		ui.Column{Title: "Последний", Width: 8, Alignment: мИнт.AlignLeft},
	}
	b.SetColumns(cols)
	b.OnDrawCell(func(info *ui.ColumnDrawInfo) {
		info.Text = fmt.Sprintf("%v:%v", info.Row, info.Col)
	})

	b.OnAction(func(ev ui.TableEvent) {
		btns := []string{"Закрыть", "Отмена"}
		var action string
		switch ev.Action {
		case мИнт.TableActionSort:
			action = "Sort table"
		case мИнт.TableActionEdit:
			action = "Edit row/cell"
		case мИнт.TableActionNew:
			action = "Add new row"
		case мИнт.TableActionDelete:
			action = "Delete row"
		default:
			action = "Unknown action"
		}

		dlg := ui.CreateConfirmationDialog(
			"<c:blue>"+action,
            "Кликните мышкой или нажмите <c:yellow>SPACE<c:> для закрытия диалога",
			btns, мИнт.DialogButton1)
        dlg.OnClose(func() {})
	})

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
