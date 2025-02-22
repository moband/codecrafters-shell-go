package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		// Uncomment this block to pass the first stage
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		command = command[:len(command)-1]

		if command == "exit 0" {
			os.Exit(0)
		}

		if len(command) > 5 && command[:5] == "echo " {
			fmt.Println(command[5:])
			continue
		}

		if len(command) > 5 && command[:5] == "type " {
			cmd := command[5:]
			switch cmd {
			case "echo", "exit", "type":
				fmt.Printf("%s is a shell builtin\n", cmd)
			default:
				pathDirs := strings.Split(os.Getenv("PATH"), ":")
				found := false
				for _, dir := range pathDirs {
					executablePath := filepath.Join(dir, cmd)
					if _, err := os.Stat(executablePath); err == nil {
						fmt.Printf("%s is %s\n", cmd, executablePath)
						found = true
						break
					}
				}
				if !found {
					fmt.Printf("%s: not found\n", cmd)
				}
			}
			continue
		}

		fmt.Println(command + ": command not found")
	}
}
