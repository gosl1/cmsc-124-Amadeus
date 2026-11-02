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

func (p *Parser) parseExpression() Expr {
	return p.parseAssignment()
}

func (p *Parser) parseAssignment() Expr {
	expr := p.parseOr()
	if p.match(EQUAL) {
		equals := p.previous()
		value := p.parseAssignment()
		variable, ok := expr.(Variable)
		if !ok {
			p.fail(equals, "Invalid assignment target.")
		}
		return Assign{Name: variable.Name, Value: value}
	}
	return expr
}

func (p *Parser) parseOr() Expr {
	expr := p.parseAnd()
	for p.match(OR) {
		expr = Binary{Left: expr, Operator: p.previous(), Right: p.parseAnd()}
	}
	return expr
}

func (p *Parser) parseAnd() Expr {
	expr := p.parseEquality()
	for p.match(AND) {
		expr = Binary{Left: expr, Operator: p.previous(), Right: p.parseEquality()}
	}
	return expr
}

func (p *Parser) parseEquality() Expr {
	expr := p.parseComparison()
	for p.match(EQUAL_EQUAL, BANG_EQUAL) {
		expr = Binary{Left: expr, Operator: p.previous(), Right: p.parseComparison()}
	}
	return expr
}

func (p *Parser) parseComparison() Expr {
	expr := p.parseTerm()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		expr = Binary{Left: expr, Operator: p.previous(), Right: p.parseTerm()}
	}
	return expr
}

func (p *Parser) parseTerm() Expr {
	expr := p.parseFactor()
	for p.match(MINUS, PLUS) {
		expr = Binary{Left: expr, Operator: p.previous(), Right: p.parseFactor()}
	}
	return expr
}

func (p *Parser) parseFactor() Expr {
	expr := p.parseUnary()
	for p.match(SLASH, STAR, PERCENT) {
		expr = Binary{Left: expr, Operator: p.previous(), Right: p.parseUnary()}
	}
	return expr
}

func (p *Parser) parseUnary() Expr {
	if p.match(BANG, MINUS) {
		return Unary{Operator: p.previous(), Right: p.parseUnary()}
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() Expr {
	if p.match(FALSE) {
		return Literal{Value: false}
	}
	if p.match(TRUE) {
		return Literal{Value: true}
	}
	if p.match(NULL) {
		return Literal{Value: nil}
	}
	if p.match(NUMBER, STRING) {
		return Literal{Value: p.previous().Literal}
	}
	if p.match(IDENTIFIER) {
		return Variable{Name: p.previous()}
	}
	if p.match(LEFT_PAREN) {
		expr := p.parseExpression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{Expression: expr}
	}

	p.fail(p.peek(), "Expect expression.")
	return nil
}
