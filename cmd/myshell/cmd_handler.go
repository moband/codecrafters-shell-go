package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

type CommandHandler struct {
	command      string
	args         []string
	outputFile   string
	stderrFile   string
	appendMode   bool
	stderrAppend bool
}

func NewCommandHandler(command string) *CommandHandler {
	parser := NewParser(command)
	args := parser.Parse()

	var outputFile, stderrFile string
	var appendMode bool
	var stderrAppend bool
	newArgs := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {

		if (args[i] == ">>" || args[i] == "1>>") && i+1 < len(args) {
			outputFile = args[i+1]
			appendMode = true
			i++

		} else if (args[i] == ">" || args[i] == "1>") && i+1 < len(args) {
			outputFile = args[i+1]
			i++
		} else if args[i] == "2>>" && i+1 < len(args) {
			stderrFile = args[i+1]
			stderrAppend = true
			i++
		} else if args[i] == "2>" && i+1 < len(args) {
			stderrFile = args[i+1]
			i++
		} else {
			newArgs = append(newArgs, args[i])
		}
	}

	return &CommandHandler{
		command:      command,
		args:         newArgs,
		outputFile:   outputFile,
		stderrFile:   stderrFile,
		appendMode:   appendMode,
		stderrAppend: stderrAppend,
	}
}

func (ch *CommandHandler) Execute(shell *Shell) {
	if len(ch.args) == 0 {
		return
	}

	var stdout, stderr *os.File

	if ch.outputFile != "" {
		stdout = os.Stdout
		file, err := createFile(ch.outputFile, ch.appendMode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n\r", err)
			return
		}
		os.Stdout = file
		defer func() {
			file.Close()
			os.Stdout = stdout
		}()
	}

	if ch.stderrFile != "" {
		stderr = os.Stderr
		file, err := createFile(ch.stderrFile, ch.stderrAppend)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n\r", err)
			return
		}
		os.Stderr = file
		defer func() {
			file.Close()
			os.Stderr = stderr
		}()
	}

	if handler, ok := shell.builtins[ch.args[0]]; ok {
		handler(ch)
	} else {
		ch.handleExternal()
	}
}

func (ch *CommandHandler) handleExit() {
	term.Restore(int(os.Stdin.Fd()), oldState)
	os.Exit(0)
}

func (ch *CommandHandler) handleEcho() {
	if len(ch.args) > 1 {
		fmt.Fprintf(os.Stdout, "%s\n\r", strings.Join(ch.args[1:], " "))
	}
}

func (ch *CommandHandler) handlePwd() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n\r", err)
	} else {
		fmt.Fprintf(os.Stdout, "%s\n\r", wd)
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
			fmt.Fprintf(os.Stderr, "%v\n\r", err)
			return
		}
		dir = homeDir
	}

	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n\r", dir)
	}
}

func (ch *CommandHandler) handleType() {
	if len(ch.args) <= 1 {
		return
	}

	cmd := ch.args[1]
	switch cmd {
	case "echo", "exit", "type", "pwd", "cd":
		fmt.Fprintf(os.Stdout, "%s is a shell builtin\n\r", cmd)
	default:
		if execPath, found := findExecutable(cmd); found {
			fmt.Fprintf(os.Stdout, "%s is %s\n\r", cmd, execPath)
		} else {
			fmt.Fprintf(os.Stderr, "%s: not found\n\r", cmd)
		}
	}
}

func (ch *CommandHandler) handleExternal() {
	if execPath, found := findExecutable(ch.args[0]); found {
		cmd := exec.Command(execPath, ch.args[1:]...)

		cmd.Stderr = &lineWriter{w: os.Stderr}
		cmd.Stdout = &lineWriter{w: os.Stdout}
		cmd.Args[0] = ch.args[0]

		err := cmd.Run()
		if err != nil {

			return
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s: command not found\n\r", ch.command)
	}
}
