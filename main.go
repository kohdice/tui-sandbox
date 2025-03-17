package main

import "github.com/rivo/tview"

func NewLayout(list tview.Primitive, panel tview.Primitive) *tview.Flex {
	bodyLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(list, 20, 1, true).
		AddItem(panel, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(bodyLayout, 0, 1, true)

	return layout
}

func NewList() *tview.List {
	list := tview.NewList()
	list.SetBorder(true).SetTitle("Application")
	return list
}

func NewLog(*tview.Application) *tview.Flex {
	logTable := tview.NewTable()
	logTable.SetBorder(true).SetTitle("Log")

	cnt := 0
	logTable.SetCellSimple(cnt, 0, "Hello")
	logTable.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	log1 := tview.NewTableCell("foo")
	logTable.SetCell(cnt, 1, log1)
	cnt++

	logTable.SetCellSimple(cnt, 0, "World")
	logTable.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	log2 := tview.NewTableCell("bar")
	logTable.SetCell(cnt, 1, log2)
	cnt++

	logPanel := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(logTable, 0, 1, false)
	return logPanel
}

func NewApp() *tview.Application {
	app := tview.NewApplication()
	pages := tview.NewPages()
	logPanel := NewLog(app)

	list := NewList()
	list.AddItem("Test", "", 'p', nil)
	list.AddItem("Quit", "", 'q', func() { app.Stop() })

	layout := NewLayout(list, logPanel)
	pages.AddPage("main", layout, true, true)

	app.SetRoot(pages, true)
	return app
}

func main() {
	app := NewApp()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
