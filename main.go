package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	s, e := tcell.NewScreen()
	if e != nil {
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		os.Exit(1)
	}
	app = tview.NewApplication().SetScreen(s)
	list = tview.NewList().ShowSecondaryText(false)
	//pages := tview.NewPages()

	searchBar = tview.NewInputField().
		SetLabel(">>> ").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite)

	preview = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetRegions(true)

	option := &Option{
		SearchMode: FirstMatch,
		Directory:  "./",
	}

	searchBar.SetChangedFunc(func(text string) {
		results, err := Search(text, option)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		UpdateList(list, results)
	})

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		currentItemText, _ := list.GetItemText(index)
		sp := strings.Split(currentItemText, ":")
		if len(sp) != 3 {
			panic("don't match format")
		}
		fileName := sp[0]
		lineNum, _ := strconv.Atoi(sp[1])

		text, err := GetPreviewContent(s, fileName, lineNum, "OneHalfDark")
		text = tview.TranslateANSI(text)
		if err != nil {
			panic(err)
		}
		preview.Clear().SetText(text)
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBar, 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(list, 0, 1, false).
			AddItem(preview, 0, 1, false), 0, 10, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
