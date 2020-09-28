package main

import (
	"os/exec"
	"strings"
)

func Search(s string, option *Option) ([]string, error) {
	var (
		cmd  string
		args []string
	)
	switch option.SearchMode {
	case Regex:
	case WordMatch:
	case FirstMatch:
		cmd = "git"
		args = []string{
			"grep", s,
		}
	case FuzzyFind:
	}

	result, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(result), "\n"), nil
}
