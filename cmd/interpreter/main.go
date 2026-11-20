package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	if len(args) == 2 && args[0] == "--eval" {
		runFileEval(args[1])
		return
	}

	if len(args) == 2 && args[0] == "--parse" {
		runFileParse(args[1])
		return
	}

	if len(args) == 2 && args[0] == "--tokenize" {
		runFileTokenize(args[1])
		return
	}

	if len(args) == 0 {
		runREPL()
		return
	}

	// Anything else (missing file, unknown flag) is a usage error.
	fmt.Fprintln(os.Stderr, "usage: run [--tokenize <file> | --parse <file> | --eval <file>]")
	os.Exit(64)

}

// runREPL reads and evaluates one expression at a time.
func runREPL() {
	reader := bufio.NewReader(os.Stdin)
	evaluator := NewEvaluator()
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if line != "" {
			tokens, failed := NewScanner(line).ScanTokens()
			if !failed && len(tokens) > 1 {
				expr, parseErr := NewParser(tokens).Parse()
				if parseErr != nil {
					if syntaxErr, ok := parseErr.(parseError); ok {
						reportParseError(syntaxErr)
					} else {
						fmt.Fprintln(os.Stderr, parseErr)
					}
				} else {
					value, evalErr := evaluator.Evaluate(expr)
					if evalErr != nil {
						reportRuntimeError(evalErr)
					} else {
						fmt.Println(formatValue(value))
					}
				}
			}
		}
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "ama: input error: %v\n", err)
			}
			fmt.Println()
			return
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

func runFileEval(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ama: cannot read '%s': %v\n", path, err)
		os.Exit(65)
	}

	tokens, failed := NewScanner(string(source)).ScanTokens()
	if failed {
		os.Exit(65)
	}

	expressions, parseErrors := parseExpressions(tokens)
	if len(parseErrors) > 0 {
		for _, parseErr := range parseErrors {
			reportParseError(parseErr)
		}
		os.Exit(65)
	}

	evaluator := NewEvaluator()
	for _, expr := range expressions {
		value, err := evaluator.Evaluate(expr)
		if err != nil {
			reportRuntimeError(err)
			os.Exit(70)
		}
		fmt.Println(formatValue(value))
	}
	os.Exit(0)
}

func reportRuntimeError(err error) {
	var runtimeErr RuntimeError
	if errors.As(err, &runtimeErr) {
		fmt.Fprintf(os.Stderr, "[line %d] Runtime error: %s\n",
			runtimeErr.Operator.Line, runtimeErr.Message)
		return
	}
	fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
}

func formatValue(value any) string {
	switch value := value.(type) {
	case nil:
		return "null"
	case string:
		return value
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	default:
		return fmt.Sprint(value)
	}
}

func runFileParse(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ama: cannot read '%s': %v\n", path, err)
		os.Exit(65)
	}

	tokens, failed := NewScanner(string(source)).ScanTokens()
	if failed {
		os.Exit(65)
	}

	expressions, parseErrors := parseExpressions(tokens)
	if len(parseErrors) > 0 {
		for _, parseErr := range parseErrors {
			reportParseError(parseErr)
		}
		os.Exit(65)
	}
	for _, expr := range expressions {
		fmt.Println(expr)
	}
	os.Exit(0)
}
