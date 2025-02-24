package main

// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"strings"

// 	"golang.org/x/term"
// )

// type CommandHandler struct {
// 	command      string
// 	args         []string
// 	outputFile   string
// 	stderrFile   string
// 	appendMode   bool
// 	stderrAppend bool
// }

// type RedirectionConfig struct {
// 	outputFile   string
// 	stderrFile   string
// 	appendMode   bool
// 	stderrAppend bool
// }

// func NewCommandHandler(command string) *CommandHandler {
// 	parser := NewParser(command)
// 	args := parser.Parse()
// 	ch := &CommandHandler{
// 		command: command,
// 		args:    make([]string, 0, len(args)),
// 	}

// 	ch.parseRedirections(args)
// 	return ch
// }

// func (ch *CommandHandler) parseRedirections(args []string) {
// 	config := &RedirectionConfig{}

// 	for i := 0; i < len(args); i++ {
// 		switch {
// 		case (args[i] == ">>" || args[i] == "1>>") && i+1 < len(args):
// 			config.outputFile = args[i+1]
// 			config.appendMode = true
// 			i++

// 		case (args[i] == ">" || args[i] == "1>") && i+1 < len(args):
// 			config.outputFile = args[i+1]
// 			i++

// 		case args[i] == "2>>" && i+1 < len(args):
// 			config.stderrFile = args[i+1]
// 			config.stderrAppend = true
// 			i++

// 		case args[i] == "2>" && i+1 < len(args):
// 			config.stderrFile = args[i+1]
// 			i++

// 		default:
// 			ch.args = append(ch.args, args[i])
// 		}
// 	}

// 	ch.outputFile = config.outputFile
// 	ch.stderrFile = config.stderrFile
// 	ch.appendMode = config.appendMode
// 	ch.stderrAppend = config.stderrAppend
// }

// func (ch *CommandHandler) setupRedirection() error {
// 	fm := &FileManager{}

// 	if ch.outputFile != "" {
// 		file, err := fm.CreateFile(ch.outputFile, ch.appendMode)
// 		if err != nil {
// 			return fmt.Errorf("failed to create output file: %w", err)
// 		}
// 		os.Stdout = file
// 		defer func() {
// 			file.Close()
// 			os.Stdout = os.NewFile(1, "/dev/stdout")
// 		}()
// 	}

// 	if ch.stderrFile != "" {
// 		file, err := fm.CreateFile(ch.stderrFile, ch.stderrAppend)
// 		if err != nil {
// 			return fmt.Errorf("failed to create error file: %w", err)
// 		}
// 		os.Stderr = file
// 		defer func() {
// 			file.Close()
// 			os.Stderr = os.NewFile(2, "/dev/stderr")
// 		}()
// 	}

// 	return nil
// }

// func (ch *CommandHandler) Execute(shell *Shell) {
// 	if len(ch.args) == 0 {
// 		return
// 	}

// 	ch.setupRedirection()

// 	if handler, ok := shell.builtins[ch.args[0]]; ok {
// 		handler(ch)
// 	} else {
// 		ch.handleExternal()
// 	}
// }

// func (ch *CommandHandler) handleExit() {
// 	term.Restore(int(os.Stdin.Fd()), oldState)
// 	os.Exit(0)
// }

// func (ch *CommandHandler) handleEcho() {
// 	if len(ch.args) > 1 {
// 		fmt.Fprintf(os.Stdout, "%s\n\r", strings.Join(ch.args[1:], " "))
// 	}
// }

// func (ch *CommandHandler) handlePwd() {
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%v\n\r", err)
// 	} else {
// 		fmt.Fprintf(os.Stdout, "%s\n\r", wd)
// 	}
// }

// func (ch *CommandHandler) handleCd() {
// 	if len(ch.args) <= 1 {
// 		return
// 	}

// 	dir := ch.args[1]
// 	if dir == "~" {
// 		homeDir, err := os.UserHomeDir()
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "%v\n\r", err)
// 			return
// 		}
// 		dir = homeDir
// 	}

// 	if err := os.Chdir(dir); err != nil {
// 		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n\r", dir)
// 	}
// }

// func (ch *CommandHandler) handleType() {
// 	if len(ch.args) <= 1 {
// 		return
// 	}

// 	cmd := ch.args[1]
// 	switch cmd {
// 	case "echo", "exit", "type", "pwd", "cd":
// 		fmt.Fprintf(os.Stdout, "%s is a shell builtin\n\r", cmd)
// 	default:
// 		if execPath, found := findExecutable(cmd); found {
// 			fmt.Fprintf(os.Stdout, "%s is %s\n\r", cmd, execPath)
// 		} else {
// 			fmt.Fprintf(os.Stderr, "%s: not found\n\r", cmd)
// 		}
// 	}
// }

// func (ch *CommandHandler) handleExternal() {
// 	if execPath, found := findExecutable(ch.args[0]); found {
// 		cmd := exec.Command(execPath, ch.args[1:]...)

// 		stdout := &lineWriter{w: os.Stdout}
// 		stderr := &lineWriter{w: os.Stderr}

// 		cmd.Stdout = stdout
// 		cmd.Stderr = stderr
// 		cmd.Args[0] = ch.args[0]
// 		cmd.Run()
// 	} else {
// 		fmt.Fprintf(os.Stderr, "%s: command not found\n\r", ch.command)
// 	}
// }
