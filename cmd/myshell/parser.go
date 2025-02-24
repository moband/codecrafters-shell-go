package main

import "strings"

type Parser struct {
	command string
}

func NewParser(command string) *Parser {
	return &Parser{command: command}
}

func (p *Parser) Parse() []string {
	var args []string
	var current strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false
	i := 0

	for i < len(p.command) {

		if p.command[i] == '\\' && !inSingleQuotes && !inDoubleQuotes && i+1 < len(p.command) {
			current.WriteByte(p.command[i+1])
			i += 2
			continue
		}

		if p.command[i] == '\'' && !inDoubleQuotes {
			inSingleQuotes = !inSingleQuotes
			i++
			continue
		}

		if p.command[i] == '"' && !inSingleQuotes {
			inDoubleQuotes = !inDoubleQuotes
			i++
			continue
		}

		if p.command[i] == '\\' && inDoubleQuotes && i+1 < len(p.command) {
			next := p.command[i+1]
			if next == '\\' || next == '$' || next == '"' || next == '\n' {
				current.WriteByte(next)
				i += 2
				continue
			}
		}

		if p.command[i] == ' ' && !inSingleQuotes && !inDoubleQuotes {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		} else {
			current.WriteByte(p.command[i])
		}
		i++
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}
