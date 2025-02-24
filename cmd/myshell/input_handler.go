package main

import (
	"fmt"
	"os"
	"sort"
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

func (ih *InputHandler) findCompletion(partial string) ([]string, bool) {
	if partial == "" {
		return nil, false
	}

	if completion := ih.findBuiltinCompletion(partial); completion != "" {
		return []string{completion}, true
	}

	completions := ih.findExecutableCompletion(partial)
	if len(completions) == 1 {
		return completions, true
	}
	sort.Strings(completions)
	return completions, false
}

func (ih *InputHandler) findBuiltinCompletion(partial string) string {
	for cmd := range ih.shell.builtins {
		if strings.HasPrefix(cmd, partial) {
			return cmd
		}
	}
	return ""
}

func (ih *InputHandler) findExecutableCompletion(partial string) []string {
	var matches []string
	seen := make(map[string]bool)

	paths := strings.Split(os.Getenv("PATH"), ":")
	for _, dir := range paths {
		files, _ := os.ReadDir(dir)
		for _, file := range files {
			name := file.Name()
			if strings.HasPrefix(name, partial) {
				if info, err := file.Info(); err == nil && !info.IsDir() && info.Mode()&0111 != 0 {
					if !seen[name] {
						matches = append(matches, name)
						seen[name] = true
					}
				}
			}
		}
	}
	return matches
}

func (ih *InputHandler) readInput() string {
	var command strings.Builder
	buffer := make([]byte, 1)
	lastWasTab := false

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
			if completions, ok := ih.findCompletion(command.String()); ok {
				fmt.Printf("\r$ %s ", completions[0])
				command.Reset()
				command.WriteString(completions[0])
				command.WriteByte(' ')

			} else if len(completions) > 1 && lastWasTab {
				fmt.Printf("\n\r%s\n\r", strings.Join(completions, "  "))
				fmt.Printf("$ %s", command.String())
				lastWasTab = false
			} else {
				fmt.Print("\a")
				lastWasTab = true
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
