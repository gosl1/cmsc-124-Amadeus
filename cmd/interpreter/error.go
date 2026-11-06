package main

import (
	"fmt"
	"os"
)

var hadError = false

// reportError prints one diagnostic to stderr and marks the scan as
// failed. Change the message format here — this is the one place it
// lives. Nothing about a rejected file goes to stdout, ever.
func reportError(line int, message string) {
	fmt.Fprintf(os.Stderr, "[Attractor Field Violation] line %d: %s\n", line, message)
	hadError = true
}

func formatParseError(token Token, message string) string {
	if token.Type == EOF {
		return fmt.Sprintf("[line %d] Error at end: %s", token.Line, message)
	}
	return fmt.Sprintf("[line %d] Error at '%s': %s", token.Line, token.Lexeme, message)
}

func reportParseError(parseErr parseError) {
	fmt.Fprintln(os.Stderr, formatParseError(parseErr.token, parseErr.message))
}
