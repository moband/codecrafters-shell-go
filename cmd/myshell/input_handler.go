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

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		j := 0
		for j < len(prefix) && j < len(strs[i]) && prefix[j] == strs[i][j] {
			j++
		}
		prefix = prefix[:j]
		if prefix == "" {
			break
		}
	}

	return prefix
}

func (ih *InputHandler) findCompletion(partial string) ([]string, bool, string) {
	if partial == "" {
		return nil, false, ""
	}

	if completion := ih.findBuiltinCompletion(partial); completion != "" {
		return []string{completion}, true, completion
	}

	completions := ih.findExecutableCompletion(partial)
	if len(completions) == 0 {
		return nil, false, ""
	}

	if len(completions) == 1 {
		return completions, true, completions[0]
	}

	sort.Strings(completions)
	commonPrefix := longestCommonPrefix(completions)
	return completions, commonPrefix == partial, commonPrefix
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
			currentCmd := command.String()
			completions, isExactMatch, commonPrefix := ih.findCompletion(currentCmd)

			if len(completions) > 0 {
				if isExactMatch && len(completions) == 1 {
					fmt.Printf("\r$ %s ", commonPrefix)
					command.Reset()
					command.WriteString(commonPrefix)
					command.WriteByte(' ')
				} else if !isExactMatch {
					fmt.Printf("\r$ %s", commonPrefix)
					command.Reset()
					command.WriteString(commonPrefix)
				} else if lastWasTab {
					fmt.Printf("\n\r%s\n\r", strings.Join(completions, "  "))
					fmt.Printf("$ %s", command.String())
					lastWasTab = false
				} else {
					fmt.Print("\a")
					lastWasTab = true
				}
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
