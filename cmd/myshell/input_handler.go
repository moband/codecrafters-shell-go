package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	tab       = 9
	backspace = 127
	ctrlC     = 3
	ctrlD     = 4
)

type InputHandler struct {
	shell *Shell
}

func NewInputHandler(shell *Shell) *InputHandler {
	return &InputHandler{shell: shell}
}

func (ih *InputHandler) findCompletion(partial string) (string, bool) {
	if partial == "" {
		return "", false
	}

	if completion := ih.findBuiltinCompletion(partial); completion != "" {
		return completion, true
	}

	if completion := ih.findExecutableCompletion(partial); completion != "" {
		return completion, true
	}

	return partial, false
}

func (ih *InputHandler) findBuiltinCompletion(partial string) string {
	for cmd := range ih.shell.builtins {
		if strings.HasPrefix(cmd, partial) {
			return cmd
		}
	}
	return ""
}

func (ih *InputHandler) findExecutableCompletion(partial string) string {
	paths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	for _, path := range paths {
		files, _ := os.ReadDir(path)
		for _, file := range files {
			if strings.HasPrefix(file.Name(), partial) {
				if info, err := file.Info(); err == nil && !info.IsDir() && info.Mode()&0111 != 0 {
					return file.Name()
				}
			}
		}
	}
	return ""
}

func (ih *InputHandler) readInput() string {
	var command strings.Builder
	buffer := make([]byte, 1)

	for {
		if _, err := os.Stdin.Read(buffer); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n\r", err)
			term.Restore(int(os.Stdin.Fd()), oldState)
			os.Exit(1)
		}

		switch buffer[0] {
		case '\r', '\n':
			fmt.Fprintf(os.Stdout, "\n\r")
			if cmd := command.String(); cmd != "" {
				return strings.TrimSpace(cmd)
			}

		case tab:
			if completed, ok := ih.findCompletion(command.String()); ok {
				fmt.Printf("\r$ %s ", completed)
				command.Reset()
				command.WriteString(completed)
				command.WriteByte(' ')
			} else {
				fmt.Print("\a")
			}
			continue

		case ctrlD, ctrlC:

			fmt.Fprintf(os.Stdout, "\n\r")
			term.Restore(int(os.Stdin.Fd()), oldState)
			os.Exit(0)

		default:
			fmt.Printf("%c", buffer[0])
			command.WriteByte(buffer[0])
			continue
		}

		break
	}

	return strings.TrimSpace(command.String())
}
