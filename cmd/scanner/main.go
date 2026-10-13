package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 2 && args[0] == "--tokenize" {
		runFileTokenize(args[1])
		return
	}

	if len(args) == 2 && args[0] == "--parse" {
		runFileParse(args[1])
		return
	}

	fmt.Fprintln(os.Stderr, "usage: run [--tokenize|--parse] <file>")
	os.Exit(65)

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

func runFileParse(path string) {
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

	expr, err := NewParser(tokens).Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(65)
	}

	fmt.Println(AstPrinter{}.Print(expr))
	os.Exit(0)
}
