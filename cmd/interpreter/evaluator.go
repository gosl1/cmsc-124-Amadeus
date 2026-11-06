package main

type RuntimeError struct {
	Operator Token
	Message  string
}

func (e RuntimeError) Error() string {
	return e.Message
}

type Evaluator struct {
}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Evaluate(expr Expr) (any, error) {
	switch expr := expr.(type) {
	case Literal:
		return expr.Value, nil

	case Grouping:
		return e.Evaluate(expr.Expression)

	case Unary:
		rightVal, err := e.Evaluate(expr.Right)
		if err != nil {
			return nil, err
		}

		switch expr.Operator.Type {
		case MINUS:
			number, ok := rightVal.(float64)
			if !ok {
				return nil, RuntimeError{
					Operator: expr.Operator,
					Message:  "Operand must be a number.",
				}
			}
			return -number, nil

		case BANG:
			return !isTruthy(rightVal), nil
		}

	}

	return nil, nil
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}

	if boolean, ok := value.(bool); ok {
		return boolean
	}

	return true
}
