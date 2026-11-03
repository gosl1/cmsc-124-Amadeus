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
