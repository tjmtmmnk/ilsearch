package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	context       *Context
	app           *tview.Application
	list          *tview.List
	searchBar     *tview.InputField
	preview       *tview.TextView
	searchResults []SearchResult
)

func main() {
	context = NewContext()
	s, e := tcell.NewScreen()
	if e != nil {
		log.Fatal("screen error")
	}
	if e := s.Init(); e != nil {
		log.Fatal("screen init error")
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
		SetScrollable(true)

	option := &Option{
		SearchMode: WordMatch,
		Directory:  "./",
	}

	searchBar.SetChangedFunc(func(text string) {
		var err error
		searchResults, err = Search(text, option)
		if err != nil {
			fmt.Println(err)
			log.Fatal("search error")
		}
		var resultTexts []string
		for _, result := range searchResults {
			resultTexts = append(resultTexts, result.text)
		}
		UpdateList(list, resultTexts)
	})

	searchBar.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			app.SetFocus(list)
		case tcell.KeyEsc:
			app.Stop()
		}
	})

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		item := searchResults[index]
		if index != item.index {
			log.Fatal("not match index")
		}

		text, err := GetPreviewContent(s, item.fileName, item.lineNum, "OneHalfDark")
		text = tview.TranslateANSI(text)
		if err != nil {
			panic(err)
		}
		preview.Clear().SetText(text)
	})

	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		item := searchResults[index]
		if index != item.index {
			log.Fatal("not match index")
		}
		OpenFile(context.workDir, item.fileName, item.lineNum)
	})

	list.SetDoneFunc(func() {
		app.SetFocus(searchBar)
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
