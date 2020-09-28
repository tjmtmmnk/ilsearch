package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	app       *tview.Application
	list      *tview.List
	searchBar *tview.InputField
	preview   *tview.TextView
)

func main() {
	app = tview.NewApplication()
	list = tview.NewList().ShowSecondaryText(false)
	//pages := tview.NewPages()

	searchBar = tview.NewInputField().
		SetLabel(">>> ").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite)

	preview = tview.NewTextView()

	option := &Option{
		SearchMode: FirstMatch,
		Directory:  "./",
	}
	results, err := Search("new", option)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, r := range results {
		list.AddItem(r, "", 0, nil)
	}
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBar, 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(list, 0, 1, false).
			AddItem(preview.SetBorder(true), 0, 1, false), 0, 10, false)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
