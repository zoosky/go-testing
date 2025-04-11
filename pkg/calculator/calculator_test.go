package calculator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAdd tests the Add method with table-driven tests
func TestAdd(t *testing.T) {
	// Create a new calculator instance
	calc := NewCalculator()

	// Define test cases
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"Positive numbers", 2, 3, 5},
		{"Negative numbers", -2, -3, -5},
		{"Mixed numbers", -2, 3, 1},
		{"Zeros", 0, 0, 0},
		{"Decimals", 1.5, 2.5, 4},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := calc.Add(tc.a, tc.b)
			assert.Equal(t, tc.expected, result, "Add(%f, %f) should equal %f", tc.a, tc.b, tc.expected)
		})
	}
}

// TestSubtract tests the Subtract method with table-driven tests
func TestSubtract(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"Positive numbers", 5, 3, 2},
		{"Negative numbers", -5, -3, -2},
		{"Mixed numbers", -5, 3, -8},
		{"Zeros", 0, 0, 0},
		{"Decimals", 5.5, 2.5, 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := calc.Subtract(tc.a, tc.b)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestMultiply tests the Multiply method with table-driven tests
func TestMultiply(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"Positive numbers", 2, 3, 6},
		{"Negative numbers", -2, -3, 6},
		{"Mixed numbers", -2, 3, -6},
		{"Zeros", 0, 5, 0},
		{"Decimals", 1.5, 2, 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := calc.Multiply(tc.a, tc.b)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestDivide tests the Divide method with table-driven tests
func TestDivide(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name        string
		a, b        float64
		expected    float64
		expectError bool
	}{
		{"Positive numbers", 6, 3, 2, false},
		{"Negative numbers", -6, -3, 2, false},
		{"Mixed numbers", -6, 3, -2, false},
		{"Division by zero", 5, 0, 0, true},
		{"Decimals", 5, 2, 2.5, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := calc.Divide(tc.a, tc.b)
			
			if tc.expectError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

// Helper function example with t.Helper()
func assertOperationResult(t *testing.T, expected, actual float64, operation string, a, b float64) {
	t.Helper() // Marks this as a helper function for better error reporting
	if expected != actual {
		t.Errorf("%s(%f, %f): expected %f, got %f", operation, a, b, expected, actual)
	}
}

// TestCalculatorWithHelper demonstrates using a test helper
func TestCalculatorWithHelper(t *testing.T) {
	calc := NewCalculator()
	
	// Using our custom helper
	assertOperationResult(t, 5, calc.Add(2, 3), "Add", 2, 3)
	assertOperationResult(t, 2, calc.Subtract(5, 3), "Subtract", 5, 3)
	assertOperationResult(t, 6, calc.Multiply(2, 3), "Multiply", 2, 3)
	
	result, err := calc.Divide(6, 3)
	assert.NoError(t, err)
	assertOperationResult(t, 2, result, "Divide", 6, 3)
}