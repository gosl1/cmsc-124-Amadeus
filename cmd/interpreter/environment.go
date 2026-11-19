package main

// Environment stores variable names and their runtime values.
type Environment struct {
	values map[string]any
}

// NewEnvironment creates an empty environment.
func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

// Define creates a variable or replaces its existing value.
func (env *Environment) Define(name Token, value any) {
	env.values[name.Lexeme] = value
}

// Get retrieves a variable's value.
func (env *Environment) Get(name Token) (any, error) {
	value, exists := env.values[name.Lexeme]
	if !exists {
		return nil, RuntimeError{
			Operator: name,
			Message:  "Undefined variable '" + name.Lexeme + "'.",
		}
	}

	return value, nil
}

// Assign updates an existing variable.
func (env *Environment) Assign(name Token, value any) error {
	if _, exists := env.values[name.Lexeme]; !exists {
		return RuntimeError{
			Operator: name,
			Message:  "Undefined variable '" + name.Lexeme + "'.",
		}
	}

	env.values[name.Lexeme] = value
	return nil
}
