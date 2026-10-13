package main

import "fmt"

// Parser turns a token stream into an abstract syntax tree.
type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (Expr, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if !p.isAtEnd() {
		return nil, p.errorToken(p.peek(), "Unexpected token after expression.")
	}
	return expr, nil
}

func (p *Parser) parseExpression() (Expr, error) {
	return p.parseAssignment()
}

func (p *Parser) parseAssignment() (Expr, error) {
	left, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	if p.match(EQUAL) {
		op := p.previous()
		right, err := p.parseAssignment()
		if err != nil {
			return nil, err
		}
		return Binary{Left: left, Operator: op, Right: right}, nil
	}
	return left, nil
}

func (p *Parser) parseOr() (Expr, error) {
	expr, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.match(OR) {
		op := p.previous()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) parseAnd() (Expr, error) {
	expr, err := p.parseEquality()
	if err != nil {
		return nil, err
	}
	for p.match(AND) {
		op := p.previous()
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) parseEquality() (Expr, error) {
	expr, err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	for p.match(EQUAL_EQUAL, BANG_EQUAL) {
		op := p.previous()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) parseComparison() (Expr, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for p.match(LESS, LESS_EQUAL, GREATER, GREATER_EQUAL) {
		op := p.previous()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) parseTerm() (Expr, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for p.match(PLUS, MINUS) {
		op := p.previous()
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) parseFactor() (Expr, error) {
	expr, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for p.match(STAR, SLASH, PERCENT) {
		op := p.previous()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) parseUnary() (Expr, error) {
	if p.match(BANG, MINUS) {
		op := p.previous()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return Unary{Operator: op, Right: right}, nil
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (Expr, error) {
	if p.match(FALSE) {
		return Literal{Value: false}, nil
	}
	if p.match(TRUE) {
		return Literal{Value: true}, nil
	}
	if p.match(NULL) {
		return Literal{Value: nil}, nil
	}
	if p.match(NUMBER) {
		return Literal{Value: p.previous().Literal}, nil
	}
	if p.match(STRING) {
		return Literal{Value: p.previous().Literal}, nil
	}
	if p.match(IDENTIFIER) {
		return Variable{Name: p.previous()}, nil
	}
	if p.match(LEFT_PAREN) {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if !p.match(RIGHT_PAREN) {
			return nil, p.errorToken(p.peek(), "Expect ')' after expression.")
		}
		return Grouping{Expression: expr}, nil
	}
	return nil, p.errorToken(p.peek(), "Expect expression.")
}

func (p *Parser) match(types ...TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) errorToken(token Token, message string) error {
	return fmt.Errorf("line %d: %s", token.Line, message)
}
