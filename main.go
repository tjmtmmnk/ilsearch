package main

import (
	"log"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/tjmtmmnk/ilsearch/pkg/history"
)

var (
	context       *Context
	app           *tview.Application
	list          *tview.List
	searchBar     *tview.InputField
	preview       *tview.TextView
	searchResults []SearchResult
	option        *Option
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := history.Init()
	err = history.WriteHistory(history.Query{
		Key: "FocusDirectories",
		Val: "/hoge",
	})
	if err != nil {
		log.Fatal(err)
	}
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

	searchBar = tview.NewInputField().
		SetLabel(">>> ").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite)

	preview = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)

	option = &Option{
		SearchMode: WordMatch,
		Directory:  "./",
		Command:    GitGrep,
	}

	searchBar.SetChangedFunc(func(text string) {
		var err error
		searchResults, err = Search(text, option)
		if err != nil {
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
	checkbox := searchOptionCheckbox()

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBar, 0, 3, true).SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(list, 0, 1, false).
			AddItem(preview, 0, 1, false), 0, 30, false)

	tree := InitTree(context.workDir)

	searchOption := tview.NewFlex().SetDirection(tview.FlexRow)
	for _, item := range checkbox.items {
		searchOption = searchOption.AddItem(item.c, 2, 0, false)
	}

	pages := tview.NewPages().
		AddPage("app", flex, true, true).
		AddPage("tree", tree, true, false).
		AddPage("option", searchOption, true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch key := event.Key(); {
		case key == tcell.KeyCtrlS:
			pages.SwitchToPage("tree")
		case key == tcell.KeyRune || key == tcell.KeyBackspace || key == tcell.KeyBackspace2:
			app.SetFocus(searchBar)
		case key == tcell.KeyDown && searchBar.HasFocus():
			app.SetFocus(list)
		case key == tcell.KeyCtrlF:
			pages.SwitchToPage("option")
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
