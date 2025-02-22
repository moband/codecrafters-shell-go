package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func parseCommand(command string) []string {
	var args []string
	var current strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false
	i := 0

	for i < len(command) {

		if command[i] == '\\' && !inSingleQuotes && !inDoubleQuotes && i+1 < len(command) {
			current.WriteByte(command[i+1])
			i += 2
			continue
		}

		if command[i] == '\'' && !inDoubleQuotes {
			inSingleQuotes = !inSingleQuotes
			i++
			continue
		}

		if command[i] == '"' && !inSingleQuotes {
			inDoubleQuotes = !inDoubleQuotes
			i++
			continue
		}

		if command[i] == '\\' && inDoubleQuotes && i+1 < len(command) {
			next := command[i+1]
			if next == '\\' || next == '$' || next == '"' || next == '\n' {
				current.WriteByte(next)
				i += 2
				continue
			}
		}

		if command[i] == ' ' && !inSingleQuotes && !inDoubleQuotes {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		} else {
			current.WriteByte(command[i])
		}
		i++
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

type CommandHandler struct {
	command    string
	args       []string
	outputFile string
	stderrFile string
}

func NewCommandHandler(command string) *CommandHandler {
	args := parseCommand(command)

	var outputFile, stderrFile string
	newArgs := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {

		if (args[i] == ">" || args[i] == "1>") && i+1 < len(args) {
			outputFile = args[i+1]
			i++
		} else if args[i] == "2>" && i+1 < len(args) {
			stderrFile = args[i+1]
			i++
		} else {
			newArgs = append(newArgs, args[i])
		}
	}

	return &CommandHandler{
		command:    command,
		args:       newArgs,
		outputFile: outputFile,
		stderrFile: stderrFile,
	}
}

func (ch *CommandHandler) Execute() {
	if len(ch.args) == 0 {
		return
	}

	var stdout *os.File

	if ch.outputFile != "" {
		stdout = os.Stdout
		file, err := os.Create(ch.outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			return
		}
		os.Stdout = file
		defer func() {
			file.Close()
			os.Stdout = stdout
		}()
	}

	if ch.stderrFile != "" {
		stderr := os.Stderr
		file, err := os.Create(ch.stderrFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			return
		}
		os.Stderr = file
		defer func() {
			file.Close()
			os.Stderr = stderr
		}()
	}

	switch ch.args[0] {
	case "exit":
		ch.handleExit()
	case "echo":
		ch.handleEcho()
	case "pwd":
		ch.handlePwd()
	case "cd":
		ch.handleCd()
	case "type":
		ch.handleType()
	default:
		ch.handleExternal()
	}
}

func (ch *CommandHandler) handleExit() {
	if ch.command == "exit 0" {
		os.Exit(0)
	}
}

func (ch *CommandHandler) handleEcho() {
	if len(ch.args) > 1 {
		fmt.Println(strings.Join(ch.args[1:], " "))
	}
}

func (ch *CommandHandler) handlePwd() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	} else {
		fmt.Println(wd)
	}
}

func (ch *CommandHandler) handleCd() {
	if len(ch.args) <= 1 {
		return
	}

	dir := ch.args[1]
	if dir == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
		dir = homeDir
	}

	if err := os.Chdir(dir); err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", dir)
	}
}

func (ch *CommandHandler) handleType() {
	if len(ch.args) <= 1 {
		return
	}

	cmd := ch.args[1]
	switch cmd {
	case "echo", "exit", "type", "pwd", "cd":
		fmt.Printf("%s is a shell builtin\n", cmd)
	default:
		pathDirs := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
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
}

func (ch *CommandHandler) handleExternal() {
	pathDirs := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	found := false
	for _, dir := range pathDirs {
		executablePath := filepath.Join(dir, ch.args[0])
		if _, err := os.Stat(executablePath); err == nil {
			found = true
			cmd := exec.Command(executablePath, ch.args[1:]...)

			if ch.outputFile != "" {
				file, err := os.Create(ch.outputFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
					return
				}
				defer file.Close()
				cmd.Stdout = file
			} else {
				cmd.Stdout = os.Stdout
			}

			if ch.stderrFile != "" {
				file, err := os.Create(ch.stderrFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
					return
				}
				defer file.Close()
				cmd.Stderr = file
			} else {
				cmd.Stderr = os.Stderr
			}

			cmd.Args[0] = ch.args[0]
			cmd.Run()
			break
		}
	}
	if !found {
		fmt.Println(ch.command + ": command not found")
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		command = command[:len(command)-1]

		handler := NewCommandHandler(command)
		handler.Execute()
	}
}
