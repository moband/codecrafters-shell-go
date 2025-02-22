package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

		if command == "pwd" {
			wd, err := os.Getwd()
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			} else {
				fmt.Println(wd)
			}
			continue
		}

		if len(command) > 3 && command[:3] == "cd " {
			dir := command[3:]
			if err := os.Chdir(dir); err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", dir)
			}
			continue
		}

		if len(command) > 5 && command[:5] == "type " {
			cmd := command[5:]
			switch cmd {
			case "echo", "exit", "type", "pwd", "cd":
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

		args := strings.Fields(command)
		if len(args) > 0 {
			cmdName := args[0]
			pathDirs := strings.Split(os.Getenv("PATH"), ":")
			found := false
			for _, dir := range pathDirs {
				executablePath := filepath.Join(dir, cmdName)
				if _, err := os.Stat(executablePath); err == nil {
					found = true
					cmd := exec.Command(executablePath, args[1:]...)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Args[0] = cmdName
					if err := cmd.Run(); err != nil {
						fmt.Fprint(os.Stderr, err)
					}
					break
				}
			}
			if !found {
				fmt.Println(command + ": command not found")
			}
		}
	}
}
