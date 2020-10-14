package main

import (
	"log"

	"github.com/rivo/tview"
)

type CheckContent int

const (
	CFirstMatch CheckContent = iota
	CWordMatch
	CRegex
	CFuzzyFind
	CGitGrep
	CRipGrep
)

type CheckboxWithContent struct {
	c       *tview.Checkbox
	content CheckContent
}
type CheckboxContainer struct {
	items []CheckboxWithContent
}

func (cc CheckboxContainer) Check(index int, newState bool) {
	if len(cc.items) <= index {
		log.Fatal("checkbox out of index")
	}
	cc.items[index].c.SetChecked(newState)
	for i, item := range cc.items {
		if i != index {
			item.c.SetChecked(false)
		}
	}
}

func (cc CheckboxContainer) CheckedItem() *CheckboxWithContent {
	for _, item := range cc.items {
		if item.c.IsChecked() {
			return &item
		}
	}
	return nil
}

func searchOptionCheckbox() CheckboxContainer {
	items := []CheckboxWithContent{
		{
			c:       tview.NewCheckbox().SetLabel("FirstMatch"),
			content: CFirstMatch,
		},
		{
			c:       tview.NewCheckbox().SetLabel("WordMatch"),
			content: CWordMatch,
		},
		{
			c:       tview.NewCheckbox().SetLabel("Regex"),
			content: CRegex,
		},
		{
			c:       tview.NewCheckbox().SetLabel("FuzzyFind"),
			content: CFuzzyFind,
		},
	}

	cc := CheckboxContainer{
		items: items,
	}

	for i, item := range items {
		item.c.SetChangedFunc(func(checked bool) {
			cc.Check(i, checked)
			if checked {
				mode := FirstMatch
				command := GitGrep
				switch item.content {
				case CFirstMatch:
					mode = FirstMatch
				case CWordMatch:
					mode = WordMatch
				case CRegex:
					mode = Regex
				case CFuzzyFind:
					mode = FuzzyFind
					command = Agrep
				}
				switch item.content {
				case CGitGrep:
					command = GitGrep
				case CRipGrep:
					command = RipGrep
				}
				option.SearchMode = mode
				option.Command = command
				setLabelByMode(mode)
			}
		})
	}

	return cc
}
