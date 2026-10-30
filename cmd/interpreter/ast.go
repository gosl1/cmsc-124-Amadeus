package main

// Expr is the interface implemented by every expression node.
type Expr interface {
	isExpr()
}

// Literal represents a literal value, such as a number, string, boolean, or null.
type Literal struct {
	Value any
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
