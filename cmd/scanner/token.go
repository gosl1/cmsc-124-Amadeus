package main

import "fmt"

// TokenType categorizes a lexeme so downstream phases can switch on a
// category instead of comparing raw strings.
type TokenType string

const (
	// grouping
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"

	// arithmetic
	PLUS    TokenType = "PLUS"
	MINUS   TokenType = "MINUS"
	STAR    TokenType = "STAR"
	SLASH   TokenType = "SLASH"
	PERCENT TokenType = "PERCENT"

	// assignment
	EQUAL TokenType = "EQUAL"

	// comparison
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"

	// logical
	BANG TokenType = "BANG"
	AND  TokenType = "AND"
	OR   TokenType = "OR"

	// indexing / access
	LEFT_BRACKET  TokenType = "LEFT_BRACKET"
	RIGHT_BRACKET TokenType = "RIGHT_BRACKET"
	DOT           TokenType = "DOT"
	COMMA         TokenType = "COMMA"

	// literals
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"

	// keywords (Amadeus vocabulary)
	RES            TokenType = "RES"            // res   — variable declaration
	DMAIL          TokenType = "DMAIL"          // dmail — print/output
	DIV            TokenType = "DIV"            // div   — if
	CON            TokenType = "CON"            // con   — else
	TIMELEAP       TokenType = "TIMELEAP"       // timeleap — while
	OPERATION      TokenType = "OPERATION"      // operation — function declaration
	ELPSY          TokenType = "ELPSY"          // elpsy — return
	READINGSTEINER TokenType = "READINGSTEINER" // readingsteiner — execution-history value
	TRUE           TokenType = "TRUE"
	FALSE          TokenType = "FALSE"
	NULL           TokenType = "NULL"

	EOF TokenType = "EOF"
)

// Token bundles a lexeme with the information downstream phases need.
type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func (t Token) String() string {
	lit := "null"
	if t.Literal != nil {
		lit = fmt.Sprintf("%v", t.Literal)
	}
	return fmt.Sprintf("Token(type=%s, lexeme=%s, literal=%s, line=%d)",
		t.Type, t.Lexeme, lit, t.Line)
}
