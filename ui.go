package main

import (
	"github.com/rivo/tview"
)

var (
	  list = tview.NewList()

    listSecond = tview.NewList()

		resultLabel = tview.NewTextView().
      SetTextAlign(tview.AlignCenter).
      SetText("")

    informationLabel = tview.NewTextView().
      SetTextAlign(tview.AlignLeft).
      SetText("(c) copy (m) move (i) info (q) quit")
    
    listTitle = tview.NewTextView().
      SetText("List One").
      SetTextAlign(tview.AlignCenter)
    
    listSecondTitle = tview.NewTextView().
      SetText("List Two").
      SetTextAlign(tview.AlignCenter)
    

)

func setupUI() *tview.Grid {

    list.SetBorder(true).SetTitle("Source")
    listSecond.SetBorder(true).SetTitle("Destination")

	  grid := tview.NewGrid().
        SetRows(3, 0, 5).
        SetColumns(60, 0, 60).
        SetBorders(true).
        AddItem(listTitle, 0, 0, 1, 1, 0, 0, false).
        AddItem(listSecondTitle, 0, 2, 1, 1, 0, 0, false).
        AddItem(list, 1, 0, 1, 1, 0, 0, true).
        AddItem(listSecond, 1, 2, 1, 1, 0, 0, true).
        AddItem(informationLabel, 2, 0, 1, 1, 0, 0, false).
        AddItem(resultLabel, 2, 1, 1, 2, 0, 0, false)
       
	return grid
}
