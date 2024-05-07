package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	selectedItems []string
	dstPath []string
)

func addToList(target *tview.List, label *tview.TextView, path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
			target.AddItem(file.Name(), "", 0, nil)
	}
	itemCount := target.GetItemCount()
	label.SetText(fmt.Sprintf("Path: %s\n Items: %d", path, itemCount))
	return nil
}

func appendToSlice(sl, sourceItems []string, baseDir string) []string {
		for _, items := range sourceItems {
		itms := filepath.Join(baseDir, filepath.Base(items))
		sl = append(sl, itms)
	}
	return sl
}

func removeFromSlice(slice []string, text string) []string {
	for i, v := range slice {
		if v == text {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}


func handleEnterKey(list *tview.List, label *tview.TextView, path string) {
	selectedIndexListOne := list.GetCurrentItem()
	selectedItemListOne, secondaryTextListOne := list.GetItemText(selectedIndexListOne)
	filePathListOne := filepath.Join(path, selectedItemListOne)

	if secondaryTextListOne == "" {
		list.SetItemText(selectedIndexListOne, selectedItemListOne, "selected")
		selectedItems = append(selectedItems, filePathListOne)
	} else {
		list.SetItemText(selectedIndexListOne, selectedItemListOne, "")
		selectedItems = removeFromSlice(selectedItems, filePathListOne)
	}
	text := strings.Join(selectedItems, "\n")
	label.SetText("Selected items: " + text)
}

func handleCopy(selectedItems []string, destinationPath []string, list, listSecond *tview.List, label, titleFirst, titleSecond *tview.TextView, path, pathSecond string) {
	copy := make(chan bool)
	go func() {
		err := copyFiles(selectedItems, destinationPath)
		if err != nil {
			label.SetText(fmt.Sprintf("Error copying: %s", err.Error()))
			copy <- false
			return
		}
		 label.SetText(fmt.Sprintf("Copied successfully from %s to %s", selectedItems, destinationPath))

		 	list.Clear()
		  if err := addToList(list, titleFirst, path); err != nil {
        label.SetText(fmt.Sprintf("Error updating list 1: %s", err))
        copy <- false
        return
      }
				listSecond.Clear()
      if err := addToList(listSecond, titleSecond, pathSecond); err != nil {
        label.SetText(fmt.Sprintf("Error updating list 2: %s", err))
        copy <- false
        return
      }       
		
		copy <- true
	}()
	<-copy
}

func handleMove(selectedItems []string, destinationPath []string, list, listSecond *tview.List, label, titleFirst, titleSecond *tview.TextView, path, pathSecond string) {
    move := make(chan bool)
    go func() {
        err := moveFiles(selectedItems, destinationPath)
        if err != nil {
            label.SetText(fmt.Sprintf("Error moving: %s", err.Error()))
						move <- false
						return
        }
          label.SetText(fmt.Sprintf("Moved successfully from %s to %s", selectedItems, destinationPath))

				list.Clear()	
        if err := addToList(list, titleFirst, path); err != nil {
            label.SetText(fmt.Sprintf("Error updating list 1: %s", err))
            move <- false
            return
        }
				listSecond.Clear()
        if err := addToList(listSecond, titleSecond, pathSecond); err != nil {
            label.SetText(fmt.Sprintf("Error updating list 2: %s", err))
            move <- false
            return
        }         
        move <- true
    }()
    <-move
}

func handleRuneKey(app *tview.Application, r rune, list, listSecond *tview.List, titleFirst, titleSecond, label *tview.TextView, path, pathSecond string) {
	selectedIndexListOne := list.GetCurrentItem()
	selectedItemListOne, _ := list.GetItemText(selectedIndexListOne)
	selectedIndexListSecond := listSecond.GetCurrentItem()
	selectedItemListSecond, _ := listSecond.GetItemText(selectedIndexListSecond)

	filePathListOne := filepath.Join(path, selectedItemListOne)
	filePathListSecond := filepath.Join(pathSecond, selectedItemListSecond)
	baseDir := filepath.Dir(filePathListSecond)
	//destinationPath := filepath.Join(baseDir, selectedItemListOne)
	

	switch r {
	case 'c':
		dstPath = nil
		handleCopy(selectedItems, appendToSlice(dstPath, selectedItems, baseDir), list, listSecond, label, titleFirst, titleSecond, path, pathSecond)
	case 'm':
		dstPath = nil
		handleMove(selectedItems, appendToSlice(dstPath, selectedItems, baseDir), list, listSecond, label, titleFirst, titleSecond, path, pathSecond)
	case 'd':
		os.Remove(filePathListOne)
	case 'i':
		fileType, fileFormatted, fileSize, modifiedTime := displaySingleFileInfo(filePathListOne)
		label.SetText(fmt.Sprintf("Path: %s\n Type: %s\n Size: %.2f %s\n LastModified: %s", filePathListOne, fileType, fileFormatted, fileSize, modifiedTime))
	case 'q':
		app.Stop()
	}
}

func eventHandler(app *tview.Application, list, listSecond *tview.List, pathSecond, path string, resultLabel *tview.TextView) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case app.GetFocus() == list:
			// Event handling for list 1
			switch event.Key() {
			case tcell.KeyRight:
				if list.GetItemCount() != 0 {
				selectedIndex := list.GetCurrentItem()
				selectedItem, _ := list.GetItemText(selectedIndex)
				filePath := filepath.Join(path, selectedItem)
				fileInfo, err := os.Stat(filePath)

				if err != nil {
					if os.IsPermission(err){
						resultLabel.SetText(fmt.Sprintf("Access denied for: %s", filePath))
						return event
					} else if os.IsNotExist(err) {
						resultLabel.SetText("Doesn't exist")
						return event
					}
				panic(err)
				}

				if fileInfo.IsDir() {
					path = filePath
					list.Clear()
					err := addToList(list, listTitle, path)
					if err != nil {
						resultLabel.SetText(fmt.Sprintf("Error: %s", err))
					}
					// clear slice when moving through directories
					selectedItems = nil
					dstPath = nil

				} else {
					fileType, fileFormatted, fileSize, creationTime := displaySingleFileInfo(filePath)
					resultLabel.SetText(fmt.Sprintf("Path: %s\n Type: %s\n Size: %.2f %s\n LastModified: %s", filePath, fileType, fileFormatted, fileSize, creationTime))
				}
			}
			case tcell.KeyLeft:
				if len(path) > 0 {
					parentDir := filepath.Dir(path)
					list.Clear()
					path = parentDir
					addToList(list, listTitle, path)

					// clear slice when moving through directories
					selectedItems = nil
					dstPath = nil
					// clear label when moving through directories
					resultLabel.SetText("")
				}
			case tcell.KeyTab:
				app.SetFocus(listSecond)

			case tcell.KeyEnter:
				handleEnterKey(list, resultLabel, path)

			case tcell.KeyRune:
				//Copying files
				if event.Rune() == 'c' {
					handleRuneKey(app, 'c', list, listSecond, listTitle, listSecondTitle, resultLabel, path, pathSecond)
				}

				//Moving files
				if event.Rune() == 'm' {
					handleRuneKey(app, 'm', list, listSecond, listTitle, listSecondTitle, resultLabel, path, pathSecond)

				}
				//Delete selected file/directory
				if event.Rune() == 'd' {
					handleRuneKey(app, 'd', list, listSecond, listTitle, listSecondTitle, resultLabel, path, pathSecond)

				}
				//Get info about file/directory
				if event.Rune() == 'i' {
					handleRuneKey(app, 'i', list, listSecond, listTitle, listSecondTitle, resultLabel, path, pathSecond)

				}
				//Stop application
				if event.Rune() == 'q' {
					handleRuneKey(app, 'q', list, listSecond, listTitle, listSecondTitle, resultLabel, path, pathSecond)
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
