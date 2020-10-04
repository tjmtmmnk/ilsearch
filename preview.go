package main

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/gdamore/tcell"
)

func GetPreviewContent(s tcell.Screen, fileName string, highlightLine int, theme string) (string, error) {
	_, h := s.Size()
	from := highlightLine - (h/2 - 1)
	if from < 0 {
		from = 0
	}
	to := highlightLine + (h/2 + 1)
	lineRange := fmt.Sprintf("%d:%d", from, to)
	cmd := []string{"bat", "--line-range", lineRange, "--highlight-line", strconv.Itoa(highlightLine), "--color=always", "--theme", theme, "--style=numbers", fileName}

	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	return string(out), err
}
