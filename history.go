package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
)

const (
	DIRECTORY = "/.ilsearch"
	FILENAME  = "ilsearch_history.json"
)

var homeDir string

type History struct {
	FocusDirectories []string `json:"focus_directories"`
	LastQuery        string   `json:"last_query"`
}

type HistoryQuery struct {
	Key string
	Val string
}

func newHistory() *History {
	return &History{
		FocusDirectories: []string{},
		LastQuery:        "",
	}
}

func getFilePath() string {
	return homeDir + DIRECTORY + "/" + FILENAME
}

func InitHistory() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	homeDir = usr.HomeDir
	path := homeDir + DIRECTORY
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0744)

		if err != nil {
			return err
		}
	}
	filePath := getFilePath()
	_, err = os.Stat(filePath)
	existFile := err == nil
	if !existFile {
		emptyHistory := newHistory()
		bytes, err := json.Marshal(emptyHistory)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filePath, bytes, 0644)
	}
	return nil
}

func updateHistory(history *History, key, s string) (*History, error) {
	switch key {
	case "FocusDirectories":
		if len(history.FocusDirectories) > 5 {
			history.FocusDirectories = history.FocusDirectories[1:]
		}
		history.FocusDirectories = append(history.FocusDirectories, s)
	case "LastQuery":
		history.LastQuery = s
	default:
		return nil, errors.New("invalid json key")
	}
	return history, nil
}

func writeHistories(queries []HistoryQuery) error {
	history, err := readHistory()
	if err != nil {
		return err
	}

	var updatedHistory *History
	for _, q := range queries {
		updatedHistory, err = updateHistory(history, q.Key, q.Val)
		if err != nil {
			return err
		}
	}

	bytes, err := json.Marshal(updatedHistory)
	if err != nil {
		return err
	}
	filePath := getFilePath()
	err = ioutil.WriteFile(filePath, bytes, 0644)
	return nil
}

func writeHistory(query HistoryQuery) error {
	return writeHistories([]HistoryQuery{query})
}

func readHistory() (*History, error) {
	filePath := getFilePath()
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	history := newHistory()
	err = json.Unmarshal(bytes, history)
	if err != nil {
		return nil, err
	}
	return history, nil
}
