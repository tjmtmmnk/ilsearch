package main

type Mode int

const (
	Regex Mode = iota
	FirstMatch
	WordMatch
	FuzzyFind
)

type Option struct {
	SearchMode Mode
	Directory  string
}
