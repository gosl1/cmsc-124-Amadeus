package main

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int

	hadError bool
}

// keywords maps each themed lexeme to its token type.
var keywords = map[string]TokenType{
	"res":            RES,
	"dmail":          DMAIL,
	"div":            DIV,
	"con":            CON,
	"timeleap":       TIMELEAP,
	"operation":      OPERATION,
	"elpsy":          ELPSY,
	"readingsteiner": READINGSTEINER,
	"true":           TRUE,
	"false":          FALSE,
	"null":           NULL,
	"and":            AND,
	"or":             OR,
}

func NewScanner(source string) *Scanner {
	//Gives a pointer to Scanner and creates the variable implicitly
	return &Scanner{
		source: source,
		tokens: []Token{},
		line:   1,
	}
}

func (s *Scanner) ScanTokens() ([]Token, bool) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, Token{Type: EOF, Lexeme: "", Literal: nil, Line: s.line})
	return s.tokens, s.hadError
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case '[':
		s.addToken(LEFT_BRACKET, nil)
	case ']':
		s.addToken(RIGHT_BRACKET, nil)
	case '.':
		s.addToken(DOT, nil)
	case '+':
		s.addToken(PLUS, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '*':
		s.addToken(STAR, nil)
	case '%':
		s.addToken(PERCENT, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL, nil)
		} else {
			s.addToken(EQUAL, nil)
		}
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL, nil)
		} else {
			s.addToken(BANG, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL, nil)
		} else {
			s.addToken(LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL, nil)
		} else {
			s.addToken(GREATER, nil)
		}
	case '/':
		if s.match('/') {
			// Line comment: scanned and discarded, runs to end of line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
		}

	// whitespace
	case ' ', '\t', '\r':
    	// discard, emit nothing

	// increment s.line when emcountering new line character
	case '\n':
    	s.line++

	case '"':
		s.stringLiteral()
	
	//case '/':
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.reportError(s.line, fmt.Sprintf("Unexpected character '%c'.", c))
		}
	}
}
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// advance consumes and returns the current character.
func (s *Scanner) advance() byte {
	c := s.source[s.current]
	s.current++
	return c
}

func (s *Scanner) addToken(tokenType TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{Type: tokenType, Lexeme: text, Literal: literal, Line: s.line})
}

// peek reads the current character WITHOUT consuming it.
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

// looks two characters ahead of the consumed character.
// peekNext looks one character further than peek — needed to decide
// whether a "." after digits is a decimal point or its own token.
func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

// match is peek plus a conditional advance: consumes the expected
// character only if it's actually there. This is what lets "=" and
// "==" be told apart with one character of lookahead.
func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

// identifier consumes the WHOLE run of identifier characters first
// (maximal munch), then does ONE table lookup on the finished text.
// A match means keyword; no match means a plain identifier.
func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tokType, ok := keywords[text]
	if !ok {
		tokType = IDENTIFIER
	}
	s.addToken(tokType, nil)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	text := s.source[s.start:s.current]
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		s.reportError(s.line, fmt.Sprintf("Malformed number '%s'.", text))
		return
	}
	s.addToken(NUMBER, value)
}

func (s *Scanner) stringLiteral() {
	var value []byte
 
	for s.peek() != '"' && !s.isAtEnd() && s.peek() != '\n' {
		if s.peek() == '\\' {
			s.advance() // consume the backslash
			if s.isAtEnd() {
				break
			}
			escaped := s.advance()
			switch escaped {
			case 'n':
				value = append(value, '\n')
			case '"':
				value = append(value, '"')
			case '\\':
				value = append(value, '\\')
			default:
				s.reportError(s.line, fmt.Sprintf("Invalid escape sequence '\\%c'.", escaped))
			}
			continue
		}
		value = append(value, s.advance())
	}
 
	if s.isAtEnd() || s.peek() == '\n' {
		s.reportError(s.line, "Unterminated string.")
		return
	}
 
	s.advance() // consume the closing "
	s.addToken(STRING, string(value))
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}