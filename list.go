package main

import "github.com/rivo/tview"

type Item struct {
	Text string
}

func UpdateList(l *tview.List, items []string) {
	l.Clear()
	for _, item := range items {
		l.AddItem(item, "", 0, nil)
	}
}
