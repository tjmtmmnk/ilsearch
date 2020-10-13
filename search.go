package main

import (
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/alessio/shellescape.v1"
)

type SearchResult struct {
	index    int
	lineNum  int
	fileName string
	text     string
}

func isValidQuery(q string) bool {
	empty := q == ""
	matched, err := regexp.Match(`\\$`, []byte(q))
	if err != nil {
		return false
	}
	return !empty && !matched
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
		cmd = []string{
			"git", "grep", "-wHn", q,
		}
	case WordMatchIgnoreCase:
		cmd = []string{
			"git", "grep", "-wiHn", q,
		}
	case FirstMatch:
		cmd = []string{
			"git", "grep", "-Hn", q,
		}
	case FuzzyFind:
		target := filepath.Join(context.workDir, "*")
		// to use glob pattern
		c := "agrep -n " + shellescape.Quote(q) + " " + target
		cmd = []string{
			"/bin/sh", "-c", c,
		}
	}

	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if len(out) == 0 {
		return []SearchResult{}, nil
	}

	if err != nil {
		return []SearchResult{}, err
	}

	var results []SearchResult
	for i, result := range strings.Split(string(out), "\n") {
		ignore, fileName, lineNum := splitResult(result, option.Command)
		if !ignore {
			results = append(results, SearchResult{
				index:    i,
				lineNum:  lineNum,
				fileName: fileName,
				text:     result,
			})
		}
	}
	return results, nil
}

func splitResult(result string, command Command) (ignore bool, fileName string, lineNum int) {
	splitted := strings.Split(result, ":")
	if len(splitted) < 2 {
		ignore = true
		return
	}
	switch command {
	case GitGrep:
		fileName = splitted[0]
		lineNum, _ = strconv.Atoi(splitted[1])
	case Agrep:
		fileName = splitted[0]
		lineNum, _ = strconv.Atoi(strings.Replace(splitted[1], " ", "", -1))
	}
	return
}
