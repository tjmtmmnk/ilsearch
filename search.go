package main

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type SearchResult struct {
	index    int
	lineNum  int
	fileName string
	text     string
}

func isValidQuery(q string) bool {
	matched, err := regexp.Match(`\\$`, []byte(q))
	if err != nil {
		return false
	}
	return !matched
}

func Search(q string, option *Option) ([]SearchResult, error) {
	var (
		cmd []string
	)
	if !isValidQuery(q) {
		return []SearchResult{}, nil
	}
	switch option.SearchMode {
	case Regex:
		cmd = []string{
			"git", "grep", "-EHn", q,
		}
	case WordMatch:
	case FirstMatch:
		cmd = []string{
			"git", "grep", "-Hn", q,
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
