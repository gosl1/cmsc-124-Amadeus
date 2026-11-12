package main

import (
	"fmt"
	"strconv"
)

// Expr is the interface implemented by every expression node.
type Expr interface {
	isExpr()
	String() string
}

// Literal represents a literal value, such as a number, string, boolean, or null.
type Literal struct {
	Value any
	Token Token
}

// Variable represents a variable reference.
type Variable struct {
	Name Token
}

// Assign represents assigning a value to a variable.
type Assign struct {
	Name  Token
	Value Expr
}

// Unary represents an expression with one operand.
// Examples: !true, -5
type Unary struct {
	Operator Token
	Right    Expr
}

// Binary represents an expression with two operands.
// Examples: 1 + 2, 5 * 3
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

// Grouping represents an expression enclosed in parentheses.
type Grouping struct {
	Expression Expr
}

// Marker methods: these make each struct satisfy the Expr interface.
func (Literal) isExpr()  {}
func (Variable) isExpr() {}
func (Assign) isExpr()   {}
func (Unary) isExpr()    {}
func (Binary) isExpr()   {}
func (Grouping) isExpr() {}

func (e Literal) String() string {
	switch value := e.Value.(type) {
	case nil:
		return "null"
	case string:
		return strconv.Quote(value)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	default:
		return fmt.Sprint(value)
	}
}

func (e Variable) String() string {
	return e.Name.Lexeme
}

func (e Assign) String() string {
	return fmt.Sprintf("(= %s %s)", e.Name.Lexeme, e.Value)
}

func (e Unary) String() string {
	return fmt.Sprintf("(%s %s)", e.Operator.Lexeme, e.Right)
}

func (e Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Operator.Lexeme, e.Left, e.Right)
}

func (e Grouping) String() string {
	return fmt.Sprintf("(group %s)", e.Expression)
}
