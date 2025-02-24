package main

// import "strings"

// type Parser struct {
// 	command string
// }

// func NewParser(command string) *Parser {
// 	return &Parser{command: command}
// }

// func (p *Parser) Parse() []string {
// 	var args []string
// 	var current strings.Builder
// 	inSingleQuotes := false
// 	inDoubleQuotes := false

// 	for i := 0; i < len(p.command); i++ {
// 		switch {
// 		case p.command[i] == '\\' && !inSingleQuotes && i+1 < len(p.command):
// 			if !inDoubleQuotes || isEscapableInDoubleQuotes(p.command[i+1]) {
// 				current.WriteByte(p.command[i+1])
// 				i++
// 			} else {
// 				current.WriteByte(p.command[i])
// 			}

// 		case p.command[i] == '\'' && !inDoubleQuotes:
// 			inSingleQuotes = !inSingleQuotes

// 		case p.command[i] == '"' && !inSingleQuotes:
// 			inDoubleQuotes = !inDoubleQuotes

// 		case p.command[i] == ' ' && !inSingleQuotes && !inDoubleQuotes:
// 			if current.Len() > 0 {
// 				args = append(args, current.String())
// 				current.Reset()
// 			}

// 		default:
// 			current.WriteByte(p.command[i])
// 		}
// 	}

// 	if current.Len() > 0 {
// 		args = append(args, current.String())
// 	}

// 	return args
// }

// func isEscapableInDoubleQuotes(ch byte) bool {
// 	return ch == '\\' || ch == '$' || ch == '"' || ch == '\n'
// }
