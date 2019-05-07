package main

/*
 * Demo includes:
 *  - How to use OnBeforeDraw event
 *  - a simple example of "DBCache" for faster drawing
 */

import (
	"fmt"
	ui "../.."
	мИнт "../../пакИнтерфейсы"
)

// number of columns in a table
const columnInTable = 6

// dbCache for data from DB. It always caches the whole table row, so it does not
// use firstCol and colCount values from OnBeforeDraw event. But you can do more
// granular storage to minimize memory usage by cache
// dbCache is quite dumb: if it detects that topRow or the number of visible rows
// is changed it invalidates the cache and reloads all the data from new row span.
// In real application, it would be good to make it smarter, e.g:
//    - if rowCount descreased and firstRow does not change - the cache is valid,
//      and redundant rereading data can be skipped
//    - usually visible area changes by 1 row, so performance-wise the cache can
//      shift row slice and read only new rows
//    - etc
type dbCache struct {
	firstRow int        // previous first visible row
	rowCount int        // previous visible row count
	data     [][]string // cache - contains at least 'rowCount' rows from DB
}

// cache data from a new row span
// It imitates a random data by selecting values from predefined arrays. Sizes
// of all arrays should be different to make TableView data look more random
func (d *dbCache) preload(firstRow, rowCount int) {
	if firstRow == d.firstRow && rowCount == d.rowCount {
		// fast path: view area is the same, return immediately
		return
	}

	// slow path: refill cache
	fNames := []string{"Джек", "Алиса", "Ричард", "Паша", "Николь", "Стивен", "Жан"}
	lNames := []string{"Смит", "Качер", "Стоун", "Белов", "Васин"}
	posts := []string{"Инженер", "Менеджер", "Охранник", "Водитель"}
	deps := []string{"ИТ", "Финансы", "Обеспечение"}
	salary := []int{40000, 38000, 41000, 32000}

	d.data = make([][]string, rowCount, rowCount)
	for i := 0; i < rowCount; i++ {
		absIndex := firstRow + i
		d.data[i] = make([]string, columnInTable, columnInTable)
		d.data[i][0] = fNames[absIndex%len(fNames)]
		d.data[i][1] = lNames[absIndex%len(lNames)]
		d.data[i][2] = fmt.Sprintf("%08d", 100+absIndex)
		d.data[i][3] = posts[absIndex%len(posts)]
		d.data[i][4] = deps[absIndex%len(deps)]
		d.data[i][5] = fmt.Sprintf("%d руб/год", salary[absIndex%len(salary)]/1000)
	}

	// do not forget to save the last values
	d.firstRow = firstRow
	d.rowCount = rowCount
}

// returns the cell value for a given col and row. Col and row are absolute
// value. But cache keeps limited number of rows to minimize memory usage.
// So, the position of the value of the cell should be calculated
// To simplify, the function just returns empty string if the cell is not
// cached. It is unlikely but can happen
func (d *dbCache) value(row, col int) string {
	rowId := row - d.firstRow
	if rowId >= len(d.data) {
		return ""
	}
	rowValues := d.data[rowId]
	if col >= len(rowValues) {
		return ""
	}
	return rowValues[col]
}

var (
	view *ui.Window
)

func createView() *ui.TableView {
	view = ui.AddWindow(0, 0, 10, 7, "Загруженные данные таблицы")
	bch := ui.CreateTableView(view, 35, 12, 1)
	ui.ActivateControl(view, bch)

	return bch
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	cache := &dbCache{firstRow: -1}
	b := createView()
	b.SetShowLines(true)
	b.SetShowRowNumber(true)
	b.SetRowCount(25)
	cols := []ui.Column{
		ui.Column{Title: "Имя", Width: 10, Alignment: мИнт.AlignLeft},
		ui.Column{Title: "Фамилия", Width: 12, Alignment: мИнт.AlignLeft},
		ui.Column{Title: "Номер", Width: 12, Alignment: мИнт.AlignRight},
		ui.Column{Title: "Адрес", Width: 12, Alignment: мИнт.AlignLeft},
		ui.Column{Title: "Отдел", Width: 15, Alignment: мИнт.AlignLeft},
		ui.Column{Title: "Доход", Width: 12, Alignment: мИнт.AlignRight},
	}
	b.SetColumns(cols)
	b.OnBeforeDraw(func(col, row, colCnt, rowCnt int) {
		cache.preload(row, rowCnt)
		l, t, w, h := b.VisibleArea()
		view.SetTitle(fmt.Sprintf("Выборка: %d:%d - %dx%d", l, t, w, h))
	})
	b.OnDrawCell(func(info *ui.ColumnDrawInfo) {
		info.Text = cache.value(info.Row, info.Col)
	})

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
