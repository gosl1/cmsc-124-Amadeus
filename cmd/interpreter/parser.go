package main

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) peek() Token {
	if p.current >= len(p.tokens) {
		return Token{Type: EOF, Line: 1}
	}
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	if p.current == 0 {
		return Token{Type: EOF, Line: p.peek().Line}
	}
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(types ...TokenType) bool {
	for _, typ := range types {
		if p.peek().Type == typ {
			return true
		}
	}

	return false
}

func (p *Parser) match(types ...TokenType) bool {
	if p.check(types...) {
		p.advance()
		return true
	}

	return false
}

type parseError struct {
	token   Token
	message string
}

func (e parseError) Error() string {
	return fmt.Sprintf("[Parser Error] line %d: %s", e.token.Line, e.message)
}

func (p *Parser) fail(token Token, message string) {
	panic(parseError{
		token:   token,
		message: message,
	})
}

func (p *Parser) consume(tokenType TokenType, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	}

	p.fail(p.peek(), message)
	return Token{}
}
