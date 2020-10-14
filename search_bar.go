package main

func setLabelByMode(mode Mode) {
	switch mode {
	case FirstMatch:
		searchBar.SetLabel("(FM)>>> ")
	case WordMatch:
		searchBar.SetLabel("(WM)>>> ")
	case Regex:
		searchBar.SetLabel("(RX)>>> ")
	case FuzzyFind:
		searchBar.SetLabel("(FF)>>> ")
	}
}
