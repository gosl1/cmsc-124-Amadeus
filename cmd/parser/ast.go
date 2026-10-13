package main

import (
	"fmt"
	"strconv"
)

// Expr is the common interface for every node in the parse tree.
type Expr interface{}

// Literal stores a primitive value.
type Literal struct {
	Value any
}

// Variable stores an identifier reference.
type Variable struct {
	Name Token
}

// Unary stores a unary operator and its operand.
type Unary struct {
	Operator Token
	Right    Expr
}

// Binary stores a left operand, an operator, and a right operand.
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

// Grouping stores one inner expression wrapped in parentheses.
type Grouping struct {
	Expression Expr
}

// AstPrinter renders the tree in a parenthesized prefix form.
type AstPrinter struct{}

func (p AstPrinter) Print(expr Expr) string {
	switch node := expr.(type) {
	case nil:
		return "null"
	case Literal:
		return printLiteral(node.Value)
	case Variable:
		return node.Name.Lexeme
	case Unary:
		return fmt.Sprintf("(%s %s)", node.Operator.Lexeme, p.Print(node.Right))
	case Binary:
		return fmt.Sprintf("(%s %s %s)", node.Operator.Lexeme, p.Print(node.Left), p.Print(node.Right))
	case Grouping:
		return p.Print(node.Expression)
	default:
		return "<?>"
	}
}

func printLiteral(value any) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case bool:
		if v {
			return "true"
		}
		return "false"
	case string:
		return strconv.Quote(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	default:
		return fmt.Sprintf("%v", v)
	}
}
