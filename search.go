package main

import (
	"os/exec"
	"strings"
)

func Search(s string, option *Option) ([]string, error) {
	var (
		cmd []string
	)
	switch option.SearchMode {
	case Regex:
	case WordMatch:
	case FirstMatch:
		cmd = []string{
			"git", "grep", "-Hn", s,
		}
	case FuzzyFind:
	}

	result, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(result), "\n"), nil
}
