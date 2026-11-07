package main

import (
	"io"
	"os"
	"testing"
)

func TestParseErrorReporting(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "unexpected trailing token",
			src:  "true false",
			want: "[line 1] Error at 'false': Unexpected token.\n",
		},
		{
			name: "missing operand at end of input",
			src:  "1 +",
			want: "[line 1] Error at end: Expect expression.\n",
		},
		{
			name: "end of input line",
			src:  "1 +\n",
			want: "[line 2] Error at end: Expect expression.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hadError = false
			tokens, failed := NewScanner(tt.src).ScanTokens()
			if failed {
				t.Fatal("scanner rejected parser test input")
			}

			_, err := NewParser(tokens).Parse()
			if err == nil {
				t.Fatal("expected parse error")
			}
			parseErr, ok := err.(parseError)
			if !ok {
				t.Fatalf("expected parseError, got %T", err)
			}

			oldStderr := os.Stderr
			reader, writer, pipeErr := os.Pipe()
			if pipeErr != nil {
				t.Fatal(pipeErr)
			}
			os.Stderr = writer
			reportParseError(parseErr)
			_ = writer.Close()
			os.Stderr = oldStderr

			got, readErr := io.ReadAll(reader)
			_ = reader.Close()
			if readErr != nil {
				t.Fatal(readErr)
			}
			if string(got) != tt.want {
				t.Errorf("diagnostic mismatch:\n got: %q\nwant: %q", got, tt.want)
			}
		})
	}
}

func TestParseExpressionsByLine(t *testing.T) {
	hadError = false
	tokens, failed := NewScanner("1 + 2\n\n// ignored\n3 * 4").ScanTokens()
	if failed {
		t.Fatal("scanner rejected parser test input")
	}

	expressions, parseErrors := parseExpressions(tokens)
	if len(parseErrors) != 0 {
		t.Fatalf("unexpected parse errors: %v", parseErrors)
	}
	if got, want := len(expressions), 2; got != want {
		t.Fatalf("parsed %d expressions, want %d", got, want)
	}
	if got, want := expressions[0].String(), "(+ 1 2)"; got != want {
		t.Errorf("first expression = %q, want %q", got, want)
	}
	if got, want := expressions[1].String(), "(* 3 4)"; got != want {
		t.Errorf("second expression = %q, want %q", got, want)
	}
}

func TestParseExpressionsContinuesAfterErrors(t *testing.T) {
	hadError = false
	tokens, failed := NewScanner("1 + )\n2 + * 3\n4").ScanTokens()
	if failed {
		t.Fatal("scanner rejected parser recovery test input")
	}

	expressions, parseErrors := parseExpressions(tokens)
	if got, want := len(parseErrors), 2; got != want {
		t.Fatalf("got %d parse errors, want %d", got, want)
	}
	if got, want := len(expressions), 1; got != want {
		t.Fatalf("got %d valid expressions, want %d", got, want)
	}
	if got, want := expressions[0].String(), "4"; got != want {
		t.Errorf("expression after errors = %q, want %q", got, want)
	}
	if got, want := parseErrors[0].token.Line, 1; got != want {
		t.Errorf("first parse error line = %d, want %d", got, want)
	}
	if got, want := parseErrors[1].token.Line, 2; got != want {
		t.Errorf("second parse error line = %d, want %d", got, want)
	}
}

func TestParseExpressionsReportsErrorAfterValidLine(t *testing.T) {
	hadError = false
	tokens, failed := NewScanner("1 + 2\n3 + )\n4").ScanTokens()
	if failed {
		t.Fatal("scanner rejected parser recovery test input")
	}

	expressions, parseErrors := parseExpressions(tokens)
	if got, want := len(parseErrors), 1; got != want {
		t.Fatalf("got %d parse errors, want %d", got, want)
	}
	if got, want := len(expressions), 2; got != want {
		t.Fatalf("got %d valid expressions, want %d", got, want)
	}
	if got, want := parseErrors[0].token.Line, 2; got != want {
		t.Errorf("parse error line = %d, want %d", got, want)
	}
}
