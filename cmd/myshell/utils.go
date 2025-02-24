package main

// import (
// 	"os"
// 	"path/filepath"
// 	"strings"
// )

// func findExecutable(cmd string) (string, bool) {
// 	pathDirs := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
// 	for _, dir := range pathDirs {
// 		executablePath := filepath.Join(dir, cmd)
// 		if _, err := os.Stat(executablePath); err == nil {
// 			return executablePath, true
// 		}
// 	}
// 	return "", false
// }
