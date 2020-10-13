package main

type Mode int
type Command int

const (
	Regex Mode = iota
	FirstMatch
	WordMatch
	WordMatchIgnoreCase
	FuzzyFind
)

const (
	GitGrep Command = iota
	RipGrep
	Agrep
)

type Option struct {
	SearchMode Mode
	Directory  string
	Command    Command
}
