package main

import "math"

// RuntimeError represents an error that occurs while evaluating an expression.
// It stores the operator responsible for the error so we can identify its
// source location when reporting the error.
type RuntimeError struct {
	Operator Token
	Message  string
}

// Error implements Go's error interface, allowing RuntimeError to be
// returned and handled like any other Go error.
func (e RuntimeError) Error() string {
	return e.Message
}

// Evaluator is responsible for traversing the AST and evaluating expressions.
type Evaluator struct {
}

// NewEvaluator creates and returns a new Evaluator instance.
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// Evaluate recursively evaluates an AST expression and returns its resulting
// runtime value. If evaluation fails, it returns an error instead.
func (e *Evaluator) Evaluate(expr Expr) (any, error) {

	// type switch determines which AST node is being evaluated.
	switch expr := expr.(type) {

	case Literal:
		// Literals already contain their runtime value, so no computation is needed.
		return expr.Value, nil

	case Grouping:
		// Evaluate the expression contained within the parentheses.
		return e.Evaluate(expr.Expression)

	case Unary:
		// Evaluate the operand first before applying the unary operator.
		rightVal, err := e.Evaluate(expr.Right)

		// Propagate any error encountered while evaluating the operand.
		if err != nil {
			return nil, err
		}

		// Determine which unary operation to perform.
		switch expr.Operator.Type {

		case MINUS:
			// Unary negation is only valid for numeric operands.
			// Type assertion checks whether the evaluated value is a float64.
			number, ok := rightVal.(float64)

			if !ok {
				// Return a runtime error if the operand is not a number.
				return nil, RuntimeError{
					Operator: expr.Operator,
					Message:  "Operand must be a number.",
				}
			}

			// Negate the number and return the result without an error.
			return -number, nil

		case BANG:
			// Logical NOT returns the opposite of the operand's truthiness.
			return !isTruthy(rightVal), nil

		}
	case Binary:
		// Evaluate the left operand first.
		leftVal, err := e.Evaluate(expr.Left)
		if err != nil {
			return nil, err
		}

		// Handle logical operators before evaluating the right operand.
		switch expr.Operator.Type {
		case AND:
			// If the left operand is falsey, return it without evaluating the right.
			if !isTruthy(leftVal) {
				return leftVal, nil
			}

		case OR:
			// If the left operand is truthy, return it without evaluating the right.
			if isTruthy(leftVal) {
				return leftVal, nil
			}
		}

		// Evaluate the right operand if necessary.
		rightVal, err := e.Evaluate(expr.Right)
		if err != nil {
			return nil, err
		}

		switch expr.Operator.Type {
		case PLUS:
			// Addition supports either two numbers or two strings.
			leftNum, leftIsNum := leftVal.(float64)
			rightNum, rightIsNum := rightVal.(float64)

			if leftIsNum && rightIsNum {
				return leftNum + rightNum, nil
			}

			leftStr, leftIsStr := leftVal.(string)
			rightStr, rightIsStr := rightVal.(string)

			if leftIsStr && rightIsStr {
				return leftStr + rightStr, nil
			}

			return nil, RuntimeError{
				Operator: expr.Operator,
				Message:  "Operands must be two numbers or two strings.",
			}

		case MINUS, STAR, SLASH, PERCENT,
			GREATER, GREATER_EQUAL, LESS, LESS_EQUAL:

			// These operators require numeric operands.
			leftNum, leftIsNum := leftVal.(float64)
			rightNum, rightIsNum := rightVal.(float64)

			if !leftIsNum || !rightIsNum {
				return nil, RuntimeError{
					Operator: expr.Operator,
					Message:  "Operands must be numbers.",
				}
			}

			switch expr.Operator.Type {
			case MINUS:
				return leftNum - rightNum, nil
			case STAR:
				return leftNum * rightNum, nil
			case SLASH:
				if rightNum == 0 {
					return nil, RuntimeError{
						Operator: expr.Operator,
						Message:  "Division by zero.",
					}
				}
				return leftNum / rightNum, nil
			case PERCENT:
				if rightNum == 0 {
					return nil, RuntimeError{
						Operator: expr.Operator,
						Message:  "Division by zero.",
					}
				}
				return math.Mod(leftNum, rightNum), nil
			case GREATER:
				return leftNum > rightNum, nil
			case GREATER_EQUAL:
				return leftNum >= rightNum, nil
			case LESS:
				return leftNum < rightNum, nil
			case LESS_EQUAL:
				return leftNum <= rightNum, nil
			}

		case EQUAL_EQUAL:
			return isEqual(leftVal, rightVal), nil

		case BANG_EQUAL:
			return !isEqual(leftVal, rightVal), nil
		}
	}

	// Temporary fallback for expression types or operators not yet implemented.
	// This should eventually be replaced with an appropriate error.
	return nil, nil
}

// isTruthy determines whether a runtime value is considered true or false.

func isTruthy(value any) bool {

	// Nil represents a falsey value.
	if value == nil {
		return false
	}

	// If the value is a boolean, return its actual boolean value.
	if boolean, ok := value.(bool); ok {
		return boolean
	}

	// Numbers, strings, and other non-nil, non-boolean values are truthy.
	return true
}

func isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	switch a := a.(type) {
	case float64:
		b, ok := b.(float64)
		return ok && a == b

	case string:
		b, ok := b.(string)
		return ok && a == b

	case bool:
		b, ok := b.(bool)
		return ok && a == b
	}

	return false
}
