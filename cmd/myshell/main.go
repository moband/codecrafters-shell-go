package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var oldState *term.State
var shell *Shell

func main() {
	shell = NewShell()
	ih := NewInputHandler(shell)
	var err error
	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set raw mode: %v\n\r", err)
		os.Exit(1)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	const prompt = "$ "

	for {
		fmt.Fprint(os.Stdout, prompt)
		command := ih.readInput()
		NewCommandHandler(command).Execute(shell)
	}
}
