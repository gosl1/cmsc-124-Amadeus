package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) >= 2 && args[0] == "--tokenize" {
		runFileTokenize(args[1])
		return
	}

	if len(args) == 0 {
		runREPL()
		return
	}

	// Anything else (missing file, unknown flag) is a usage error.
	fmt.Fprintln(os.Stderr, "usage: run [--tokenize <file>]")
	os.Exit(64)

}

// runREPL reads one line at a time, scans it, and prints its tokens.
func runREPL() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if line != "" {
			hadError = false
			tokens, failed := NewScanner(line).ScanTokens()
			// Same rule as file mode: a rejected line prints nothing to stdout.
			if !failed {
				for _, token := range tokens {
					fmt.Println(token)
				}
			}
		}
		if err != nil { // io.EOF (Ctrl-D / Ctrl-Z) ends the session cleanly
			fmt.Println()
			os.Exit(0)
		}
	}
}

func runFileTokenize(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ama: cannot read '%s': %v\n", path, err)
		os.Exit(65)
	}
	scanner := NewScanner(string(source))

	tokens, failed := scanner.ScanTokens()

	if failed {
		os.Exit(65)
	}
	for _, token := range tokens {
		fmt.Println(token)
	}
	os.Exit(0)
}
