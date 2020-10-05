package main

type Mode int

const (
	Regex Mode = iota
	FirstMatch
	WordMatch
	WordMatchIgnoreCase
	FuzzyFind
)

type Option struct {
	SearchMode Mode
	Directory  string
}
