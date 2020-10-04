package main

import (
	"os/exec"
	"strconv"
	"strings"
)

type SearchResult struct {
	index    int
	lineNum  int
	fileName string
	text     string
}

func Search(s string, option *Option) ([]SearchResult, error) {
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

	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return []SearchResult{}, err
	}

	var results []SearchResult
	for i, result := range strings.Split(string(out), "\n") {
		splittedResult := strings.Split(result, ":")
		if len(splittedResult) < 2 {
			continue
		}
		fileName := splittedResult[0]
		lineNum, _ := strconv.Atoi(splittedResult[1])
		results = append(results, SearchResult{
			index:    i,
			lineNum:  lineNum,
			fileName: fileName,
			text:     result,
		})
	}
	return results, nil
}
