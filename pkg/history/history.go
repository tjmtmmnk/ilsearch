package history

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
)

const (
	Directory      = "/.ilsearch"
	FileName       = "ilsearch_history.json"
	MaxDirectories = 5
)

var homeDir string

type History struct {
	FocusDirectories []string `json:"focus_directories"`
	LastQuery        string   `json:"last_query"`
}

type Query struct {
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
	return homeDir + Directory + "/" + FileName
}

func Init() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	homeDir = usr.HomeDir
	path := homeDir + Directory
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

func updateHistory(history *History, query Query) (*History, error) {
	switch query.Key {
	case "FocusDirectories":
		if len(history.FocusDirectories) > MaxDirectories-1 {
			history.FocusDirectories = history.FocusDirectories[1:]
		}
		history.FocusDirectories = append(history.FocusDirectories, query.Val)
	case "LastQuery":
		history.LastQuery = query.Val
	default:
		return nil, errors.New("invalid json key")
	}
	return history, nil
}

func WriteHistories(queries []Query) error {
	history, err := ReadHistory()
	if err != nil {
		return err
	}

	var updatedHistory *History
	for _, q := range queries {
		updatedHistory, err = updateHistory(history, q)
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

func WriteHistory(query Query) error {
	return WriteHistories([]Query{query})
}

func ReadHistory() (*History, error) {
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
