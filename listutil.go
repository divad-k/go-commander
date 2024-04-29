package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func addToList(target *tview.List, label *tview.TextView, path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			target.AddItem(file.Name(), "", 0, nil)
		} else {
			target.AddItem(file.Name(), "", 0, nil)
		}
	}
	label.SetText(fmt.Sprintf("Path: %s", path))
}

func eventHandler(app *tview.Application, list, listSecond *tview.List, pathSecond, path string, resultLabel *tview.TextView) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case app.GetFocus() == list:
			// Event handling for list 1
			switch event.Key() {
			case tcell.KeyRight:
				selectedIndex := list.GetCurrentItem()
				selectedItem, _ := list.GetItemText(selectedIndex)
				filePath := filepath.Join(path, selectedItem)
				fileInfo, err := os.Stat(filePath)

				if err != nil {
					panic(err)
				}

				if fileInfo.IsDir() {
					path = filePath
					list.Clear()
					addToList(list, listTitle, path)
				} else {
					fileType, fileFormatted, fileSize, creationTime := displaySingleFileInfo(filePath)
					resultLabel.SetText(fmt.Sprintf("Path: %s\n Type: %s\n Size: %.2f %s\n LastModified: %s", filePath, fileType, fileFormatted, fileSize, creationTime))
				}
			case tcell.KeyLeft:
				if len(path) > 0 {
					parentDir := filepath.Dir(path)
					list.Clear()
					path = parentDir
					addToList(list, listTitle, path)
				}
			case tcell.KeyTab:
				app.SetFocus(listSecond)

			case tcell.KeyEnter:
				list.SetSelectedFocusOnly(true)

			case tcell.KeyRune:
				selectedIndexListOne := list.GetCurrentItem()
				selectedItemListOne, _ := list.GetItemText(selectedIndexListOne)
				selectedIndexListSecond := listSecond.GetCurrentItem()
				selectedItemListSecond, _ := listSecond.GetItemText(selectedIndexListSecond)

				filePathListOne := filepath.Join(path, selectedItemListOne)
				filePathListSecond := filepath.Join(pathSecond, selectedItemListSecond)
				baseDir := filepath.Dir(filePathListSecond)
				destinationPath := filepath.Join(baseDir, selectedItemListOne)

				//Copying files
				if event.Rune() == 'c' {

					resultLabel.SetText(fmt.Sprintf("list1 %s, list2 %s", filePathListOne, destinationPath))

					copy := make(chan bool)

					go func() {
						err := copyFile(filePathListOne, destinationPath)
						if err != nil {
							resultLabel.SetText(fmt.Sprintf("Error copying: %s", err.Error()))
						} else {
							resultLabel.SetText(fmt.Sprintf("Copied successfully from %s to %s", filePathListOne, destinationPath))
							listSecond.Clear()
							addToList(listSecond, listSecondTitle, pathSecond) //refresh  list 2 after files/dir were copied

						}
						copy <- true

					}()
					<-copy
				}

				//Moving files
				if event.Rune() == 'm' {

					move := make(chan bool)

					go func() {
						err := moveFile(filePathListOne, destinationPath)
						if err != nil {
							resultLabel.SetText(fmt.Sprintf("Error moving: %s", err.Error()))
						} else {
							resultLabel.SetText(fmt.Sprintf("Moved successfully from %s to %s", filePathListOne, destinationPath))
							list.Clear()
							addToList(list, listTitle, path) //refresh  list 1 after files/dir were moved
							listSecond.Clear()
							addToList(listSecond, listSecondTitle, pathSecond) //refresh  list 2 after files/dir were moved

						}
						move <- true

					}()
					<-move
				}
				//Delete selected file/directory
				if event.Rune() == 'd' {
					os.Remove(filePathListOne)
				}
				//Get info about file/directory
				if event.Rune() == 'i' {
					fileType, fileFormatted, fileSize, creationTime := displaySingleFileInfo(filePathListOne)
					resultLabel.SetText(fmt.Sprintf("Path: %s\n Type: %s\n Size: %.2f %s\n LastModified: %s", filePathListOne, fileType, fileFormatted, fileSize, creationTime))
				}
				//Stop application
				if event.Rune() == 'q' {
					app.Stop()
				}
			}
		case app.GetFocus() == listSecond:
			// Event handling for list 2
			switch event.Key() {
			case tcell.KeyRight:
				selectedIndex := listSecond.GetCurrentItem()
				selectedItem, _ := listSecond.GetItemText(selectedIndex)
				filePath := filepath.Join(pathSecond, selectedItem)
				fileInfo, err := os.Stat(filePath)

				if err != nil {
					panic(err)
				}

				if fileInfo.IsDir() {
					pathSecond = filePath
					listSecond.Clear()
					addToList(listSecond, listSecondTitle, pathSecond)
				} else {
					fileType, fileFormatted, fileSize, creationTime := displaySingleFileInfo(filePath)
					resultLabel.SetText(fmt.Sprintf("Path: %s\n Type: %s\n Size: %.2f %s\n LastModified: %s", filePath, fileType, fileFormatted, fileSize, creationTime))
				}
			case tcell.KeyLeft:
				if len(path) > 0 {
					parentDir := filepath.Dir(pathSecond)
					listSecond.Clear()
					pathSecond = parentDir
					addToList(listSecond, listSecondTitle, pathSecond)
				}
			case tcell.KeyTab:
				app.SetFocus(list)
			}
		}

		return event
	})
}
