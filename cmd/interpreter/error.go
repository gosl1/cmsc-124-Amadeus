package main

import (
	"fmt"
	"os"
)

func formatParseError(token Token, message string) string {
	if token.Type == EOF {
		return fmt.Sprintf("[line %d] Error at end: %s", token.Line, message)
	}
	return fmt.Sprintf("[line %d] Error at '%s': %s", token.Line, token.Lexeme, message)
}

func reportParseError(parseErr parseError) {
	fmt.Fprintln(os.Stderr, formatParseError(parseErr.token, parseErr.message))
}
