// Package calculator provides mathematical operations
package calculator

import (
	"errors"
)

// Calculator performs mathematical operations
type Calculator struct{}

// NewCalculator creates a new Calculator instance
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Add adds two numbers and returns the result
func (c *Calculator) Add(a, b float64) float64 {
	return a + b
}

// Subtract subtracts b from a and returns the result
func (c *Calculator) Subtract(a, b float64) float64 {
	return a - b
}

// Multiply multiplies two numbers and returns the result
func (c *Calculator) Multiply(a, b float64) float64 {
	return a * b
}

// Divide divides a by b and returns the result
// Returns an error if b is zero
func (c *Calculator) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}