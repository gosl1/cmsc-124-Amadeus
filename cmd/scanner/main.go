package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) >= 2 && args[0] == "--tokenize" {
		runFileTokenize(args[1])
		return
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
