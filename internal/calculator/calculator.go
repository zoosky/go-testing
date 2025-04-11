// Package calculator provides internal mathematical operations
package calculator

import (
	"go-testing/pkg/calculator"
)

// Calculator wraps the public calculator with any internal functionality
type Calculator struct {
	*calculator.Calculator
}

// NewCalculator creates a new Calculator instance
func NewCalculator() *Calculator {
	return &Calculator{
		Calculator: calculator.NewCalculator(),
	}
}