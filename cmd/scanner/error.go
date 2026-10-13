package main

import (
	"fmt"
	"os"
)


// reportError prints one diagnostic to stderr and marks the scan as
// failed. Change the message format here — this is the one place it
// lives. Nothing about a rejected file goes to stdout, ever.
func (s *Scanner) reportError(line int, message string) {
	fmt.Fprintf(os.Stderr, "[Attractor Field Violation] line %d: %s\n", line, message)
	s.hadError = true
}