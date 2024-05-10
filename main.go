package main

import (
	"os"
	"path/filepath"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	currentPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(currentPath)
	pathSecond := filepath.Dir(currentPath)

	grid := setupUI()
	addToList(list, listTitle, path)
	addToList(listSecond, listSecondTitle, path)
	eventHandler(app, list, listSecond, pathSecond, path, resultLabel)

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
