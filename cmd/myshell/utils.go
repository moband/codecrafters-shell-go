package main

import (
	"os"
	"path/filepath"
	"strings"
)

func findExecutable(cmd string) (string, bool) {
	pathDirs := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	for _, dir := range pathDirs {
		executablePath := filepath.Join(dir, cmd)
		if _, err := os.Stat(executablePath); err == nil {
			return executablePath, true
		}
	}
	return "", false
}

func createFile(path string, appendMode bool) (*os.File, error) {
	flag := os.O_WRONLY | os.O_CREATE
	if appendMode {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}
	return os.OpenFile(path, flag, 0644)
}
