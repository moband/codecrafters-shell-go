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

func (ih *InputHandler) findCompletion(partial string) string {
	if partial == "" {
		return ""
	}

	for cmd := range shell.builtins {
		if strings.HasPrefix(cmd, partial) {
			return cmd
		}
	}
	return partial
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
			if completed := ih.findCompletion(command.String()); completed != command.String() {
				fmt.Printf("\r$ %s ", completed)
				command.Reset()
				command.WriteString(completed)
				command.WriteByte(' ')
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
