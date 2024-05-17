package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	CopyKey   = 'c'
	MoveKey   = 'm'
	DeleteKey = 'd'
	InfoKey   = 'i'
	QuitKey   = 'q'
	TabKey    = tcell.KeyTab
	EnterKey  = tcell.KeyEnter
	LeftKey   = tcell.KeyLeft
	RightKey  = tcell.KeyRight
)

var (
	selectedItems []string
	dstPath       []string
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

func handleInfoKey(path string, label *tview.TextView) {
	if list.GetItemCount() != 0 {
		selectedIndexListOne := list.GetCurrentItem()
		selectedItemListOne, _ := list.GetItemText(selectedIndexListOne)
		filePathListOne := filepath.Join(path, selectedItemListOne)

		fileType, formattedSize, sizeUnit, creationTime, mode, err := displaySingleFileInfo(filePathListOne)
		if err != nil {
			label.SetText(fmt.Sprintf("Error:", err))
		}
		label.SetText(fmt.Sprintf("Path: %s\n Type: %s\n Size: %.2f %s\n LastModified: %s \n Permissions: %s", filePathListOne, fileType, formattedSize, sizeUnit, creationTime, mode))
	}
}

func handleEnterKey(list *tview.List, label *tview.TextView, path string) {
	if list.GetItemCount() != 0 {
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
	} else {
		label.SetText("No items in list")
	}
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
		label.SetText(fmt.Sprintf("Moved successfully from %s to %s\n", selectedItems, destinationPath))

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
	switch r {

	case CopyKey:
		dstPath = nil
		handleCopy(selectedItems, appendToSlice(dstPath, selectedItems, pathSecond), list, listSecond, label, titleFirst, titleSecond, path, pathSecond)
		selectedItems = nil

	case MoveKey:
		dstPath = nil
		handleMove(selectedItems, appendToSlice(dstPath, selectedItems, pathSecond), list, listSecond, label, titleFirst, titleSecond, path, pathSecond)
		selectedItems = nil

	case DeleteKey:
		for _, items := range selectedItems {
			os.RemoveAll(items)
		}
		list.Clear()
		addToList(list, listTitle, path)

	case InfoKey:
		handleInfoKey(path, label)

	case QuitKey:
		app.Stop()
	}
}

func handleRightKey(event *tcell.EventKey, list *tview.List, path string, title, resultLabel *tview.TextView) (string, *tcell.EventKey) {
	if list.GetItemCount() != 0 {
		selectedIndex := list.GetCurrentItem()
		selectedItem, _ := list.GetItemText(selectedIndex)

		filePath := filepath.Join(path, selectedItem)
		fileInfo, err := os.Stat(filePath)

		if err != nil {
			if os.IsPermission(err) {
				resultLabel.SetText(fmt.Sprintf("Access denied for: %s", filePath))
				return path, event
			} else if os.IsNotExist(err) {
				resultLabel.SetText("Doesn't exist")
				return path, event
			}
			panic(err)
		}

		if fileInfo.IsDir() {
			path = filePath
			list.Clear()
			if err := addToList(list, title, path); err != nil {
				resultLabel.SetText(fmt.Sprintf("Error: %s", err))
			}
			// clear slices when moving through directories
			selectedItems = nil
			dstPath = nil

		} else {
			resultLabel.SetText("File")
		}
	}
	return path, event
}

func handleLeftKey(event *tcell.EventKey, list *tview.List, path string, title, resultLabel *tview.TextView) (string, *tcell.EventKey) {
	if len(path) > 0 {
		parentDir := filepath.Dir(path)
		list.Clear()
		path = parentDir
		addToList(list, title, path)

		// clear slices when moving through directories
		selectedItems = nil
		dstPath = nil

		// clear label when moving through directories
		resultLabel.SetText("")
	}
	return path, event
}

func eventHandler(app *tview.Application, list, listSecond *tview.List, pathSecond, path string, resultLabel *tview.TextView) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case app.GetFocus() == list:
			// Event handling for list 1
			switch event.Key() {
			case RightKey:
				path, event = handleRightKey(event, list, path, listTitle, resultLabel)
			case LeftKey:
				path, event = handleLeftKey(event, list, path, listTitle, resultLabel)
			case TabKey:
				selectedItems = nil
				app.SetFocus(listSecond)
			case EnterKey:
				handleEnterKey(list, resultLabel, path)
			case tcell.KeyRune:
				handleRuneKey(app, event.Rune(), list, listSecond, listTitle, listSecondTitle, resultLabel, path, pathSecond)
			}
		case app.GetFocus() == listSecond:
			// Event handling for list 2
			switch event.Key() {
			case tcell.KeyRight:
				pathSecond, event = handleRightKey(event, listSecond, pathSecond, listSecondTitle, resultLabel)
			case tcell.KeyLeft:
				pathSecond, event = handleLeftKey(event, listSecond, pathSecond, listSecondTitle, resultLabel)
			case tcell.KeyTab:
				selectedItems = nil
				app.SetFocus(list)
			case tcell.KeyEnter:
				handleEnterKey(listSecond, resultLabel, pathSecond)
			case tcell.KeyRune:
				handleRuneKey(app, event.Rune(), listSecond, list, listSecondTitle, listTitle, resultLabel, pathSecond, path)
			}
		}

		return event
	})
}
