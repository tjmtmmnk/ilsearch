package main

import (
	"log"
	"os"
	"path/filepath"
)

type Context struct {
	workDir string
}

func NewContext() *Context {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	evaledWorkDir, err := filepath.EvalSymlinks(workDir)
	if err != nil {
		log.Fatal(err)
	}
	return &Context{
		workDir: evaledWorkDir,
	}
}
